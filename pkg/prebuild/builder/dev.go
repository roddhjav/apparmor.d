// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package builder

import (
	"github.com/roddhjav/apparmor.d/pkg/prebuild/cfg"
	"github.com/roddhjav/apparmor.d/pkg/util"
)

var (
	regDev = util.ToRegexRepl([]string{
		`Cx`, `cx`,
		`PUx`, `pux`,
		`Px`, `px`,
		`Ux`, `ux`,
	})
)

type Dev struct {
	cfg.Base
}

func init() {
	RegisterBuilder(&Dev{
		Base: cfg.Base{
			Keyword: "dev",
			Msg:     "Apply test development changes",
		},
	})
}

func (b Dev) Apply(profile string) string {
	return regDev.Replace(profile)
}
