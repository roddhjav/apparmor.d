// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"slices"
	"strings"

	"github.com/roddhjav/apparmor.d/pkg/util"
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
	Kind() Kind             // Kind of the rule
	Constraint() Constraint // Where the rule can be found (preamble, profile, any)
	String() string         // Render the rule as a string
	Validate() error        // Validate the rule. Return an error if the rule is invalid
	Compare(other Rule) int // Compare two rules. Return 0 if they are identical
	Merge(other Rule) bool  // Merge rules of same kind together. Return true if merged
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

func (r Rules) DeleteKind(kind Kind) Rules {
	res := make(Rules, 0)
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

func (r Rules) Filter(filter Kind) Rules {
	res := make(Rules, 0)
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

func (r Rules) GetVariables() []*Variable {
	res := make([]*Variable, 0)
	for _, rule := range r {
		switch rule := rule.(type) {
		case *Variable:
			res = append(res, rule)
		}
	}
	return res
}

func (r Rules) GetIncludes() []*Include {
	res := make([]*Include, 0)
	for _, rule := range r {
		switch rule := rule.(type) {
		case *Include:
			res = append(res, rule)
		}
	}
	return res
}

// Merge merge similar rules together.
// Steps:
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

// Format the rules for better readability before printing it.
// Follow: https://apparmor.pujol.io/development/guidelines/#the-file-block
func (r Rules) Format() Rules {
	const prefixOwner = "      "
	suffixMaxlen := 36
	transitions := append(requirements[FILE]["transition"], "m")

	paddingIndex := []int{}
	paddingMaxLenght := 0
	for i, rule := range r {
		if rule == nil {
			continue
		}

		if rule.Kind() == FILE {
			rule := r[i].(*File)

			// Add padding to align with other transition rule
			isTransition := util.Intersect(transitions, rule.Access)
			if len(isTransition) > 0 {
				ruleLen := len(rule.Path) + 1
				paddingMaxLenght = max(ruleLen, paddingMaxLenght)
				paddingIndex = append(paddingIndex, i)
			}

			// Add suffix to align comment on udev/data rule
			if rule.Comment != "" && strings.HasPrefix(rule.Path, "@{run}/udev/data/") {
				suffixlen := suffixMaxlen - len(rule.Path)
				if suffixlen < 0 {
					suffixlen = 0
				}
				rule.Suffix = strings.Repeat(" ", suffixlen)
			}
		}
	}
	if len(paddingIndex) > 1 {
		r.setPadding(paddingIndex, paddingMaxLenght)
	}

	hasOwnerRule := false
	for i := len(r) - 1; i >= 0; i-- {
		if r[i] == nil {
			hasOwnerRule = false
			continue
		}

		// File rule
		if r[i].Kind() == FILE {
			rule := r[i].(*File)

			// Add prefix before rule path to align with other rule
			if rule.Owner {
				hasOwnerRule = true
			} else if hasOwnerRule {
				rule.Prefix = prefixOwner
			}

			// Do not add new line on executable rule
			isTransition := util.Intersect(transitions, rule.Access)
			if len(isTransition) > 0 {
				continue
			}

			// Add a new line between Files rule of different group type
			j := i - 1
			if j < 0 || r[j] == nil || r[j].Kind() != FILE {
				continue
			}
			letterI := getLetterIn(fileAlphabet, rule.Path)
			letterJ := getLetterIn(fileAlphabet, r[j].(*File).Path)
			groupI, ok1 := fileAlphabetGroups[letterI]
			groupJ, ok2 := fileAlphabetGroups[letterJ]
			if letterI != letterJ && !(ok1 && ok2 && groupI == groupJ) {
				hasOwnerRule = false
				r = r.Insert(i, nil)
			}
		}
	}
	return r
}

// setPadding adds padding to the rule path to align with other rules.
func (r *Rules) setPadding(paddingIndex []int, paddingMaxLenght int) {
	for _, i := range paddingIndex {
		(*r)[i].(*File).Padding = strings.Repeat(" ", paddingMaxLenght-len((*r)[i].(*File).Path))
	}
}
