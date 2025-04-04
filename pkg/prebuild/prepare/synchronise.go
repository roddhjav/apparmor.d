// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package prepare

import (
	"github.com/roddhjav/apparmor.d/pkg/paths"
	"github.com/roddhjav/apparmor.d/pkg/prebuild"
)

type Synchronise struct {
	prebuild.Base
	Paths []string // File or directory to sync into the build directory.
}

func init() {
	RegisterTask(&Synchronise{
		Base: prebuild.Base{
			Keyword: "synchronise",
			Msg:     "Initialize a new clean apparmor.d build directory",
		},
		Paths: []string{"apparmor.d", "share"},
	})
}

func (p Synchronise) Apply() ([]string, error) {
	res := []string{}
	if err := prebuild.Root.Join("systemd").RemoveAll(); err != nil {
		return res, err
	}
	if err := prebuild.RootApparmord.RemoveAll(); err != nil {
		return res, err
	}

	for _, name := range p.Paths {
		src := paths.New(name)
		dst := prebuild.Root.Join(name)
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
