// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package prepare

import (
	"github.com/roddhjav/apparmor.d/pkg/paths"
	"github.com/roddhjav/apparmor.d/pkg/prebuild/cfg"
)

type Ignore struct {
	cfg.Base
}

func init() {
	RegisterTask(&Ignore{
		Base: cfg.Base{
			Keyword: "ignore",
			Msg:     "Ignore profiles and files from:",
		},
	})
}

func (p Ignore) Apply() ([]string, error) {
	res := []string{}
	for _, name := range []string{"main", cfg.Distribution} {
		for _, ignore := range cfg.Ignore.Read(name) {
			profile := cfg.Root.Join(ignore)
			if profile.NotExist() {
				files, err := cfg.RootApparmord.ReadDirRecursiveFiltered(nil, paths.FilterNames(ignore))
				if err != nil {
					return res, err
				}
				for _, path := range files {
					if err := path.RemoveAll(); err != nil {
						return res, err
					}
				}
			} else {
				if err := profile.RemoveAll(); err != nil {
					return res, err
				}
			}
		}
		res = append(res, cfg.IgnoreDir.Join(name+".ignore").String())
	}
	return res, nil
}
