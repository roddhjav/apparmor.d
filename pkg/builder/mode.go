// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package builder

import (
	"github.com/roddhjav/apparmor.d/pkg/tasks"
	"github.com/roddhjav/apparmor.d/pkg/util"
)

// DeployMode sets the default deploy mode (enforce/complain) on every profile,
// except those a user flags.d drop-in already assigns a mode to, those win.
type DeployMode struct {
	tasks.BaseTask
	mode       string
	overridden map[string]bool
}

// NewDeployMode creates a new DeployMode builder. An empty mode is a no-op
// (profiles keep the mode they were built with).
func NewDeployMode(mode string, overridden map[string]bool) *DeployMode {
	return &DeployMode{
		BaseTask: tasks.BaseTask{
			Keyword: "deploy-mode",
			Msg:     "Build: set the default deploy mode",
		},
		mode:       mode,
		overridden: overridden,
	}
}

func (b DeployMode) Apply(opt *Option, profile string) (string, error) {
	if b.mode == "" {
		return profile, nil
	}
	name := opt.Name
	if m := regProfileName.FindStringSubmatch(profile); m != nil {
		name = m[1]
	}
	if b.overridden[name] {
		return profile, nil // a user flags.d drop-in sets this profile's mode
	}
	return util.SetMode(profile, b.mode)
}
