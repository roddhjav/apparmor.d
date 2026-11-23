// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2023-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package util

import (
	"encoding/hex"
	"regexp"
	"strings"
)

// regHex matches the fields the kernel may log as hex. Quoted values are
// matched first so they are skipped, and a key must start a field: a key
// inside a value (fsname=, a decoded text) is never decoded.
var regHex = regexp.MustCompile(`"[^"]*"|(^|\s)(name|srcname|comm|profile|addr|peer_addr)=("@[0-9A-F]+"|@?[0-9A-F]+\b)`)

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
	return regHex.ReplaceAllStringFunc(str, func(s string) string {
		m := regHex.FindStringSubmatch(s)
		if m[2] == "" { // Quoted value
			return s
		}
		lead, name, hexa := m[1], m[2], strings.Trim(m[3], `"`)
		prefix := ""
		if strings.HasPrefix(hexa, "@") {
			prefix, hexa = "@", hexa[1:]
		}
		// Short addr values are most likely a plain name.
		if strings.HasSuffix(name, "addr") && len(hexa) < 5 {
			return s
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
		return lead + name + "=\"" + prefix + value + "\""
	})
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
