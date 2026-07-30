// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package configure

import (
	"fmt"

	"github.com/roddhjav/apparmor.d/pkg/paths"
	"github.com/roddhjav/apparmor.d/pkg/tasks"
	"github.com/roddhjav/apparmor.d/pkg/util"
)

type SetFlags struct {
	tasks.BaseTask
	sources paths.PathList
}

// NewSetFlags creates a SetFlags task that applies flags from the given
// sources, each a flags file or a directory of *.conf files, read in order
// so a later source overrides an earlier one. aa-install passes the
// flags.d config dirs; prebuild can pass the dists/flags.d files of the
// target distribution.
func NewSetFlags(sources paths.PathList) *SetFlags {
	return &SetFlags{
		BaseTask: tasks.BaseTask{
			Keyword: "setflags",
			Msg:     "Set flags as defined in:",
		},
		sources: sources,
	}
}

func (p SetFlags) Apply() ([]string, error) {
	entries := util.ReadFlagDirs(p.sources...)
	if len(entries) == 0 {
		return nil, nil
	}
	res, err := p.applyFlags(entries)
	if err != nil {
		return res, err
	}
	return append(res, existingPaths(p.sources)...), nil
}

func (p SetFlags) applyFlags(entries map[string][]string) ([]string, error) {
	var res []string
	for profile, flags := range entries {
		file := p.RootApparmor.Join(profile)
		if !file.Exist() {
			if renamed := p.RootApparmor.Join(profile + "." + p.Pkgname); renamed.Exist() {
				file = renamed
			} else {
				res = append(res, fmt.Sprintf("Profile %s not found, ignoring", profile))
				continue
			}
		}
		if len(flags) == 0 {
			continue
		}
		out, err := file.ReadFileAsString()
		if err != nil {
			return res, err
		}
		out, err = util.ApplyFlags(out, flags)
		if err != nil {
			return res, err
		}
		if err := file.WriteFile([]byte(out)); err != nil {
			return res, err
		}
	}
	return res, nil
}
