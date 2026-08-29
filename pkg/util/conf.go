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

// EffectiveConfFiles returns the *.conf files of the given files or
// directories, in order. A same-named file from a later source fully
// replaces an earlier one (admin config overrides vendor config).
func EffectiveConfFiles(sources ...*paths.Path) paths.PathList {
	var order []string
	byName := map[string]*paths.Path{}
	add := func(file *paths.Path) {
		if _, seen := byName[file.Base()]; !seen {
			order = append(order, file.Base())
		}
		byName[file.Base()] = file
	}
	for _, src := range sources {
		if src == nil {
			continue
		}
		if src.IsDir() {
			files, err := src.ReadDir(paths.FilterOutDirectories(), paths.FilterSuffixes(".conf"))
			if err != nil {
				continue
			}
			for _, file := range files {
				add(file)
			}
		} else if src.Exist() {
			add(src)
		}
	}
	res := make(paths.PathList, 0, len(order))
	for _, name := range order {
		res = append(res, byName[name])
	}
	return res
}

// ReadFlagDirs merges the flags manifests from the effective files of
// sources; a later file overrides an earlier one per profile.
func ReadFlagDirs(sources ...*paths.Path) map[string][]string {
	res := map[string][]string{}
	for _, file := range EffectiveConfFiles(sources...) {
		maps.Copy(res, ReadFlagsFile(file))
	}
	return res
}

// ReadConfDirs returns the filtered lines of the effective *.conf files in dirs.
func ReadConfDirs(dirs ...*paths.Path) []string {
	var res []string
	for _, file := range EffectiveConfFiles(dirs...) {
		res = append(res, file.MustReadFilteredFileAsLines()...)
	}
	return res
}
