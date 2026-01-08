// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package prepare

import (
	"github.com/roddhjav/apparmor.d/pkg/paths"
	"github.com/roddhjav/apparmor.d/pkg/prebuild"
)

type SystemdDefault struct {
	prebuild.Base
}

func init() {
	RegisterTask(&SystemdDefault{
		Base: prebuild.Base{
			Keyword: "systemd-default",
			Msg:     "Configure systemd unit drop in files to a profile for some units",
		},
	})
}

func (p SystemdDefault) Apply() ([]string, error) {
	return []string{}, paths.CopyTo(prebuild.SystemdDir.Join("default"), prebuild.Root.Join("systemd"))
}
