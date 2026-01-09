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
	tasks.BaseTask
}

// NewSystemd creates a new SystemdDefault task.
func NewSystemd() *SystemdDefault {
	return &SystemdDefault{
		BaseTask: tasks.BaseTask{
			Keyword: "systemd-default",
			Msg:     "Configure systemd unit drop in files to a profile for some units",
		},
	}
}

func (p SystemdDefault) Apply() ([]string, error) {
	return []string{}, paths.CopyTo(prebuild.SystemdDir.Join("default"), p.Root.Join("systemd"))
}
