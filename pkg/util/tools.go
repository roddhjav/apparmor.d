// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2023 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package util

import (
	"encoding/hex"
	"regexp"
)

var isHexa = regexp.MustCompile("^[0-9A-Fa-f]+$")

type RegexRepl struct {
	Regex *regexp.Regexp
	Repl  string
}

// DecodeHex decode a string if it is hexa.
func DecodeHex(str string) string {
	if isHexa.MatchString(str) {
		bs, _ := hex.DecodeString(str)
		return string(bs)
	}
	return str
}

// RemoveDuplicate filter out all duplicates from a slice. Also filter out empty string
func RemoveDuplicate[T comparable](inlist []T) []T {
	var empty T
	list := []T{}
	keys := map[T]bool{}
	keys[empty] = true
	for _, item := range inlist {
		if _, ok := keys[item]; !ok {
			keys[item] = true
			list = append(list, item)
		}
	}
	return list
}

// ToRegexRepl convert slice of regex into a slice of RegexRepl
func ToRegexRepl(in []string) []RegexRepl {
	out := make([]RegexRepl, 0)
	idx := 0
	for idx < len(in)-1 {
		regex, repl := in[idx], in[idx+1]
		out = append(out, RegexRepl{
			Regex: regexp.MustCompile(regex),
			Repl:  repl,
		})
		idx = idx + 2
	}
	return out
}
