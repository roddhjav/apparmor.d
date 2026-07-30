// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package configure

import (
	"github.com/roddhjav/apparmor.d/pkg/paths"
	"github.com/roddhjav/apparmor.d/pkg/prebuild"
	"github.com/roddhjav/apparmor.d/pkg/tasks"
	"github.com/roddhjav/apparmor.d/pkg/util"
)

type Ignore struct {
	tasks.BaseTask
	userOnly bool
	userDirs paths.PathList
}

// NewIgnore creates an Ignore task that removes profiles listed in the
// distribution ignore files (dists/ignore.d/). Intended for prebuild.
func NewIgnore() *Ignore {
	return &Ignore{
		BaseTask: tasks.BaseTask{
			Keyword: "ignore",
			Msg:     "Ignore profiles and files from:",
		},
	}
}

// NewUserIgnore creates an Ignore task that removes profiles listed in every
// *.conf file of dirs, read in order. Intended for the aa-install pipeline.
func NewUserIgnore(dirs paths.PathList) *Ignore {
	return &Ignore{
		BaseTask: tasks.BaseTask{
			Keyword: "ignore",
			Msg:     "Ignore profiles and files from:",
		},
		userOnly: true,
		userDirs: dirs,
	}
}

func (p Ignore) Apply() ([]string, error) {
	res := []string{}
	if p.userOnly {
		user := util.ReadConfDirs(p.userDirs...)
		if len(user) == 0 {
			return res, nil
		}
		if err := p.removeEntries(user); err != nil {
			return res, err
		}
		return append(res, existingPaths(p.userDirs)...), nil
	}
	for _, name := range []string{"main", tasks.Distribution} {
		if err := p.removeEntries(prebuild.Ignore.Read(name)); err != nil {
			return res, err
		}
		res = append(res, prebuild.IgnoreDir.Join(name+".conf").String())
	}
	return res, nil
}

func (p Ignore) removeEntries(entries []string) error {
	for _, ignore := range entries {
		profile := p.Root.Join(ignore)
		if profile.NotExist() {
			files, err := p.RootApparmor.ReadDirRecursiveFiltered(nil, paths.FilterNames(ignore))
			if err != nil {
				return err
			}
			for _, path := range files {
				if err := path.RemoveAll(); err != nil {
					return err
				}
			}
		} else {
			if err := profile.RemoveAll(); err != nil {
				return err
			}
		}
	}
	return nil
}
