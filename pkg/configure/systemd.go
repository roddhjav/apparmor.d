// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package configure

import (
	"github.com/roddhjav/apparmor.d/pkg/paths"
	"github.com/roddhjav/apparmor.d/pkg/prebuild"
	"github.com/roddhjav/apparmor.d/pkg/tasks"
)

type SystemdDefault struct {
	tasks.Base
}

func init() {
	RegisterTask(&SystemdDefault{
		Base: tasks.Base{
			Keyword: "systemd-default",
			Msg:     "Configure systemd unit drop in files to a profile for some units",
		},
	})
}

func (p SystemdDefault) Apply() ([]string, error) {
	return []string{}, paths.CopyTo(prebuild.SystemdDir.Join("default"), prebuild.Root.Join("systemd"))
}
