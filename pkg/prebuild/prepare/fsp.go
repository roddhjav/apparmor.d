// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package prepare

import (
	"strings"

	"github.com/arduino/go-paths-helper"
	"github.com/roddhjav/apparmor.d/pkg/prebuild/cfg"
	"github.com/roddhjav/apparmor.d/pkg/util"
)

type FullSystemPolicy struct {
	cfg.Base
}

func init() {
	RegisterTask(&FullSystemPolicy{
		Base: cfg.Base{
			Keyword: "fsp",
			Msg:     "Configure AppArmor for full system policy",
		},
	})
}

func (p FullSystemPolicy) Apply() ([]string, error) {
	res := []string{}

	// Install full system policy profiles
	if err := util.CopyTo(paths.New("apparmor.d/groups/_full/"), cfg.Root.Join("apparmor.d")); err != nil {
		return res, err
	}

	// Set systemd profile name
	path := cfg.RootApparmord.Join("tunables/multiarch.d/system")
	content, err := path.ReadFile()
	if err != nil {
		return res, err
	}
	out := strings.Replace(string(content), "@{systemd}=unconfined", "@{systemd}=systemd", -1)
	out = strings.Replace(out, "@{systemd_user}=unconfined", "@{systemd_user}=systemd-user", -1)
	if err := path.WriteFile([]byte(out)); err != nil {
		return res, err
	}

	// Fix conflicting x modifiers in abstractions - FIXME: Temporary solution
	path = cfg.RootApparmord.Join("abstractions/gstreamer")
	content, err = path.ReadFile()
	if err != nil {
		return res, err
	}
	out = string(content)
	regFixConflictX := util.ToRegexRepl([]string{`.*gst-plugin-scanner.*`, ``})
	out = regFixConflictX.Replace(out)
	if err := path.WriteFile([]byte(out)); err != nil {
		return res, err
	}

	// Set systemd unit drop-in files
	return res, util.CopyTo(cfg.SystemdDir.Join("full"), cfg.Root.Join("systemd"))
}
