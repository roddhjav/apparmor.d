// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2023-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package util

import (
	"encoding/hex"
	"regexp"
	"slices"
	"strings"

	"github.com/roddhjav/apparmor.d/pkg/paths"
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

// RemoveDuplicate filter out all duplicates from a slice. Also filter out empty element.
func RemoveDuplicate[T comparable](inlist []T) []T {
	var empty T
	list := []T{}
	seen := map[T]bool{}
	seen[empty] = true
	for _, item := range inlist {
		if _, ok := seen[item]; !ok {
			seen[item] = true
			list = append(list, item)
		}
	}
	return list
}

// Intersect returns the intersection between two collections.
// From https://github.com/samber/lo
func Intersect[T comparable](list1 []T, list2 []T) []T {
	result := []T{}
	seen := map[T]struct{}{}

	for _, elem := range list1 {
		seen[elem] = struct{}{}
	}

	for _, elem := range list2 {
		if _, ok := seen[elem]; ok {
			result = append(result, elem)
		}
	}

	return result
}

// CopyTo recursivelly copy all files from a source path to a destination path.
func CopyTo(src *paths.Path, dst *paths.Path) error {
	files, err := src.ReadDirRecursiveFiltered(nil,
		paths.FilterOutDirectories(),
		paths.FilterOutNames("README.md"),
	)
	if err != nil {
		return err
	}
	for _, file := range files {
		destination, err := file.RelFrom(src)
		if err != nil {
			return err
		}
		destination = dst.JoinPath(destination)
		if err := destination.Parent().MkdirAll(); err != nil {
			return err
		}
		if err := file.CopyTo(destination); err != nil {
			return err
		}
	}
	return nil
}

// Filter out comments and empty line from a string
func Filter(src string) string {
	return regFilter.Replace(src)
}

// ReadFile read a file and return its content as a string.
func ReadFile(path *paths.Path) (string, error) {
	content, err := path.ReadFile()
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// MustReadFile read a file and return its content as a string. Panic if an error occurs.
func MustReadFile(path *paths.Path) string {
	content, err := path.ReadFile()
	if err != nil {
		panic(err)
	}
	return string(content)
}

// MustReadFileAsLines read a file and return its content as a slice of string.
// It panics if an error occurs and filter out comments and empty lines.
func MustReadFileAsLines(path *paths.Path) []string {
	res := strings.Split(Filter(MustReadFile(path)), "\n")
	if slices.Contains(res, "") {
		idx := slices.Index(res, "")
		res = slices.Delete(res, idx, idx+1)
	}
	return res
}
