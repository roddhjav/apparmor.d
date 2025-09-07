// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package builder

import (
	"github.com/roddhjav/apparmor.d/pkg/prebuild"
	"github.com/roddhjav/apparmor.d/pkg/util"
)

var (
	regHotfix = util.ToRegexRepl([]string{
		`Cx`, `cx`,
		`PUx`, `pux`,
		`Px`, `px`,
		`Ux`, `ux`,
	})
)

type Hotfix struct {
	prebuild.Base
}

func init() {
	RegisterBuilder(&Hotfix{
		Base: prebuild.Base{
			Keyword: "hotfix",
			Msg:     "Fix: temporary solution for #74, #80 & #235",
		},
	})
}

func (b Hotfix) Apply(opt *Option, profile string) (string, error) {
	return regHotfix.Replace(profile), nil
}
