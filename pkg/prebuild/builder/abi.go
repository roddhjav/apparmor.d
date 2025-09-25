// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package builder

import (
	"github.com/roddhjav/apparmor.d/pkg/prebuild"
	"github.com/roddhjav/apparmor.d/pkg/util"
)

var (
	regAbi4To3 = util.ToRegexRepl([]string{
		`abi/4.0`, `abi/3.0`,
		`  userns,`, `  # userns,`,
		`  mqueue`, `  # mqueue`,
		`  all`, `  # all`,
		`  deny mqueue`, `  # deny mqueue`,
	})
	regApparmor41To40 = util.ToRegexRepl([]string{
		`priority=[0-9\-]*`, ``,
	})
)

type ABI3 struct {
	prebuild.Base
}

type APPARMOR40 struct {
	prebuild.Base
}

func init() {
	RegisterBuilder(&ABI3{
		Base: prebuild.Base{
			Keyword: "abi3",
			Msg:     "Build: convert all profiles from abi 4.0 to abi 3.0",
		},
	})
	RegisterBuilder(&APPARMOR40{
		Base: prebuild.Base{
			Keyword: "apparmor4.0",
			Msg:     "Build: convert all profiles from apparmor 4.1 to 4.0 or less",
		},
	})
}

func (b ABI3) Apply(opt *Option, profile string) (string, error) {
	return regAbi4To3.Replace(profile), nil
}

func (b APPARMOR40) Apply(opt *Option, profile string) (string, error) {
	return regApparmor41To40.Replace(profile), nil
}
