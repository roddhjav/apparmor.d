// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package configure

import (
	"fmt"
	"os"

	"github.com/roddhjav/apparmor.d/pkg/prebuild"
	"github.com/roddhjav/apparmor.d/pkg/tasks"
)

var ext = "." + prebuild.Pkgname

type Overwrite struct {
	tasks.BaseTask
	Optional bool
}

// NewOverwrite creates a new Overwrite task with optional configuration.
func NewOverwrite(optional bool) *Overwrite {
	return &Overwrite{
		BaseTask: tasks.BaseTask{
			Keyword: "overwrite",
			Msg:     "Overwrite dummy upstream profiles",
		},
		Optional: optional,
	}
}

func (p Overwrite) Apply() ([]string, error) {
	res := []string{}
	if prebuild.ABI == 3 {
		return res, nil
	}

	disableDir := p.RootApparmor.Join("disable")
	if err := disableDir.Mkdir(); err != nil {
		return res, err
	}

	path := prebuild.DistDir.Join("overwrite")
	if !path.Exist() {
		return res, fmt.Errorf("%s not found", path)
	}
	for _, name := range path.MustReadFilteredFileAsLines() {
		origin := p.RootApparmor.Join(name)
		dest := p.RootApparmor.Join(name + ext)
		if !dest.Exist() && p.Optional {
			continue
		}
		if origin.Exist() {
			if err := origin.Rename(dest); err != nil {
				return res, err
			}
		}
		originRel, err := origin.RelFrom(dest)
		if err != nil {
			return res, err
		}
		if err := os.Symlink(originRel.String(), disableDir.Join(name).String()); err != nil {
			return res, err
		}
	}

	return res, nil
}
