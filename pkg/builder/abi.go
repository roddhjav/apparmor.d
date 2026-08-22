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
		`(?m)^[ \t]*if .+\{[ \t]*\n([\s\S]*?)^[ \t]*\}[ \t]*$\n?`, `${1}`,
		`(?m)^[ \t]*\} else if .+\{[ \t]*\n`, ``,
		`(?m)^[ \t]*\} else[ \t]*\{[ \t]*\n`, ``,
		`(?m)^\$\{`, `#$${`, // boolean variables
	})
)

type ABI4 struct {
	tasks.BaseTask
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

func (b ABI4) Apply(opt *Option, profile string) (string, error) {
	return regAbi5To4.Replace(profile), nil
}
