// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package prepare

import (
	"github.com/roddhjav/apparmor.d/pkg/paths"
	"github.com/roddhjav/apparmor.d/pkg/prebuild/cfg"
	"github.com/roddhjav/apparmor.d/pkg/util"
)

type Synchronise struct {
	cfg.Base
}

func init() {
	RegisterTask(&Synchronise{
		Base: cfg.Base{
			Keyword: "synchronise",
			Msg:     "Initialize a new clean apparmor.d build directory",
		},
	})
}

func (p Synchronise) Apply() ([]string, error) {
	res := []string{}
	dirs := paths.PathList{cfg.RootApparmord, cfg.Root.Join("root"), cfg.Root.Join("systemd")}
	for _, dir := range dirs {
		if err := dir.RemoveAll(); err != nil {
			return res, err
		}
	}
	for _, name := range []string{"apparmor.d", "root"} {
		if err := util.CopyTo(paths.New(name), cfg.Root.Join(name)); err != nil {
			return res, err
		}
	}
	return res, nil
}
