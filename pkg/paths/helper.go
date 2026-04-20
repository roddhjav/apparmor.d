// This file is part of PathsHelper library.
// Copyright (C) 2018-2025 Arduino AG (http://www.arduino.cc/)
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package paths

import (
	"fmt"
	"regexp"
)

var (
	Comment   = `#`
	regFilter = []*regexp.Regexp{
		regexp.MustCompile(`\s*` + Comment + `.*`),
		regexp.MustCompile(`(?m)^(?:[\t\s]*(?:\r?\n|\r))+`),
	}
)

// Filter out comments and empty lines from a string.
func Filter(src string) string {
	for _, re := range regFilter {
		src = re.ReplaceAllLiteralString(src, "")
	}
	return src
}

// PathListFromArgs resolves CLI-style arguments into a PathList. Each arg may
// be a file path, a directory (recursed into, skipping README.md), or a bare
// name looked up under magicRoot.
func PathListFromArgs(args []string, magicRoot *Path) (PathList, error) {
	res := PathList{}
	for _, arg := range args {
		path := New(arg)
		switch {
		case !path.Exist():
			magic := magicRoot.Join(arg)
			if !magic.Exist() {
				return nil, fmt.Errorf("file %s not found", path)
			}
			res = append(res, magic)
		case path.IsDir():
			files, err := path.ReadDirRecursiveFiltered(nil,
				FilterOutDirectories(),
				FilterOutNames("README.md"),
			)
			if err != nil {
				return nil, err
			}
			res = append(res, files...)
		default:
			res = append(res, path)
		}
	}
	return res, nil
}
