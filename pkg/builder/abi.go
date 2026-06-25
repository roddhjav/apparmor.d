// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package builder

import (
	"github.com/roddhjav/apparmor.d/pkg/tasks"
	"github.com/roddhjav/apparmor.d/pkg/util"
)

var (
	regAbi5To4 = util.ToRegexRepl([]string{
		`abi/5.0`, `abi/4.0`,
		`(?m)^([ \t]*)if (.+\S)\s*\{[ \t]*$([\s\S]*?)^([ \t]*)\}[ \t]*$`, `${1}# if ${2}${3}${4}# }`,
		`(?m)^([ \t]*)\} else if (.+\S)\s*\{[ \t]*$`, `${1}# } else if ${2}`,
		`(?m)^([ \t]*)\} else\s*\{[ \t]*$`, `${1}# } else`,
	})
	regAbi4To5 = util.ToRegexRepl([]string{
		`abi/4.0`, `abi/5.0`,
	})
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

type ABI5 struct {
	tasks.BaseTask
}

type ABI4 struct {
	tasks.BaseTask
}

type ABI3 struct {
	tasks.BaseTask
}

type APPARMOR40 struct {
	tasks.BaseTask
}

// NewABI3 creates a new ABI3 builder.
func NewABI3() *ABI3 {
	return &ABI3{
		BaseTask: tasks.BaseTask{
			Keyword: "abi3",
			Msg:     "Build: convert all profiles from abi 4.0 to abi 3.0",
		},
	}
}

// NewABI4 creates a new ABI4 builder.
func NewABI4() *ABI4 {
	return &ABI4{
		BaseTask: tasks.BaseTask{
			Keyword: "abi4",
			Msg:     "Build: convert all profiles from abi 5.0 to abi 4.0",
		},
	}
}

// NewABI5 creates a new ABI5 builder.
func NewABI5() *ABI5 {
	return &ABI5{
		BaseTask: tasks.BaseTask{
			Keyword: "abi5",
			Msg:     "Build: convert all profiles from abi 4.0 to abi 5.0",
		},
	}
}

// NewAPPARMOR40 creates a new APPARMOR40 builder.
func NewAPPARMOR40() *APPARMOR40 {
	return &APPARMOR40{
		BaseTask: tasks.BaseTask{
			Keyword: "apparmor4.0",
			Msg:     "Build: convert all profiles from apparmor 4.1 to 4.0 or less",
		},
	}
}

func (b ABI5) Apply(opt *Option, profile string) (string, error) {
	return regAbi4To5.Replace(profile), nil
}

func (b ABI4) Apply(opt *Option, profile string) (string, error) {
	return regAbi5To4.Replace(profile), nil
}

func (b ABI3) Apply(opt *Option, profile string) (string, error) {
	return regAbi4To3.Replace(profile), nil
}

func (b APPARMOR40) Apply(opt *Option, profile string) (string, error) {
	return regApparmor41To40.Replace(profile), nil
}
