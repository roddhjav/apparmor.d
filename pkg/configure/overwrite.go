// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package configure

import (
	"os"

	"github.com/roddhjav/apparmor.d/pkg/aa"
	"github.com/roddhjav/apparmor.d/pkg/tasks"
)

// Overwrite replaces upstream profiles by our own at install time. A listed
// name is a path reserved by an upstream project, which distributions ship or
// not: our profile (if any) is always renamed <name>.<pkgname> so it never
// conflicts, and the upstream profile is disabled with a disable/<name> link
// when the target ships it.
type Overwrite struct {
	tasks.BaseTask
	entries []string
}

// NewOverwrite creates an Overwrite task from the overwrite.d config entries.
func NewOverwrite(entries []string) *Overwrite {
	return &Overwrite{
		BaseTask: tasks.BaseTask{
			Keyword: "overwrite",
			Msg:     "Overwrite dummy upstream profiles",
		},
		entries: entries,
	}
}

func (p Overwrite) Apply() ([]string, error) {
	res := []string{}
	disableDir := p.RootApparmor.Join("disable")
	for _, name := range p.entries {
		done := false

		// The <name> path is reserved for the upstream profile, whether or
		// not the distribution ships it today. Always install ours under
		// <name>.<pkgname> so it never conflicts with the upstream file.
		if origin := p.RootApparmor.Join(name); origin.Exist() {
			if err := origin.Rename(p.RootApparmor.Join(name + "." + p.Pkgname)); err != nil {
				return res, err
			}
			done = true
		}

		// Nothing to disable on the distributions that do not ship it.
		if aa.MagicRoot.Join(name).Exist() {
			if err := disableDir.MkdirAll(); err != nil {
				return res, err
			}
			if err := os.Symlink("../"+name, disableDir.Join(name).String()); err != nil {
				return res, err
			}
			done = true
		}

		if done {
			res = append(res, name)
		}
	}
	return res, nil
}
