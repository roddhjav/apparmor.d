// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package configure

import (
	"os"

	"github.com/roddhjav/apparmor.d/pkg/aa"
	"github.com/roddhjav/apparmor.d/pkg/tasks"
)

// Overwrite replaces upstream profiles by our own at install time. For each
// listed profile present on the install target, the upstream profile is
// disabled with a disable/<name> link and our profile (if any) is renamed
// <name>.<pkgname>. Profiles the target does not ship are left untouched.
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
		if !aa.MagicRoot.Join(name).Exist() {
			continue
		}
		if origin := p.RootApparmor.Join(name); origin.Exist() {
			if err := origin.Rename(p.RootApparmor.Join(name + "." + p.Pkgname)); err != nil {
				return res, err
			}
		}
		if err := disableDir.MkdirAll(); err != nil {
			return res, err
		}
		if err := os.Symlink("../"+name, disableDir.Join(name).String()); err != nil {
			return res, err
		}
		res = append(res, name)
	}
	return res, nil
}
