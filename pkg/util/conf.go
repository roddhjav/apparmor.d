// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package util

import (
	"maps"
	"strings"

	"github.com/roddhjav/apparmor.d/pkg/paths"
)

// ReadFlagsFile parses a flags manifest file: one "profile flag1,flag2" per
// line, # comments filtered. A profile without flags maps to an empty list.
func ReadFlagsFile(path *paths.Path) map[string][]string {
	res := map[string][]string{}
	if !path.Exist() {
		return res
	}
	for _, line := range path.MustReadFilteredFileAsLines() {
		manifest := strings.Split(line, " ")
		profile := manifest[0]
		flags := []string{}
		if len(manifest) > 1 {
			flags = strings.Split(manifest[1], ",")
		}
		res[profile] = flags
	}
	return res
}

// ReadFlagDirs merges every flags manifest found in sources. A source is
// either a flags file, read as is, or a directory whose *.conf files are
// read in alphabetical order. Sources are read in the given order, so a
// later source overrides an earlier one on per-profile collision.
func ReadFlagDirs(sources ...*paths.Path) map[string][]string {
	res := map[string][]string{}
	for _, src := range sources {
		if src == nil {
			continue
		}
		files := paths.PathList{src}
		if src.IsDir() {
			var err error
			files, err = src.ReadDir(paths.FilterOutDirectories(), paths.FilterSuffixes(".conf"))
			if err != nil {
				continue
			}
		}
		for _, file := range files {
			maps.Copy(res, ReadFlagsFile(file))
		}
	}
	return res
}

// ReadConfDirs returns the concatenated filtered lines of every *.conf file
// in dirs. Directories are read in the given order, files within a directory
// in alphabetical order.
func ReadConfDirs(dirs ...*paths.Path) []string {
	var res []string
	for _, dir := range dirs {
		if dir == nil || !dir.IsDir() {
			continue
		}
		files, err := dir.ReadDir(paths.FilterOutDirectories(), paths.FilterSuffixes(".conf"))
		if err != nil {
			continue
		}
		for _, file := range files {
			res = append(res, file.MustReadFilteredFileAsLines()...)
		}
	}
	return res
}
