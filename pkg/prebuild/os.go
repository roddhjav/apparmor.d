// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package prebuild

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
	id_like := Release["ID_LIKE"]
	for main, based := range supportedDists {
		if main == id || main == id_like {
			return main
		} else if slices.Contains(based, id) {
			return main
		} else if slices.Contains(based, id_like) {
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
