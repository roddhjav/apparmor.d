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
	if err := util.CopyTo(paths.New("apparmor.d/groups/_full/"), prebuild.Root.Join("apparmor.d")); err != nil {
		return res, err
	}

	// Set systemd profile name
	path := prebuild.RootApparmord.Join("tunables/multiarch.d/system")
	out, err := util.ReadFile(path)
	if err != nil {
		return res, err
	}
	out = strings.Replace(out, "@{p_systemd}=unconfined", "@{p_systemd}=systemd", -1)
	out = strings.Replace(out, "@{p_systemd_user}=unconfined", "@{p_systemd_user}=systemd-user", -1)
	if err := path.WriteFile([]byte(out)); err != nil {
		return res, err
	}

	// Fix conflicting x modifiers in abstractions - FIXME: Temporary solution
	path = prebuild.RootApparmord.Join("abstractions/gstreamer")
	out, err = util.ReadFile(path)
	if err != nil {
		return res, err
	}
	regFixConflictX := util.ToRegexRepl([]string{`.*gst-plugin-scanner.*`, ``})
	out = regFixConflictX.Replace(out)
	if err := path.WriteFile([]byte(out)); err != nil {
		return res, err
	}

	// Set systemd unit drop-in files
	return res, util.CopyTo(prebuild.SystemdDir.Join("full"), prebuild.Root.Join("systemd"))
}
