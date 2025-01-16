// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package prepare

import (
	"fmt"
	"os"

	"github.com/roddhjav/apparmor.d/pkg/prebuild"
)

var ext = "." + prebuild.Pkgname

type Overwrite struct {
	prebuild.Base
	Optional bool
}

func init() {
	RegisterTask(&Overwrite{
		Base: prebuild.Base{
			Keyword: "overwrite",
			Msg:     "Overwrite dummy upstream profiles",
		},
		Optional: false,
	})
}

func (p Overwrite) Apply() ([]string, error) {
	res := []string{}
	if prebuild.ABI == 3 {
		return res, nil
	}

	disableDir := prebuild.RootApparmord.Join("disable")
	if err := disableDir.Mkdir(); err != nil {
		return res, err
	}

	path := prebuild.DistDir.Join("overwrite")
	if !path.Exist() {
		return res, fmt.Errorf("%s not found", path)
	}
	for _, name := range path.MustReadFilteredFileAsLines() {
		origin := prebuild.RootApparmord.Join(name)
		dest := prebuild.RootApparmord.Join(name + ext)
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
