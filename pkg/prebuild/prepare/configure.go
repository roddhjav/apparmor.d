// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package prepare

import (
	"fmt"

	"github.com/roddhjav/apparmor.d/pkg/prebuild/cfg"
	"github.com/roddhjav/apparmor.d/pkg/util"
)

type Configure struct {
	cfg.Base
}

func init() {
	RegisterTask(&Configure{
		Base: cfg.Base{
			Keyword: "configure",
			Msg:     "Set distribution specificities",
		},
	})
}

func (p Configure) Apply() ([]string, error) {
	res := []string{}
	switch cfg.Distribution {
	case "arch", "opensuse":
		if cfg.Overwrite {
			if err := cfg.Overwrite.Apply(); err != nil {
				return res, err
			}
		}

	case "ubuntu":
		if err := cfg.DebianHide.Init(); err != nil {
			return res, err
		}

		if cfg.Overwrite {
			if err := cfg.Overwrite.Apply(); err != nil {
				return res, err
			}
		} else {
			if err := util.CopyTo(cfg.DistDir.Join("ubuntu"), cfg.RootApparmord); err != nil {
				return res, err
			}
		}

	case "debian", "whonix":
		if err := cfg.DebianHide.Init(); err != nil {
			return res, err
		}

		// Copy Debian specific abstractions
		if err := util.CopyTo(cfg.DistDir.Join("ubuntu"), cfg.RootApparmord); err != nil {
			return res, err
		}

	default:
		return []string{}, fmt.Errorf("%s is not a supported distribution", cfg.Distribution)

	}
	return res, nil
}
