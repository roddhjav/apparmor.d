// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package configure

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/roddhjav/apparmor.d/pkg/paths"
	"github.com/roddhjav/apparmor.d/pkg/prebuild"
	"github.com/roddhjav/apparmor.d/pkg/tasks"
	"github.com/roddhjav/apparmor.d/pkg/util"
)

// SystemdFSPRel is the generated full system policy systemd drop-ins
// directory, relative to the build root. It mirrors /etc/systemd.
const SystemdFSPRel = "systemd-fsp"

type SystemdDefault struct {
	tasks.BaseTask
}

type SystemdFSP struct {
	tasks.BaseTask
	dirs paths.PathList
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

// NewSystemdFSP creates a new SystemdFSP task that generates the full system
// policy systemd unit drop-ins from the fsp.d config dirs.
func NewSystemdFSP(dirs paths.PathList) *SystemdFSP {
	return &SystemdFSP{
		BaseTask: tasks.BaseTask{
			Keyword: "systemd-fsp",
			Msg:     "Configure systemd unit drop in files for full system policy",
		},
		dirs: dirs,
	}
}

// Apply writes a system/<unit>.service.d/apparmor.conf drop-in setting the
// AppArmorProfile of every unit listed in the fsp.d config files (one
// "<unit> <profile>" per line). A later file overrides an earlier one per unit.
func (p SystemdFSP) Apply() ([]string, error) {
	res := []string{}
	units := map[string]string{}
	for _, line := range util.ReadConfDirs(p.dirs...) {
		fields := strings.Fields(line)
		if len(fields) != 2 {
			return res, fmt.Errorf("invalid fsp.d entry %q, want '<unit> <profile>'", line)
		}
		units[strings.TrimSuffix(fields[0], ".service")] = fields[1]
	}

	root := p.Root.Join(SystemdFSPRel)
	if err := root.RemoveAll(); err != nil {
		return res, err
	}
	for _, unit := range slices.Sorted(maps.Keys(units)) {
		dropin := root.Join("system", unit+".service.d", "apparmor.conf")
		if err := dropin.Parent().MkdirAll(); err != nil {
			return res, err
		}
		content := "[Service]\nAppArmorProfile=" + units[unit] + "\n"
		if err := dropin.WriteFile([]byte(content)); err != nil {
			return res, err
		}
		res = append(res, unit)
	}
	return res, nil
}
