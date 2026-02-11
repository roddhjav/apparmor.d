// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package builder

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/roddhjav/apparmor.d/pkg/prebuild"
	"github.com/roddhjav/apparmor.d/pkg/tasks"
)

var (
	regProfileName = regexp.MustCompile(`(?m)^profile\s+(\S+)\s+`)
	profileModes   = []string{
		"enforce", "complain", "kill", "default_allow", "unconfined", "prompt",
	}
)

type ProfileMode struct {
	tasks.BaseTask
	modes map[string]string
}

// NewProfileMode creates a new ProfileMode builder.
func NewProfileMode() *ProfileMode {
	modes := make(map[string]string)
	for _, name := range []string{"main", tasks.Distribution} {
		for profile, flags := range prebuild.Flags.Read(name) {
			if len(flags) > 0 {
				modes[profile] = flags[0]
			}
		}
	}
	return &ProfileMode{
		BaseTask: tasks.BaseTask{
			Keyword: "profile-mode",
			Msg:     "Build: set modes (complain, enforce...) as definied in dist/flags",
		},
		modes: modes,
	}
}

func (b ProfileMode) Apply(opt *Option, profile string) (string, error) {
	matches := regProfileName.FindStringSubmatch(profile)
	if matches == nil {
		return profile, nil
	}

	name := matches[1]
	mode, present := b.modes[name]
	if !present {
		return profile, nil
	}
	if !slices.Contains(profileModes, mode) {
		return profile, fmt.Errorf("unknown profile mode: %s", mode)
	}

	return setMode(profile, mode)
}

func setMode(profile string, mode string) (string, error) {
	flags := []string{}
	matches := regFlags.FindStringSubmatch(profile)
	if len(matches) != 0 {
		flags = strings.Split(matches[1], ",")
	}

	// Remove all conflicting mode flags
	flags = slices.DeleteFunc(flags, func(f string) bool {
		return slices.Contains(profileModes, f)
	})

	// "enforce" is the default (no mode flag needed), otherwise add the mode
	if mode != "enforce" {
		flags = append(flags, mode)
	}

	// Remove all flags definition, then set the new flags
	profile = regFlags.ReplaceAllLiteralString(profile, "")
	if len(flags) > 0 {
		flagsStr := " flags=(" + strings.Join(flags, ",") + ") {\n"
		profile = regProfileHeader.ReplaceAllLiteralString(profile, flagsStr)
	}
	return profile, nil
}
