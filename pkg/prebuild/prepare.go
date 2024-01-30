// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2023 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package prebuild

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/arduino/go-paths-helper"
	"github.com/roddhjav/apparmor.d/pkg/logging"
)

// Prepare the build directory with the following tasks
var (
	Prepares = []PrepareFunc{
		Synchronise,
		Ignore,
		Merge,
		Configure,
		SetFlags,
	}
	PrepareMsg = map[string]string{
		"Synchronise":         "Initialize a new clean apparmor.d build directory",
		"Ignore":              "Ignore profiles and files from:",
		"Merge":               "Merge all profiles",
		"Configure":           "Set distribution specificities",
		"SetFlags":            "Set flags on some profiles",
		"SetDefaultSystemd":   "Set systemd unit drop in files to ensure some service start after apparmor",
		"SetFullSystemPolicy": "Configure AppArmor for full system policy",
	}
)

type PrepareFunc func() ([]string, error)

// Initialize a new clean apparmor.d build directory
func Synchronise() ([]string, error) {
	res := []string{}
	dirs := paths.PathList{RootApparmord, Root.Join("root"), Root.Join("systemd")}
	for _, dir := range dirs {
		if err := dir.RemoveAll(); err != nil {
			return res, err
		}
	}
	for _, name := range []string{"apparmor.d", "root"} {
		if err := copyTo(paths.New(name), Root.Join(name)); err != nil {
			return res, err
		}
	}
	return res, nil
}

// Ignore profiles and files as defined in dists/ignore/
func Ignore() ([]string, error) {
	res := []string{}
	for _, name := range []string{"main.ignore", Distribution + ".ignore"} {
		path := DistDir.Join("ignore", name)
		if !path.Exist() {
			continue
		}
		lines, _ := path.ReadFileAsLines()
		for _, line := range lines {
			if strings.HasPrefix(line, "#") || line == "" {
				continue
			}
			profile := Root.Join(line)
			if profile.NotExist() {
				files, err := RootApparmord.ReadDirRecursiveFiltered(nil, paths.FilterNames(line))
				if err != nil {
					return res, err
				}
				for _, path := range files {
					if err := path.RemoveAll(); err != nil {
						return res, err
					}
				}
			} else {
				if err := profile.RemoveAll(); err != nil {
					return res, err
				}
			}
		}
		res = append(res, path.String())
	}
	return res, nil
}

// Merge all profiles in a new apparmor.d directory
func Merge() ([]string, error) {
	res := []string{}
	dirToMerge := []string{
		"groups/*/*", "groups",
		"profiles-*-*/*", "profiles-*",
	}

	idx := 0
	for idx < len(dirToMerge)-1 {
		dirMoved, dirRemoved := dirToMerge[idx], dirToMerge[idx+1]
		files, err := filepath.Glob(RootApparmord.Join(dirMoved).String())
		if err != nil {
			return res, err
		}
		for _, file := range files {
			err := os.Rename(file, RootApparmord.Join(filepath.Base(file)).String())
			if err != nil {
				return res, err
			}
		}

		files, err = filepath.Glob(RootApparmord.Join(dirRemoved).String())
		if err != nil {
			return []string{}, err
		}
		for _, file := range files {
			if err := paths.New(file).RemoveAll(); err != nil {
				return res, err
			}
		}
		idx = idx + 2
	}
	return res, nil
}

// Set the distribution specificities
func Configure() ([]string, error) {
	res := []string{}
	switch Distribution {
	case "arch", "opensuse":

	case "debian", "ubuntu", "whonix":
		// Copy Ubuntu specific profiles
		if err := copyTo(DistDir.Join("ubuntu"), RootApparmord); err != nil {
			return res, err
		}

	default:
		return []string{}, fmt.Errorf("%s is not a supported distribution", Distribution)

	}
	return res, nil
}

// Set flags on some profiles according to manifest defined in `dists/flags/`
func SetFlags() ([]string, error) {
	res := []string{}
	for _, name := range []string{"main.flags", Distribution + ".flags"} {
		path := FlagDir.Join(name)
		if !path.Exist() {
			continue
		}
		lines, _ := path.ReadFileAsLines()
		for _, line := range lines {
			if strings.HasPrefix(line, "#") || line == "" {
				continue
			}
			manifest := strings.Split(line, " ")
			profile := manifest[0]
			file := RootApparmord.Join(profile)
			if !file.Exist() {
				logging.Warning("Profile %s not found", profile)
				continue
			}

			// If flags is set, overwrite profile flag
			if len(manifest) > 1 {
				flags := " flags=(" + manifest[1] + ") {"
				content, err := file.ReadFile()
				if err != nil {
					return res, err
				}

				// Remove all flags definition, then set manifest' flags
				out := regFlags.ReplaceAllLiteralString(string(content), "")
				out = regProfileHeader.ReplaceAllLiteralString(out, flags)
				if err := file.WriteFile([]byte(out)); err != nil {
					return res, err
				}
			}
		}
		res = append(res, path.String())
	}
	return res, nil
}

// Set systemd unit drop in files to ensure some service start after apparmor
func SetDefaultSystemd() ([]string, error) {
	return []string{}, copyTo(paths.New("systemd/default/"), Root.Join("systemd"))
}

// Set AppArmor for (experimental) full system policy.
// See https://apparmor.pujol.io/full-system-policy/
func SetFullSystemPolicy() ([]string, error) {
	res := []string{}
	// Install full system policy profiles
	if err := copyTo(paths.New("apparmor.d/groups/_full/"), Root.Join("apparmor.d")); err != nil {
		return res, err
	}

	// Set systemd profile name
	path := RootApparmord.Join("tunables/multiarch.d/system")
	content, err := path.ReadFile()
	if err != nil {
		return res, err
	}
	out := strings.Replace(string(content), "@{systemd}=unconfined", "@{systemd}=systemd", -1)
	if err := path.WriteFile([]byte(out)); err != nil {
		return res, err
	}

	// Set systemd unit drop-in files
	return res, copyTo(paths.New("systemd/full/"), Root.Join("systemd"))
}
