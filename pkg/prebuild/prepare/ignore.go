// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package prepare

import (
	"fmt"
	"strings"

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
			// Ignore file from share/
			path := prebuild.Root.Join(ignore)
			if path.Exist() {
				if err := path.RemoveAll(); err != nil {
					return res, err
				}
				continue
			}

			// Ignore file from apparmor.d/
			profile := strings.TrimPrefix(ignore, prebuild.Src+"/")
			if strings.HasPrefix(ignore, prebuild.Src) {
				path = prebuild.RootApparmord.Join(profile)
			}
			if path.Exist() {
				if err := path.RemoveAll(); err != nil {
					return res, err
				}

			} else {
				files, err := prebuild.RootApparmord.ReadDirRecursiveFiltered(nil, paths.FilterNames(profile))
				if err != nil {
					return res, err
				}
				if len(files) == 0 {
					return res, fmt.Errorf("%s.ignore: no files found for '%s'", name, profile)
				}
				for _, path := range files {
					if err := path.RemoveAll(); err != nil {
						return res, err
					}
				}

			}
		}
		res = append(res, prebuild.IgnoreDir.Join(name+".ignore").String())
	}
	return res, nil
}
