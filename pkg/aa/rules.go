// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"slices"
)

type requirement map[string][]string

type constraint uint

const (
	anyKind      constraint = iota // The rule can be found in either preamble or profile
	preambleKind                   // The rule can only be found in the preamble
	blockKind                      // The rule can only be found in a profile
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
	Validate() error
	Less(other any) bool
	Equals(other any) bool
	String() string
	Constraint() constraint
	Kind() Kind
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
		if rule.Kind() == item.Kind() && rule.Equals(item) {
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
			typeOfI := r[i].Kind()
			typeOfJ := r[j].Kind()
			if typeOfI != typeOfJ {
				continue
			}

			// If rules are identical, merge them
			if r[i].Equals(r[j]) {
				r = r.Delete(j)
				j--
				continue
			}

			// File rule
			if typeOfI == FILE && typeOfJ == FILE {
				// Merge access
				fileI := r[i].(*File)
				fileJ := r[j].(*File)
				if fileI.Path == fileJ.Path {
					fileI.Access = append(fileI.Access, fileJ.Access...)
					slices.SortFunc(fileI.Access, cmpFileAccess)
					fileI.Access = slices.Compact(fileI.Access)
					r = r.Delete(j)
					j--
				}
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
		if a.Equals(b) {
			return 0
		}
		if a.Less(b) {
			return -1
		}
		return 1
	})
	return r
}

// Format the rules for better readability before printing it.
// Follow: https://apparmor.pujol.io/development/guidelines/#the-file-block
func (r Rules) Format() Rules {
	const prefixOwner = "      "

	hasOwnerRule := false
	for i := len(r) - 1; i > 0; i-- {
		j := i - 1
		typeOfI := r[i].Kind()
		typeOfJ := r[j].Kind()

		// File rule
		if typeOfI == FILE && typeOfJ == FILE {
			letterI := getLetterIn(fileAlphabet, r[i].(*File).Path)
			letterJ := getLetterIn(fileAlphabet, r[j].(*File).Path)

			// Add prefix before rule path to align with other rule
			if r[i].(*File).Owner {
				hasOwnerRule = true
			} else if hasOwnerRule {
				r[i].(*File).Prefix = prefixOwner
			}

			if letterI != letterJ {
				// Add a new empty line between Files rule of different type
				hasOwnerRule = false
				r = r.Insert(i, nil)
			}
		}
	}
	return r
}
