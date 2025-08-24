// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"fmt"
	"slices"
	"strings"
)

// Must is a helper that wraps a call to a function returning (any, error) and
// panics if the error is non-nil.
func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func merge(kind Kind, key string, a, b []string) []string {
	a = append(a, b...)
	switch kind {
	case FILE:
		slices.SortFunc(a, compareFileAccess)
	case VARIABLE:
		slices.SortFunc(a, func(s1, s2 string) int {
			return compare(s1, s2)
		})
	default:
		slices.SortFunc(a, func(i, j string) int {
			return requirementsWeights[kind][key][i] - requirementsWeights[kind][key][j]
		})
	}
	return slices.Compact(a)
}

func length(prefix string, value any) int {
	var res int
	switch value := value.(type) {
	case bool:
		if value {
			return len(prefix) + 1
		}
		return 0
	case string:
		if value != "" {
			res = len(value) + len(prefix) + 1
		}
		return res
	case []string:
		for _, v := range value {
			lenV := len(v)
			if lenV > 0 {
				res += lenV + 1 // Space between values
			}
		}
		if len(value) > 1 {
			res += 2 // Brackets on slices
		}
		return res
	default:
		panic("length: unsupported type")
	}
}

func setPaddings(max []int, prefixes []string, values []any) []string {
	if len(max) != len(values) || len(max) != len(prefixes) {
		panic("setPaddings: max, prefix, and values must have the same length")
	}
	res := make([]string, len(max))
	for i, v := range values {
		if max[i] == 0 {
			res[i] = ""
			continue
		}
		count := max[i] - length(prefixes[i], v)
		if count > 0 {
			res[i] = strings.Repeat(" ", count)
		}
	}
	return res
}

func compare(a, b any) int {
	switch a := a.(type) {
	case int:
		return a - b.(int)
	case string:
		a = strings.ToLower(a)
		b := strings.ToLower(b.(string))
		if a == b {
			return 0
		}
		for i := 0; i < len(a) && i < len(b); i++ {
			if a[i] != b[i] {
				return stringWeights[a[i]] - stringWeights[b[i]]
			}
		}
		return len(a) - len(b)
	case []string:
		return slices.CompareFunc(a, b.([]string), func(s1, s2 string) int {
			return compare(s1, s2)
		})
	case bool:
		return boolToInt(a) - boolToInt(b.(bool))
	default:
		panic("compare: unsupported type")
	}
}

// compareFileAccess compares two access strings for file rules.
// It is aimed to be used in slices.SortFunc.
func compareFileAccess(i, j string) int {
	if slices.Contains(requirements[FILE]["access"], i) &&
		slices.Contains(requirements[FILE]["access"], j) {
		return requirementsWeights[FILE]["access"][i] - requirementsWeights[FILE]["access"][j]
	}
	if slices.Contains(requirements[FILE]["transition"], i) &&
		slices.Contains(requirements[FILE]["transition"], j) {
		return requirementsWeights[FILE]["transition"][i] - requirementsWeights[FILE]["transition"][j]
	}
	if slices.Contains(requirements[FILE]["access"], i) {
		return -1
	}
	return 1
}

func validateValues(kind Kind, key string, values []string) error {
	for _, v := range values {
		if v == "" {
			continue
		}
		if !slices.Contains(requirements[kind][key], v) {
			return fmt.Errorf("invalid mode '%s'", v)
		}
	}
	return nil
}

func tokenToSlice(token string) []string {
	res := []string{}
	token = strings.Trim(token, "()\n ")
	if strings.ContainsAny(token, ", ") {
		var sep string
		token = strings.ReplaceAll(token, "  ", " ")
		switch {
		case strings.Contains(token, ","):
			sep = ","
		case strings.Contains(token, " "):
			sep = " "
		}
		for _, v := range strings.Split(token, sep) {
			res = append(res, strings.Trim(v, " "))
		}
	} else {
		res = append(res, token)
	}
	return res
}

// Helper function to convert a string to a slice of rule values according to
// the rule requirements as defined in the requirements map.
func toValues(kind Kind, key string, input string) ([]string, error) {
	req, ok := requirements[kind][key]
	if !ok {
		return nil, fmt.Errorf("unrecognized requirement '%s' for rule %s", key, kind)
	}

	res := tokenToSlice(input)
	for idx := range res {
		res[idx] = strings.Trim(res[idx], `" `)
		if res[idx] == "" {
			res = slices.Delete(res, idx, idx+1)
			continue
		}
		if !slices.Contains(req, res[idx]) {
			return nil, fmt.Errorf("unrecognized %s for rule %s: %s", key, kind, res[idx])
		}
	}
	slices.SortFunc(res, func(i, j string) int {
		return requirementsWeights[kind][key][i] - requirementsWeights[kind][key][j]
	})
	return slices.Compact(res), nil
}

// Helper function to convert an access string to a slice of access according to
// the rule requirements as defined in the requirements map.
func toAccess(kind Kind, input string) ([]string, error) {
	var res []string

	switch kind {
	case FILE:
		raw := strings.Split(input, "")
		trans := []string{}
		for _, access := range raw {
			if slices.Contains(requirements[FILE]["access"], access) {
				res = append(res, access)
			} else {
				trans = append(trans, access)
			}
		}

		transition := strings.Join(trans, "")
		if len(transition) > 0 {
			if slices.Contains(requirements[FILE]["transition"], transition) {
				res = append(res, transition)
			} else {
				return nil, fmt.Errorf("unrecognized transition: %s", transition)
			}
		}

	case FILE + "-log":
		raw := strings.Split(input, "")
		for _, access := range raw {
			if slices.Contains(requirements[FILE]["access"], access) {
				res = append(res, access)
			} else if maskToAccess[access] != "" {
				res = append(res, maskToAccess[access])
			} else {
				return nil, fmt.Errorf("toAccess: unrecognized file access '%s' for %s", input, kind)
			}
		}

	default:
		return toValues(kind, "access", input)
	}

	slices.SortFunc(res, compareFileAccess)
	return slices.Compact(res), nil
}
