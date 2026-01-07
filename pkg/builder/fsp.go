// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package builder

import (
	"github.com/roddhjav/apparmor.d/pkg/tasks"
	"github.com/roddhjav/apparmor.d/pkg/util"
)

var (
	regFullSystemPolicy = util.ToRegexRepl([]string{
		`(PU|U)x,`, `Px,`,
	})
)

type FullSystemPolicy struct {
	tasks.Base
}

func init() {
	RegisterBuilder(&FullSystemPolicy{
		Base: tasks.Base{
			Keyword: "fsp",
			Msg:     "Feat: prevent unconfined transitions in profile rules",
		},
	})
}

func (b FullSystemPolicy) Apply(opt *Option, profile string) (string, error) {
	return regFullSystemPolicy.Replace(profile), nil
}
