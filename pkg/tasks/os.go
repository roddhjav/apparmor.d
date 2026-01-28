// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package tasks

import (
	"os"
	"slices"
	"strings"

	"github.com/roddhjav/apparmor.d/pkg/paths"
)

var (
	Distribution = getDistribution()
	Release      = getOSRelease()
	Family       = getFamily()

	// DebianDir is the directory where the debian specific files are stored
	// DebianDir *paths.Path = paths.New("debian")

	// // DebianHide is the path to the debian/common.hide file
	// DebianHide = DebianHider{path: DebianDir.Join("common.hide")}

	// Ignore = Ignorer{}
	// Flags  = Flagger{}
)

var (
	osReleaseFile  = "/etc/os-release"
	supportedDists = map[string][]string{
		"arch":     {},
		"debian":   {},
		"ubuntu":   {"neon"},
		"opensuse": {"suse", "opensuse-tumbleweed"},
		"whonix":   {},
	}
	famillyDists = map[string][]string{
		"apt":    {"debian", "ubuntu", "whonix"},
		"pacman": {"arch"},
		"zypper": {"opensuse"},
	}
)

func getOSRelease() map[string]string {
	var lines []string
	var err error
	for _, name := range []string{osReleaseFile, "/usr/lib/os-release"} {
		path := paths.New(name)
		if path.Exist() {
			lines, err = path.ReadFileAsLines()
			if err != nil {
				panic(err)
			}
			break
		}
	}
	os := map[string]string{}
	for _, line := range lines {
		item := strings.Split(line, "=")
		if len(item) == 2 {
			os[item[0]] = strings.Trim(item[1], "\"")
		}
	}
	return os
}

func getDistribution() string {
	dist, present := os.LookupEnv("DISTRIBUTION")
	if present {
		return dist
	}

	id := Release["ID"]
	if id == "ubuntu" {
		return id
	}
	idLike := Release["ID_LIKE"]
	for main, based := range supportedDists {
		if main == id || main == idLike {
			return main
		} else if slices.Contains(based, id) {
			return main
		} else if slices.Contains(based, idLike) {
			return main
		}
	}
	return id
}

func getFamily() string {
	for familly, dist := range famillyDists {
		if slices.Contains(dist, Distribution) {
			return familly
		}
	}
	return ""
}
