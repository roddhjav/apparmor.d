// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2023-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package util

import (
	"encoding/hex"
	"regexp"
	"strings"
)

var (
	regHex = map[string]*regexp.Regexp{
		"name":    regexp.MustCompile(`name=[0-9A-F]+`),
		"comm":    regexp.MustCompile(`comm=[0-9A-F]+`),
		"profile": regexp.MustCompile(`profile=[0-9A-F]+`),

		// Unix socket address, abstract ones are logged as "@<hex>". Short
		// values are ignored as they are most likely a plain name.
		"addr":      regexp.MustCompile(`\baddr=(?:"@[0-9A-F]{5,}"|@?[0-9A-F]{5,}\b)`),
		"peer_addr": regexp.MustCompile(`\bpeer_addr=(?:"@[0-9A-F]{5,}"|@?[0-9A-F]{5,}\b)`),
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
		str = aa.Regex.ReplaceAllString(str, aa.Repl)
	}
	return str
}

// DecodeHexInString decode and replace all hex value in a given string of "key=value" format.
func DecodeHexInString(str string) string {
	for name, re := range regHex {
		str = re.ReplaceAllStringFunc(str, func(s string) string {
			hexa := strings.Trim(s[len(name)+1:], `"`)
			prefix := ""
			if strings.HasPrefix(hexa, "@") {
				prefix, hexa = "@", hexa[1:]
			}
			bs, err := hex.DecodeString(hexa)
			if err != nil || !isUntrusted(bs) {
				return s
			}
			// Abstract unix addresses may end with NUL bytes, each one is
			// matched by a ? in unix rules.
			value := string(bs)
			trimmed := strings.TrimRight(value, "\x00")
			value = trimmed + strings.Repeat("?", len(value)-len(trimmed))
			return name + "=\"" + prefix + value + "\""
		})
	}
	return str
}

// isUntrusted reports whether the kernel audit would have hex encoded bs: it
// contains a double quote, a space, a control or a non ASCII character.
func isUntrusted(bs []byte) bool {
	for _, b := range bs {
		if b == '"' || b < 0x21 || b > 0x7e {
			return true
		}
	}
	return false
}
