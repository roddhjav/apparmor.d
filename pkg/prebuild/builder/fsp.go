// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package builder

import (
	"github.com/roddhjav/apparmor.d/pkg/prebuild"
	"github.com/roddhjav/apparmor.d/pkg/util"
)

var (
	regFullSystemPolicy = util.ToRegexRepl([]string{
		`(PU|U)x,`, `Px,`,
	})
)

type FullSystemPolicy struct {
	prebuild.Base
}

func init() {
	RegisterBuilder(&FullSystemPolicy{
		Base: prebuild.Base{
			Keyword: "fsp",
			Msg:     "Prevent unconfined transitions in profile rules",
		},
	})
}

func (b FullSystemPolicy) Apply(opt *Option, profile string) (string, error) {
	return regFullSystemPolicy.Replace(profile), nil
}
