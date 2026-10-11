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

// NewSystemdDefault creates a new SystemdDefault task.
func NewSystemdDefault() *SystemdDefault {
	return &SystemdDefault{
		BaseTask: tasks.BaseTask{
			Keyword: "systemd-default",
			Msg:     "Configure systemd unit drop in files to a profile for some units",
		},
	}
}

func (p SystemdDefault) Apply() ([]string, error) {
	// Regenerate the systemd drop-in dir from scratch: p.Root is reused
	// across builds (only apparmor.d/ and share/ are re-synchronised), so a
	// leftover full-system-policy drop-in from a prior --fsp run would
	// otherwise survive into a non-fsp build. This runs before the FSP task,
	// which layers its own drop-ins on top.
	dst := p.Root.Join("systemd")
	if err := dst.RemoveAll(); err != nil {
		return []string{}, err
	}
	return []string{}, paths.CopyTo(prebuild.SystemdDir.Join("default"), dst)
}
