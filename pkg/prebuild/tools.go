// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2023-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package prebuild

import (
	"os"
	"strings"

	"github.com/arduino/go-paths-helper"
	"golang.org/x/exp/slices"
)

var (
	osReleaseFile  = "/etc/os-release"
	supportedDists = map[string][]string{
		"arch":     {},
		"debian":   {},
		"ubuntu":   {},
		"opensuse": {"suse"},
		"whonix":   {},
	}

func NewOSRelease() map[string]string {
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

func getSupportedDistribution() string {
	dist, present := os.LookupEnv("DISTRIBUTION")
	if present {
		return dist
	}

	os := NewOSRelease()
	id := os["ID"]
	if id == "ubuntu" {
		return id
	}
	id_like := os["ID_LIKE"]
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
	files, err := src.ReadDirRecursiveFiltered(nil, paths.FilterOutDirectories(), paths.FilterOutNames("README.md"))
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

// Displace files in the package sources
func displaceFiles(files []string) error {
	const ext = ".apparmor.d"
	for _, name := range files {
		origin := RootApparmord.Join(name)
		dest := RootApparmord.Join(name + ext)
		if err := origin.Rename(dest); err != nil {
			return err
		}
		file, err := paths.New("debian/apparmor.d.displace").Append()
		if err != nil {
			return err
		}
		if _, err := file.WriteString("/etc/apparmor.d/" + name + ext + "\n"); err != nil {
			return err
		}
	}
	return nil
}

func overwriteProfile(path *paths.Path) []string {
	res := []string{}
	lines, err := path.ReadFileAsLines()
	if err != nil {
		panic(err)
	}
	for _, line := range lines {
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}
		res = append(res, line)
	}
	return res
}

