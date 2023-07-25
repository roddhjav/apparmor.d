// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2023 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package prebuild

import (
	"os"
	"strings"

	"github.com/arduino/go-paths-helper"
	"golang.org/x/exp/slices"
)

var (
	osReleaseFile   = "/etc/os-release"
	firstPartyDists = []string{"arch", "debian", "ubuntu", "opensuse", "whonix"}
	supportedDists  = map[string][]string{
		"arch":     {},
		"debian":   {},
		"ubuntu":   {},
		"opensuse": {"suse"},
		"whonix":   {},
	}
)

func getSupportedDistribution() string {
	dist, present := os.LookupEnv("DISTRIBUTION")
	if present {
		return dist
	}

	lines, err := paths.New(osReleaseFile).ReadFileAsLines()
	if err != nil {
		panic(err)
	}

	id := ""
	id_like := ""
	for _, line := range lines {
		item := strings.Split(line, "=")
		if item[0] == "ID" {
			id = strings.Split(strings.Trim(item[1], "\""), " ")[0]
		} else if item[0] == "ID_LIKE" {
			id_like = strings.Split(strings.Trim(item[1], "\""), " ")[0]
		}
	}

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

func copyTo(src *paths.Path, dst *paths.Path) error {
	files, err := src.ReadDirRecursiveFiltered(nil, paths.FilterOutDirectories())
	if err != nil {
		return err
	}
	for _, file := range files {
		destination, err := file.RelFrom(src)
		if err != nil {
			return err
		}
		destination = dst.JoinPath(destination)
		if err := file.CopyTo(destination); err != nil {
			return err
		}
	}
	return nil
}
