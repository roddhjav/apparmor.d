// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package prepare

import (
	"fmt"

	"github.com/roddhjav/apparmor.d/pkg/prebuild"
	"github.com/roddhjav/apparmor.d/pkg/util"
)

type Configure struct {
	prebuild.Base
}

func init() {
	RegisterTask(&Configure{
		Base: prebuild.Base{
			Keyword: "configure",
			Msg:     "Set distribution specificities",
		},
	})
}

func (p Configure) Apply() ([]string, error) {
	res := []string{}

	switch prebuild.Distribution {
	case "arch", "opensuse":

	case "ubuntu":
		if err := prebuild.DebianHide.Init(); err != nil {
			return res, err
		}

		if prebuild.ABI == 3 {
			if err := util.CopyTo(prebuild.DistDir.Join("ubuntu"), prebuild.RootApparmord); err != nil {
				return res, err
			}
		}

	case "debian", "whonix":
		if err := prebuild.DebianHide.Init(); err != nil {
			return res, err
		}

		// Copy Debian specific abstractions
		if err := util.CopyTo(prebuild.DistDir.Join("ubuntu"), prebuild.RootApparmord); err != nil {
			return res, err
		}

	default:
		return []string{}, fmt.Errorf("%s is not a supported distribution", prebuild.Distribution)

	}
	return res, nil
}
