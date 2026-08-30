// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package builder

import (
	"strings"

	"github.com/roddhjav/apparmor.d/pkg/tasks"
	"github.com/roddhjav/apparmor.d/pkg/util"
)

var (
	regAbi5To4 = util.ToRegexRepl([]string{
		`abi/5.0`, `abi/4.0`,
		`(?m)^\$\{`, `#$${`, // boolean variables
	})
)

// stripConditions removes the if/else scaffolding lines, keeping the rules.
func stripConditions(profile string) string {
	lines := strings.Split(profile, "\n")
	res := make([]string, 0, len(lines))
	stack := []bool{} // for each open block, whether it is a condition
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		closes := strings.HasPrefix(trimmed, "}")
		opens := strings.HasSuffix(trimmed, "{")
		drop := false

		if closes && len(stack) > 0 {
			drop = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
		}
		if opens {
			isCond := strings.HasPrefix(trimmed, "if ") ||
				(closes && strings.HasPrefix(trimmed, "} else"))
			stack = append(stack, isCond)
			drop = isCond
		}
		if !drop {
			res = append(res, line)
		}
	}
	return strings.Join(res, "\n")
}

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
	return regAbi5To4.Replace(stripConditions(profile)), nil
}
