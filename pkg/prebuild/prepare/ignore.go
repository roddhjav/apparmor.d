// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package prepare

import (
	"github.com/roddhjav/apparmor.d/pkg/paths"
	"github.com/roddhjav/apparmor.d/pkg/prebuild"
)

type Ignore struct {
	prebuild.Base
}

func init() {
	RegisterTask(&Ignore{
		Base: prebuild.Base{
			Keyword: "ignore",
			Msg:     "Ignore profiles and files from:",
		},
	})
}

func (p Ignore) Apply() ([]string, error) {
	res := []string{}
	for _, name := range []string{"main", prebuild.Distribution} {
		for _, ignore := range prebuild.Ignore.Read(name) {
			profile := prebuild.Root.Join(ignore)
			if profile.NotExist() {
				files, err := prebuild.RootApparmord.ReadDirRecursiveFiltered(nil, paths.FilterNames(ignore))
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
		res = append(res, prebuild.IgnoreDir.Join(name+".ignore").String())
	}
	return res, nil
}
