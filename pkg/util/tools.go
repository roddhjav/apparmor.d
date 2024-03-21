// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2023-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package util

import (
	"encoding/hex"
	"regexp"
)

type RegexReplList []RegexRepl

type RegexRepl struct {
	Regex *regexp.Regexp
	Repl  string
}

// ToRegexRepl convert slice of regex into a slice of RegexRepl
func ToRegexRepl(in []string) RegexReplList {
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

func (rr RegexReplList) Replace(str string) string {
	for _, aa := range rr {
		str = aa.Regex.ReplaceAllLiteralString(str, aa.Repl)
	}
	return str
}

// DecodeHexInString decode and replace all hex value in a given string constitued of "key=value".
func DecodeHexInString(str string) string {
	toDecode := []string{"name", "comm", "profile"}
	for _, name := range toDecode {
		exp := name + `=[0-9A-F]+`
		re := regexp.MustCompile(exp)
		str = re.ReplaceAllStringFunc(str, func(s string) string {
			hexa := s[len(name)+1:]
			bs, _ := hex.DecodeString(hexa)
			return name + "=\"" + string(bs) + "\""
		})
	}
	return str
}

// RemoveDuplicate filter out all duplicates from a slice. Also filter out empty element.
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
