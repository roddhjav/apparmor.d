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
	Path string
}

func init() {
	RegisterTask(&Synchronise{
		Base: prebuild.Base{
			Keyword: "synchronise",
			Msg:     "Initialize a new clean apparmor.d build directory",
		},
		Path: "",
	})
}

func (p Synchronise) Apply() ([]string, error) {
	res := []string{}
	dirs := paths.PathList{prebuild.RootApparmord, prebuild.Root.Join("share"), prebuild.Root.Join("systemd")}
	for _, dir := range dirs {
		if err := dir.RemoveAll(); err != nil {
			return res, err
		}
	}
	if p.Path == "" {
		for _, name := range []string{"apparmor.d", "share"} {
			if err := paths.CopyTo(paths.New(name), prebuild.Root.Join(name)); err != nil {
				return res, err
			}
		}
	} else {
		file := paths.New(p.Path)
		destination, err := file.RelFrom(paths.New("apparmor.d"))
		if err != nil {
			return res, err
		}
		destination = prebuild.RootApparmord.JoinPath(destination)
		if err := destination.Parent().MkdirAll(); err != nil {
			return res, err
		}
		if err := file.CopyTo(destination); err != nil {
			return res, err
		}
		res = append(res, destination.String())
	}
	return res, nil
}
