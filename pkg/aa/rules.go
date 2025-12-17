// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"slices"
)

type requirement map[string][]string

type Constraint uint

const (
	AnyRule      Constraint = iota // The rule can be found in either preamble or profile
	PreambleRule                   // The rule can only be found in the preamble
	BlockRule                      // The rule can only be found in a profile
)

// Kind represents an AppArmor rule kind.
type Kind string

func (k Kind) String() string {
	return string(k)
}

func (k Kind) Tok() string {
	if t, ok := tok[k]; ok {
		return t
	}
	return string(k)
}

// Rule generic interface for all AppArmor rules
type Rule interface {
	Kind() Kind              // Kind of the rule
	Constraint() Constraint  // Where the rule can be found (preamble, profile, any)
	String() string          // Render the rule as a string
	Validate() error         // Validate the rule. Return an error if the rule is invalid
	Compare(other Rule) int  // Compare two rules. Return 0 if they are identical
	Merge(other Rule) bool   // Merge rules of same kind together. Return true if merged
	Padding(i int) string    // Padding for rule items at index i
	Lengths() []int          // Length of each item in the rule
	setPaddings(max []int)   // Set paddings for each item in the rule
	addLine(other Rule) bool // Check either a new line should be added before the rule
}

type Rules []Rule

func (r Rules) Validate() error {
	for _, rule := range r {
		if rule == nil {
			continue
		}
		if err := rule.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (r Rules) String() string {
	return renderTemplate("rules", r)
}

// Index returns the index of the first occurrence of rule rin r, or -1 if not present.
func (r Rules) Index(item Rule) int {
	for idx, rule := range r {
		if rule == nil {
			continue
		}
		if rule.Kind() == item.Kind() && rule.Compare(item) == 0 {
			return idx
		}
	}
	return -1
}

// IndexOf returns the index of the first occurrence of item in r, or -1 if not present.
func (r Rules) IndexOf(item Rule) int {
	for idx, rr := range r {
		if rr.Kind() == item.Kind() && rr.Compare(item) == 0 {
			return idx
		}
	}
	return -1
}

// Contains checks if the rule is in the slice
func (r Rules) Contains(rule Rule) bool {
	return r.IndexOf(rule) != -1
}

// Remove removes the first occurrence of rule from the slice and returns the new slice.
func (r Rules) Remove(rule Rule) Rules {
	idx := r.IndexOf(rule)
	if idx == -1 {
		return r
	}
	return append(r[:idx], r[idx+1:]...)
}

// Replace replaces the elements r[i] by the given rules, and returns the
// modified slice.
func (r Rules) Replace(i int, rules ...Rule) Rules {
	return append(r[:i], append(rules, r[i+1:]...)...)
}

// Insert inserts the rules into r at index i, returning the modified slice.
func (r Rules) Insert(i int, rules ...Rule) Rules {
	return append(r[:i], append(rules, r[i:]...)...)
}

// Delete removes the elements r[i] from r, returning the modified slice.
func (r Rules) Delete(i int) Rules {
	return append(r[:i], r[i+1:]...)
}

// DeleteKind removes all rules of the given kind from the slice and returns the new slice.
func (r Rules) DeleteKind(kind Kind) Rules {
	res := make(Rules, 0, len(r))
	for _, rule := range r {
		if rule == nil {
			continue
		}
		if rule.Kind() != kind {
			res = append(res, rule)
		}
	}
	return res
}

// FilterOut removes all rules of the given kind from the slice and returns the new slice.
func (r Rules) FilterOut(filter Kind) Rules {
	res := make(Rules, 0, len(r))
	for _, rule := range r {
		if rule == nil {
			continue
		}
		if rule.Kind() != filter {
			res = append(res, rule)
		}
	}
	return res
}

// Filter returns all rules of the given kind from the slice.
func (r Rules) Filter(filter Kind) Rules {
	res := make(Rules, 0, len(r))
	for _, rule := range r {
		if rule == nil {
			continue
		}
		if rule.Kind() == filter {
			res = append(res, rule)
		}
	}
	return res
}

// GetVariables returns all Variable rules from the slice.
func (r Rules) GetVariables() []*Variable {
	res := make([]*Variable, 0, len(r))
	for _, rule := range r {
		switch rule := rule.(type) {
		case *Variable:
			res = append(res, rule)
		}
	}
	return res
}

// GetIncludes returns all Include rules from the slice.
func (r Rules) GetIncludes() []*Include {
	res := make([]*Include, 0, len(r))
	for _, rule := range r {
		switch rule := rule.(type) {
		case *Include:
			res = append(res, rule)
		}
	}
	return res
}

// Merge merge similar rules together:
//   - Remove identical rules
//   - Merge rule access. Eg: for same path, 'r' and 'w' becomes 'rw'
//
// Note: logs.regCleanLogs helps a lot to do a first cleaning
func (r Rules) Merge() Rules {
	for i := 0; i < len(r); i++ {
		for j := i + 1; j < len(r); j++ {
			if r[i] == nil && r[j] == nil {
				r = r.Delete(j)
				j--
				continue
			}
			if r[i] == nil || r[j] == nil {
				continue
			}
			if r[i].Kind() != r[j].Kind() {
				continue
			}

			// If rules are identical, merge them. Ignore comments
			if r[i].Kind() != COMMENT && r[i].Compare(r[j]) == 0 {
				r = r.Delete(j)
				j--
				continue
			}

			if r[i].Merge(r[j]) {
				r = r.Delete(j)
				j--
			}
		}
	}
	return r
}

// Sort the rules according to the guidelines:
// https://apparmor.pujol.io/development/guidelines/#guidelines
func (r Rules) Sort() Rules {
	slices.SortFunc(r, func(a, b Rule) int {
		kindOfA := a.Kind()
		kindOfB := b.Kind()
		if kindOfA != kindOfB {
			if kindOfA == INCLUDE && a.(*Include).IfExists {
				kindOfA = "include_if_exists"
			}
			if kindOfB == INCLUDE && b.(*Include).IfExists {
				kindOfB = "include_if_exists"
			}
			return ruleWeights[kindOfA] - ruleWeights[kindOfB]
		}
		return a.Compare(b)
	})
	return r
}

// setPaddings set paddings for each element in each rules
func (r *Rules) setPaddings(paddingsIndex map[Kind][]int, paddingsMaxLen map[Kind][]int) {
	for kind, index := range paddingsIndex {
		if len(index) <= 1 {
			continue
		}
		for _, i := range index {
			(*r)[i].setPaddings(paddingsMaxLen[kind])
		}
	}
}

// Format the rules for better readability before printing it. Format supposes
// the rules are merged and sorted.
// Follow: https://apparmor.pujol.io/development/guidelines/#the-file-block
func (r Rules) Format() Rules {
	// Insert new line between rule of different type/subtype.
	for i := len(r) - 1; i >= 0; i-- {
		j := i - 1
		if j < 0 || r[j] == nil {
			continue
		}
		if r[i].addLine(r[j]) {
			r = r.Insert(i, nil)
		}
	}

	// Find max paddings for each element in each rules
	paddingsIndex := map[Kind][]int{}
	paddingsMaxLen := map[Kind][]int{}
	for i, rule := range r {
		if rule == nil {
			r.setPaddings(paddingsIndex, paddingsMaxLen)
			paddingsIndex = map[Kind][]int{}
			paddingsMaxLen = map[Kind][]int{}
			continue
		}

		lengths := rule.Lengths()
		paddingsIndex[rule.Kind()] = append(paddingsIndex[rule.Kind()], i)
		for idx, length := range lengths {
			if _, ok := paddingsMaxLen[rule.Kind()]; !ok {
				paddingsMaxLen[rule.Kind()] = make([]int, len(lengths))
			}
			paddingsMaxLen[rule.Kind()][idx] = max(paddingsMaxLen[rule.Kind()][idx], length)
		}
	}
	r.setPaddings(paddingsIndex, paddingsMaxLen)
	return r
}

// ParaRules is a slice of Rules grouped by paragraph
type ParaRules []Rules

// Flatten flattens the ParaRules into a single Rules slice
func (r ParaRules) Flatten() Rules {
	totalLen := 0
	for i := range r {
		totalLen += len(r[i])
	}

	res := make(Rules, 0, totalLen)
	for i := range r {
		res = append(res, r[i]...)
	}

	return res
}
