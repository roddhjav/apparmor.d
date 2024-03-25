// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package prepare

import (
	"strings"

	"github.com/arduino/go-paths-helper"
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
	for _, name := range []string{"main.ignore", cfg.Distribution + ".ignore"} {
		path := cfg.DistDir.Join("ignore", name)
		if !path.Exist() {
			continue
		}
		lines, _ := path.ReadFileAsLines()
		for _, line := range lines {
			if strings.HasPrefix(line, "#") || line == "" {
				continue
			}
			profile := cfg.Root.Join(line)
			if profile.NotExist() {
				files, err := cfg.RootApparmord.ReadDirRecursiveFiltered(nil, paths.FilterNames(line))
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
		res = append(res, path.String())
	}
	return res, nil
}
