// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2023 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package util

import (
	"encoding/hex"
	"regexp"
)

var isHexa = regexp.MustCompile("^[0-9A-Fa-f]+$")

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
