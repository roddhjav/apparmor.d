// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2023-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package util

import (
	"encoding/hex"
	"regexp"
)

var (
	Comment   = `#`
	regFilter = ToRegexRepl([]string{
		`\s*` + Comment + `.*`, ``,
		`(?m)^(?:[\t\s]*(?:\r?\n|\r))+`, ``,
	})
	regHex = map[string]*regexp.Regexp{
		"name":    regexp.MustCompile(`name=[0-9A-F]+`),
		"comm":    regexp.MustCompile(`comm=[0-9A-F]+`),
		"profile": regexp.MustCompile(`profile=[0-9A-F]+`),
	}
)

type RegexReplList []RegexRepl

type RegexRepl struct {
	Regex *regexp.Regexp
	Repl  string
}

// ToRegexRepl convert slice of regex into a slice of RegexRepl
func ToRegexRepl(in []string) RegexReplList {
	out := make([]RegexRepl, 0, len(in)/2)
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

// DecodeHexInString decode and replace all hex value in a given string of "key=value" format.
func DecodeHexInString(str string) string {
	for name, re := range regHex {
		str = re.ReplaceAllStringFunc(str, func(s string) string {
			hexa := s[len(name)+1:]
			bs, _ := hex.DecodeString(hexa)
			return name + "=\"" + string(bs) + "\""
		})
	}
	return str
}

// Filter out comments and empty line from a string
func Filter(src string) string {
	return regFilter.Replace(src)
}
