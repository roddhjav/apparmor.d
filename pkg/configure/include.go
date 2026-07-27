// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package configure

import (
	"github.com/roddhjav/apparmor.d/pkg/paths"
	"github.com/roddhjav/apparmor.d/pkg/tasks"
	"github.com/roddhjav/apparmor.d/pkg/util"
)

type Include struct {
	tasks.BaseTask
	userDirs paths.PathList
}

// NewInclude creates an Include task that only keeps the profiles listed in
// the *.conf files of dirs, read in order. Intended for the aa-install
// pipeline.
func NewInclude(dirs paths.PathList) *Include {
	return &Include{
		BaseTask: tasks.BaseTask{
			Keyword: "include",
			Msg:     "Only keep profiles and groups listed in manifest include files",
		},
		userDirs: dirs,
	}
}

// Active reports whether any user include rules are defined. When true,
// this task takes over the selection role otherwise performed by Install.
func (p Include) Active() bool {
	return len(util.ReadConfDirs(p.userDirs...)) > 0
}

func (p Include) Apply() ([]string, error) {
	res := []string{}
	entries := util.ReadConfDirs(p.userDirs...)
	if len(entries) == 0 {
		return res, nil
	}

	keep := map[string]bool{}
	var keepDirs []*paths.Path
	for _, entry := range entries {
		dir := p.RootApparmor.Join(entry)
		if dir.IsDir() {
			keepDirs = append(keepDirs, dir)
			continue
		}
		keep[entry] = true
	}

	files, err := p.RootApparmor.ReadDirRecursiveFiltered(skipSystemDirs, paths.FilterOutDirectories())
	if err != nil {
		return res, err
	}
	for _, file := range files {
		if keep[file.Base()] {
			continue
		}
		if file.IsInsideAnyDir(keepDirs) {
			continue
		}
		if err := file.Remove(); err != nil {
			return res, err
		}
	}

	res = append(res, existingPaths(p.userDirs)...)
	return res, nil
}
