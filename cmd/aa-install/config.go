// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package main

import (
	"fmt"
	"slices"
	"strings"

	"github.com/roddhjav/apparmor.d/pkg/paths"
	"github.com/roddhjav/apparmor.d/pkg/util"
)

var (
	// deployModes are the modes accepted as a default deploy mode.
	deployModes = []string{"enforce", "complain"}

	// includeModes are the modes accepted for the include key.
	includeModes = []string{"default", "full", "all"}

	// reloadModes are the modes accepted for the reload key.
	reloadModes = []string{"yes", "no"}

	// vendorConfigDir holds the vendor configuration defaults
	vendorConfigDir = paths.New("/usr/share/apparmor")
)

type conf struct {
	mode          string
	include       string
	reload        bool
	flagDirs      paths.PathList
	ignoreDirs    paths.PathList
	includeDirs   paths.PathList
	overwriteDirs paths.PathList
}

// configTier returns the drop-in directories for name, vendor tier first so
// entries in configDir override it.
func configTier(configDir *paths.Path, name string) paths.PathList {
	return paths.PathList{vendorConfigDir.Join(name), configDir.Join(name)}
}

// loadConfig resolves the configuration from vendorConfigDir and configDir.
// The deploy mode is the "default" key of the modes files, overridden by
// cli arguments.
func loadConfig(configDir *paths.Path) (*conf, error) {
	res := &conf{
		flagDirs:      configTier(configDir, "flags.d"),
		ignoreDirs:    configTier(configDir, "ignore.d"),
		includeDirs:   configTier(configDir, "include.d"),
		overwriteDirs: configTier(configDir, "overwrite.d"),
	}

	modes := readModeConfig(vendorConfigDir.Join("modes"), configDir.Join("modes"))
	res.mode = modes["default"]
	if res.mode == "" {
		res.mode = "complain"
	}
	switch {
	case complain:
		res.mode = "complain"
	case enforce:
		res.mode = "enforce"
	}
	if !slices.Contains(deployModes, res.mode) {
		return nil, fmt.Errorf("invalid default mode %q in %s", res.mode, configDir.Join("modes"))
	}
	res.include = modes["include"]
	if res.include == "" {
		res.include = "default"
	}
	if all {
		res.include = "all"
	}
	if !slices.Contains(includeModes, res.include) {
		return nil, fmt.Errorf("invalid include mode %q in %s", res.include, configDir.Join("modes"))
	}
	reload := modes["reload"]
	if reload == "" {
		reload = "yes"
	}
	if !slices.Contains(reloadModes, reload) {
		return nil, fmt.Errorf("invalid reload mode %q in %s", reload, configDir.Join("modes"))
	}
	res.reload = reload == "yes" && !noReload
	return res, nil
}

// readModeConfig parses a general aa-install config file (one "key value"
// per line, # comments filtered). Only the last existing file is read: a
// later file fully replaces an earlier one.
func readModeConfig(files ...*paths.Path) map[string]string {
	res := map[string]string{}
	for i := len(files) - 1; i >= 0; i-- {
		if !files[i].Exist() {
			continue
		}
		for _, line := range files[i].MustReadFilteredFileAsLines() {
			if key, value, ok := strings.Cut(line, " "); ok {
				res[strings.TrimSpace(key)] = strings.TrimSpace(value)
			}
		}
		break
	}
	return res
}

// userModeOverrides returns the set of profiles a flags.d drop-in assigns a
// mode to; the default deploy mode must not override these.
func userModeOverrides(dirs paths.PathList) map[string]bool {
	res := map[string]bool{}
	for profile, flags := range util.ReadFlagDirs(dirs...) {
		for _, f := range flags {
			if slices.Contains(util.ProfileModes, f) {
				res[profile] = true
				break
			}
		}
	}
	return res
}
