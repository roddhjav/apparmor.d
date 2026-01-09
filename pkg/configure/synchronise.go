// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package configure

import (
	"github.com/roddhjav/apparmor.d/pkg/paths"
	"github.com/roddhjav/apparmor.d/pkg/tasks"
)

type Synchronise struct {
	tasks.BaseTask
	Sources []*paths.Path // Files or directories to sync into the build directory.
}

// NewSynchronise creates a new Synchronise task.
func NewSynchronise(sources []*paths.Path) *Synchronise {
	return &Synchronise{
		BaseTask: tasks.BaseTask{
			Keyword: "synchronise",
			Msg:     "Initialize a new clean apparmor.d directory",
		},
		Sources: sources,
	}
}

func (p Synchronise) Apply() ([]string, error) {
	res := []string{}
	for _, src := range p.Sources {
		dst := p.Root.Join(src.Base())
		if err := dst.RemoveAll(); err != nil {
			return res, err
		}

		if src.IsDir() {
			if err := src.CopyFS(dst); err != nil {
				return res, err
			}
		} else {
			if err := dst.Parent().MkdirAll(); err != nil {
				return res, err
			}
			if err := src.CopyTo(dst); err != nil {
				return res, err
			}
		}
		res = append(res, dst.String())
	}
	return res, nil
}
