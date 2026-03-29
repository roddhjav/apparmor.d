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

// boolToInt converts a boolean to an integer (true = 1, false = 0).
func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// merge merges two slices of strings according to their kind (file, variable, other) and requirements matrix.
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

// length returns the length of a value with its prefix for padding calculation.
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
		if res > 0 && prefix != "" {
			res += len(prefix) // Add prefix length when slice is non-empty
		}
		return res
	default:
		panic("length: unsupported type")
	}
}

// setPaddings returns a slice of paddings for each value to align them according to the max lengths provided.
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

// compare compares two values of the same type and returns -1, 0, or 1
// if a is less than, equal to, or greater than b, respectively.
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
		b := b.([]string)
		// Compare using the formatted representation so that parenthesized
		// multi-element values (e.g., "(kill term)") sort before single
		// element values (e.g., "hup") based on the '(' character weight.
		sa := strings.Join(a, " ")
		if len(a) > 1 {
			sa = "(" + sa + ")"
		}
		sb := strings.Join(b, " ")
		if len(b) > 1 {
			sb = "(" + sb + ")"
		}
		return compare(sa, sb)
	case bool:
		return boolToInt(a) - boolToInt(b.(bool))
	default:
		panic("compare: unsupported type")
	}
}

// compareFileAccess compares two access strings for file rules.
// It is aimed to be used in slices.SortFunc.
func compareFileAccess(i, j string) int {
	accessWeights := requirementsWeights[FILE]["access"]
	transitionWeights := requirementsWeights[FILE]["transition"]
	wi, iIsAccess := accessWeights[i]
	wj, jIsAccess := accessWeights[j]
	if iIsAccess && jIsAccess {
		return wi - wj
	}
	wi, iIsTransition := transitionWeights[i]
	wj, jIsTransition := transitionWeights[j]
	if iIsTransition && jIsTransition {
		return wi - wj
	}
	if iIsAccess {
		return -1
	}
	return 1
}

// validateValues checks if all values in the slice are valid according to the requirements matrix.
func validateValues(kind Kind, key string, values []string) error {
	for _, v := range values {
		if v == "" {
			continue
		}

		// Skip variable references — they will be expanded at runtime
		if strings.Contains(v, "@{") {
			continue
		}

		v = strings.Trim(v, "`\"") // Strip surrounding quotes for validation
		if !slices.Contains(requirements[kind][key], v) {
			// Check for prefix-based values (e.g., "kill.signal=hup" matches "kill.signal=")
			found := false
			for _, req := range requirements[kind][key] {
				if strings.HasSuffix(req, "=") && strings.HasPrefix(v, req) {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("invalid mode '%s'", v)
			}
		}
	}
	return nil
}

// flagMatch checks if a value matches a conflict pattern.
// Supports exact match and prefix match for patterns ending with "=".
func flagMatch(value, pattern string) bool {
	if value == pattern {
		return true
	}
	// Prefix match: "attach_disconnected.ipc=/foo" matches "attach_disconnected.ipc="
	if strings.HasSuffix(pattern, "=") && strings.HasPrefix(value, pattern) {
		return true
	}
	// Also match the other way: value "attach_disconnected.ipc=" prefix matches actual value
	if strings.HasSuffix(value, "=") && strings.HasPrefix(pattern, value) {
		return true
	}
	return false
}

// validateConflicts checks if any values in the slice conflict with each other.
// Conflicts are defined in the conflicts map as pairs of mutually exclusive values.
func validateConflicts(kind Kind, key string, values []string) error {
	conflictPairs, ok := conflicts[kind][key]
	if !ok {
		return nil
	}
	for _, pair := range conflictPairs {
		if len(pair) != 2 {
			continue
		}
		hasFirst := false
		hasSecond := false
		for _, v := range values {
			if flagMatch(v, pair[0]) {
				hasFirst = true
			}
			if flagMatch(v, pair[1]) {
				hasSecond = true
			}
		}
		if hasFirst && hasSecond {
			return fmt.Errorf("conflicting %s '%s' and '%s'", key, pair[0], pair[1])
		}
	}
	return nil
}

// validateAAREPattern checks for invalid AARE (AppArmor Regular Expression) patterns.
func validateAAREPattern(path string) error {
	// Check for empty character class: []
	if strings.Contains(path, "[]") {
		return fmt.Errorf("empty character class '[]' in path '%s'", path)
	}
	// Check for empty alternation: {}
	if strings.Contains(path, "{}") {
		return fmt.Errorf("empty alternation '{}' in path '%s'", path)
	}
	// Check for single-entry alternation: {word} (no comma inside)
	// Skip variable references @{...} and handle nesting
	for i := 0; i < len(path); i++ {
		if path[i] == '{' && (i == 0 || (path[i-1] != '@' && path[i-1] != '\\')) {
			// Find matching closing brace, accounting for nesting
			depth := 1
			end := -1
			for j := i + 1; j < len(path); j++ {
				if path[j] == '{' {
					depth++
				} else if path[j] == '}' {
					depth--
					if depth == 0 {
						end = j
						break
					}
				}
			}
			if end > 0 {
				inner := path[i+1 : end]
				if !strings.Contains(inner, ",") {
					return fmt.Errorf("single-entry alternation '{%s}' in path '%s'", inner, path)
				}
			}
		}
	}
	return nil
}

// tokenToSlice splits a token string into a slice of strings based on commas or spaces.
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

	tokens := tokenToSlice(input)
	res := make([]string, 0, len(tokens))
	for _, token := range tokens {
		token = strings.Trim(token, `" `)
		if token == "" {
			continue
		}
		if !slices.Contains(req, token) {
			return nil, fmt.Errorf("unrecognized %s for rule %s: %s", key, kind, token)
		}
		res = append(res, token)
	}
	slices.SortFunc(res, func(i, j string) int {
		return requirementsWeights[kind][key][i] - requirementsWeights[kind][key][j]
	})
	return slices.Compact(res), nil
}

