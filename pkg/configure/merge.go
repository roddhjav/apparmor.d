// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package configure

import (
	"os"
	"path/filepath"

	"github.com/roddhjav/apparmor.d/pkg/paths"
	"github.com/roddhjav/apparmor.d/pkg/tasks"
)

type Merge struct {
	tasks.BaseTask
}

// NewMerge creates a new Merge task.
func NewMerge() *Merge {
	return &Merge{
		BaseTask: tasks.BaseTask{
			Keyword: "merge",
			Msg:     "Merge profiles (from group/, profiles-*-*/) to a unified apparmor.d directory",
		},
	}
}

func (p Merge) Apply() ([]string, error) {
	res := []string{}
	dirToMerge := []string{
		"groups/*/*", "groups",
		"profiles-*-*/*", "profiles-*",
	}

	idx := 0
	for idx < len(dirToMerge)-1 {
		dirMoved, dirRemoved := dirToMerge[idx], dirToMerge[idx+1]
		files, err := filepath.Glob(p.RootApparmor.Join(dirMoved).String())
		if err != nil {
			return res, err
		}
		for _, file := range files {
			err := os.Rename(file, p.RootApparmor.Join(filepath.Base(file)).String())
			if err != nil {
				return res, err
			}
		}

		files, err = filepath.Glob(p.RootApparmor.Join(dirRemoved).String())
		if err != nil {
			return []string{}, err
		}
		for _, file := range files {
			if err := paths.New(file).RemoveAll(); err != nil {
				return res, err
			}
		}
		idx = idx + 2
	}

	// Namespaces directory
	nsRoot := p.RootApparmor.Join("namespaces")
	if !nsRoot.Exist() {
		return res, nil
	}
	dirs, err := nsRoot.ReadDir(paths.FilterDirectories())
	if err != nil {
		return res, err
	}
	for _, dir := range dirs {
		nsName := dir.Base()
		files, err := dir.ReadDir(paths.FilterOutDirectories())
		if err != nil {
			return res, err
		}
		for _, file := range files {
			destPath := p.RootApparmor.Join(":" + nsName + ":" + file.Base())
			err := os.Rename(file.String(), destPath.String())
			if err != nil {
				return res, err
			}
		}
	}
	if err := nsRoot.RemoveAll(); err != nil {
		return res, err
	}
	return res, nil
}
