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
	src      *paths.Path // restore mode: source of the re-applied profiles
}

// NewInclude creates an Include task that only keeps the profiles listed in
// the *.conf files of dirs, plus the alwaysKeep ones (modes: include full).
func NewInclude(dirs paths.PathList) *Include {
	return &Include{
		BaseTask: tasks.BaseTask{
			Keyword: "include",
			Msg:     "Only keep profiles and groups listed in manifest include files",
		},
		userDirs: dirs,
	}
}

// NewRestoreInclude creates an Include task that re-applies the listed
// profiles from src after the ignore task (modes: include default).
func NewRestoreInclude(dirs paths.PathList, src *paths.Path) *Include {
	return &Include{
		BaseTask: tasks.BaseTask{
			Keyword: "include",
			Msg:     "Re-apply ignored profiles listed in manifest include files",
		},
		userDirs: dirs,
		src:      src,
	}
}

// Active reports whether any user include rules are defined. When true,
// this task takes over the selection role otherwise performed by Install.
func (p Include) Active() bool {
	return len(util.ReadConfDirs(p.userDirs...)) > 0
}

func (p Include) Apply() ([]string, error) {
	entries := util.ReadConfDirs(p.userDirs...)
	if len(entries) == 0 {
		return []string{}, nil
	}
	if p.src != nil {
		return p.restore(entries)
	}
	return p.only(entries)
}

// only removes every profile not listed in entries.
func (p Include) only(entries []string) ([]string, error) {
	res := []string{}
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
		if keep[file.Base()] || isAlwaysKept(file.Base()) {
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

// restore copies the listed entries back from the source directory.
func (p Include) restore(entries []string) ([]string, error) {
	res := []string{}
	for _, entry := range entries {
		var files paths.PathList
		var err error
		if src := p.src.Join(entry); src.IsDir() {
			files, err = src.ReadDirRecursiveFiltered(nil, paths.FilterOutDirectories())
		} else if src.Exist() {
			files = paths.PathList{src}
		} else {
			files, err = p.src.ReadDirRecursiveFiltered(skipSystemDirs, paths.FilterNames(entry))
		}
		if err != nil {
			return res, err
		}
		for _, file := range files {
			if err := p.restoreFile(file); err != nil {
				return res, err
			}
		}
	}
	return append(res, existingPaths(p.userDirs)...), nil
}

// restoreFile copies a source file to its build location when missing.
func (p Include) restoreFile(file *paths.Path) error {
	rel, err := file.RelFrom(p.src)
	if err != nil {
		return err
	}
	dst := p.RootApparmor.JoinPath(rel)
	if dst.Exist() {
		return nil
	}
	if err := dst.Parent().MkdirAll(); err != nil {
		return err
	}
	return file.CopyTo(dst)
}
