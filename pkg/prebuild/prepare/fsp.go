// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package prepare

import (
	"strings"

	"github.com/roddhjav/apparmor.d/pkg/paths"
	"github.com/roddhjav/apparmor.d/pkg/prebuild"
	"github.com/roddhjav/apparmor.d/pkg/util"
)

type FullSystemPolicy struct {
	prebuild.Base
}

func init() {
	RegisterTask(&FullSystemPolicy{
		Base: prebuild.Base{
			Keyword: "fsp",
			Msg:     "Configure AppArmor for full system policy",
		},
	})
}

func (p FullSystemPolicy) Apply() ([]string, error) {
	res := []string{}

	// Install full system policy profiles
	if err := paths.New("apparmor.d/groups/_full/").CopyFS(prebuild.Root.Join("apparmor.d")); err != nil {
		return res, err
	}

	// Set systemd profile name
	path := prebuild.RootApparmord.Join("tunables/multiarch.d/profiles")
	out, err := path.ReadFileAsString()
	if err != nil {
		return res, err
	}
	out = strings.ReplaceAll(out, "@{p_systemd}=unconfined", "@{p_systemd}=systemd")
	out = strings.ReplaceAll(out, "@{p_systemd_user}=unconfined", "@{p_systemd_user}=systemd-user")
	if err := path.WriteFile([]byte(out)); err != nil {
		return res, err
	}

	// Fix conflicting x modifiers in abstractions - FIXME: Temporary solution
	path = prebuild.RootApparmord.Join("abstractions/gstreamer")
	out, err = path.ReadFileAsString()
	if err != nil {
		return res, err
	}
	regFixConflictX := util.ToRegexRepl([]string{`.*gst-plugin-scanner.*`, ``})
	out = regFixConflictX.Replace(out)
	if err := path.WriteFile([]byte(out)); err != nil {
		return res, err
	}

	// Set systemd unit drop-in files
	return res, prebuild.SystemdDir.Join("full").CopyFS(prebuild.Root.Join("systemd"))
}
