// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package builder

import (
	"regexp"

	"github.com/roddhjav/apparmor.d/pkg/prebuild"
	"github.com/roddhjav/apparmor.d/pkg/tasks"
	"github.com/roddhjav/apparmor.d/pkg/util"
)

var (
	regProfileName = regexp.MustCompile(`(?m)^profile\s+(\S+)\s+`)
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

	return util.SetMode(profile, mode)
}
