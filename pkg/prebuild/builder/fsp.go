// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package builder

import (
	"github.com/roddhjav/apparmor.d/pkg/prebuild/cfg"
	"github.com/roddhjav/apparmor.d/pkg/util"
)

var (
	regFullSystemPolicy = util.ToRegexRepl([]string{
		`r(PU|U)x,`, `rPx,`,
	})
)

type FullSystemPolicy struct {
	cfg.Base
}

func init() {
	RegisterBuilder(&FullSystemPolicy{
		Base: cfg.Base{
			Keyword: "fsp",
			Msg:     "Prevent unconfined transitions in profile rules",
		},
	})
}

func (b FullSystemPolicy) Apply(opt *Option, profile string) (string, error) {
	return regFullSystemPolicy.Replace(profile), nil
}
