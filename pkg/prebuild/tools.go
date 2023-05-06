// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2023 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package prebuild

import (
	"os"
	"strings"

	"github.com/arduino/go-paths-helper"
	"github.com/roddhjav/apparmor.d/pkg/aa"
	"golang.org/x/exp/slices"
)

var (
	osReleaseFile   = "/etc/os-release"
	firstPartyDists = []string{"arch", "debian", "ubuntu", "opensuse", "whonix"}
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

	if slices.Contains(firstPartyDists, id) {
		return id
	} else if slices.Contains(firstPartyDists, id_like) {
		return id_like
	}
	return id
}

func setLibexec(libexec string) error {
	aa.Tunables["libexec"] = []string{libexec}
	file, err := RootApparmord.Join("tunables", "multiarch.d", "apparmor.d").Append()
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(`@{libexec}=` + libexec)
	return err
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
