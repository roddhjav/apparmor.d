// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package builder

import (
	"github.com/roddhjav/apparmor.d/pkg/prebuild/cfg"
	"github.com/roddhjav/apparmor.d/pkg/util"
)

var (
	regAbi4To3 = util.ToRegexRepl([]string{ // Currently Abi3 -> Abi4
		`abi/3.0`, `abi/4.0`,
		`# userns,`, `userns,`,
		`# mqueue`, `mqueue`,
	})
)

type ABI3 struct {
	cfg.Base
}

func init() {
	RegisterBuilder(&ABI3{
		Base: cfg.Base{
			Keyword: "abi3",
			Msg:     "Convert all profiles from abi 4.0 to abi 3.0",
		},
	})
}

func (b ABI3) Apply(profile string) (string, error) {
	return regAbi4To3.Replace(profile), nil
}
