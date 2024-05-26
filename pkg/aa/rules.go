// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"fmt"
	"slices"
	"strings"
)

const (
	tokALLOW = "allow"
	tokAUDIT = "audit"
	tokDENY  = "deny"
)

type requirement map[string][]string

type constraint uint

const (
	anyKind      constraint = iota // The rule can be found in either preamble or profile
	preambleKind                   // The rule can only be found in the preamble
	blockKind                      // The rule can only be found in a profile
)

// Rule generic interface for all AppArmor rules
type Rule interface {
	Validate() error
	Less(other any) bool
	Equals(other any) bool
	String() string
	Constraint() constraint
	Kind() string
}

type Rules []Rule

func (r Rules) Validate() error {
	for _, rule := range r {
		if err := rule.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (r Rules) String() string {
	return renderTemplate("rules", r)
}

func (r Rules) IndexOf(rule Rule) int {
	for idx, rr := range r {
		if rr.Kind() == rule.Kind() && rr.Equals(rule) {
			return idx
		}
	}
	return -1
}

func (r Rules) Contains(rule Rule) bool {
	return r.IndexOf(rule) != -1
}

func (r Rules) Add(rule Rule) Rules {
	if r.Contains(rule) {
		return r
	}
	return append(r, rule)
}

func (r Rules) Remove(rule Rule) Rules {
	idx := r.IndexOf(rule)
	if idx == -1 {
		return r
	}
	return append(r[:idx], r[idx+1:]...)
}

func (r Rules) Insert(idx int, rules ...Rule) Rules {
	return append(r[:idx], append(rules, r[idx:]...)...)
}

func (r Rules) Sort() Rules {
	return r
}

func (r Rules) DeleteKind(kind string) Rules {
	res := make(Rules, 0)
	for _, rule := range r {
		if rule.Kind() != kind {
			res = append(res, rule)
		}
	}
	return res
}

func (r Rules) Filter(filter string) Rules {
	res := make(Rules, 0)
	for _, rule := range r {
		if rule.Kind() != filter {
			res = append(res, rule)
		}
	}
	return res
}

func (r Rules) GetVariables() []*Variable {
	res := make([]*Variable, 0)
	for _, rule := range r {
		switch rule.(type) {
		case *Variable:
			res = append(res, rule.(*Variable))
		}
	}
	return res
}

func (r Rules) GetIncludes() []*Include {
	res := make([]*Include, 0)
	for _, rule := range r {
		switch rule.(type) {
		case *Include:
			res = append(res, rule.(*Include))
		}
	}
	return res
}

// Must is a helper that wraps a call to a function returning (any, error) and
// panics if the error is non-nil.
func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func validateValues(rule string, key string, values []string) error {
	for _, v := range values {
		if v == "" {
			continue
		}
		if !slices.Contains(requirements[rule][key], v) {
			return fmt.Errorf("invalid mode '%s'", v)
		}
	}
	return nil
}

// Helper function to convert a string to a slice of rule values according to
// the rule requirements as defined in the requirements map.
func toValues(rule string, key string, input string) ([]string, error) {
	var sep string
	req, ok := requirements[rule][key]
	if !ok {
		return nil, fmt.Errorf("unrecognized requirement '%s' for rule %s", key, rule)
	}

	switch {
	case strings.Contains(input, ","):
		sep = ","
	case strings.Contains(input, " "):
		sep = " "
	}
	res := strings.Split(input, sep)
	for idx := range res {
		res[idx] = strings.Trim(res[idx], `" `)
		if !slices.Contains(req, res[idx]) {
			return nil, fmt.Errorf("unrecognized %s: %s", key, res[idx])
		}
	}
	slices.SortFunc(res, func(i, j string) int {
		return requirementsWeights[rule][key][i] - requirementsWeights[rule][key][j]
	})
	return slices.Compact(res), nil
}

// Helper function to convert an access string to a slice of access according to
// the rule requirements as defined in the requirements map.
func toAccess(rule string, input string) ([]string, error) {
	var res []string

	switch rule {
	case tokFILE:
		raw := strings.Split(input, "")
		trans := []string{}
		for _, access := range raw {
			if slices.Contains(requirements[tokFILE]["access"], access) {
				res = append(res, access)
			} else {
				trans = append(trans, access)
			}
		}

		transition := strings.Join(trans, "")
		if len(transition) > 0 {
			if slices.Contains(requirements[tokFILE]["transition"], transition) {
				res = append(res, transition)
			} else {
				return nil, fmt.Errorf("unrecognized transition: %s", transition)
			}
		}

	case tokFILE + "-log":
		raw := strings.Split(input, "")
		for _, access := range raw {
			if slices.Contains(requirements[tokFILE]["access"], access) {
				res = append(res, access)
			} else if maskToAccess[access] != "" {
				res = append(res, maskToAccess[access])
			} else {
				return nil, fmt.Errorf("toAccess: unrecognized file access '%s'", input)
			}
		}

	default:
		return toValues(rule, "access", input)
	}

	slices.SortFunc(res, cmpFileAccess)
	return slices.Compact(res), nil
}