// Helper function to convert an access string to a slice of access according to
// the rule requirements as defined in the requirements matrix.
func toAccess(kind Kind, input string) ([]string, error) {
	var res []string

	switch kind {
	case FILE:
		accessWeights := requirementsWeights[FILE]["access"]
		transitionWeights := requirementsWeights[FILE]["transition"]
		raw := strings.Split(input, "")
		trans := []string{}
		// Track positions of access vs transition chars to detect interleaving
		lastTransPos := -1
		firstAccessAfterTrans := false
		for i, access := range raw {
			if _, ok := accessWeights[access]; ok {
				res = append(res, access)
				if lastTransPos >= 0 {
					firstAccessAfterTrans = true
				}
			} else {
				_ = i
				if firstAccessAfterTrans {
					// Transition char after access char that was after a transition char
					// e.g., "prx" → p(trans) r(access) x(trans) — invalid interleaving
					return nil, fmt.Errorf("invalid access mode: access and transition chars interleaved in '%s'", input)
				}
				trans = append(trans, access)
				lastTransPos = i
			}
		}

		transition := strings.Join(trans, "")
		if len(transition) > 0 {
			resolved := resolveTransition(transition, transitionWeights)
			if resolved != "" {
				res = append(res, resolved)
			} else {
				return nil, fmt.Errorf("unrecognized transition: %s", transition)
			}
		}

	case FILE + "-log":
		accessWeights := requirementsWeights[FILE]["access"]
		raw := strings.Split(input, "")
		for _, access := range raw {
			if _, ok := accessWeights[access]; ok {
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

// resolveTransition tries to match a transition string against known transitions.
// It handles exact match, case-insensitive match, and repeated transitions (e.g., "pxpxpx" → "px").
func resolveTransition(transition string, weights map[string]int) string {
	// Exact match
	if _, ok := weights[transition]; ok {
		return transition
	}
	// Case-insensitive match
	for t := range weights {
		if strings.EqualFold(transition, t) {
			return t
		}
	}
	// Check if it's a repeated valid transition (e.g., "pxpxpx" → "px")
	for t := range weights {
		if len(t) > 0 && len(transition) > len(t) && len(transition)%len(t) == 0 {
			repeated := strings.Repeat(t, len(transition)/len(t))
			if repeated == transition || strings.EqualFold(repeated, transition) {
				return t
			}
		}
	}
	return ""
}

// allLocalOnly returns true if all access types are in the local-only list.
func allLocalOnly(access, localOnly []string) bool {
	for _, a := range access {
		if !slices.Contains(localOnly, a) {
			return false
		}
	}
	return true
}
