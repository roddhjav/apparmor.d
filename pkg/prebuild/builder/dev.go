// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package builder

import (
	"github.com/roddhjav/apparmor.d/pkg/prebuild"
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
	prebuild.Base
}

func init() {
	RegisterBuilder(&Dev{
		Base: prebuild.Base{
			Keyword: "dev",
			Msg:     "Apply test development changes",
		},
	})
}

func (b Dev) Apply(opt *Option, profile string) (string, error) {
	return regDev.Replace(profile), nil
}
