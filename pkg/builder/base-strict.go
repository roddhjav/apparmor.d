// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package builder

import (
	"strings"

	"github.com/roddhjav/apparmor.d/pkg/tasks"
)

type BaseStrict struct {
	tasks.Base
}

func init() {
	RegisterBuilder(&BaseStrict{
		Base: tasks.Base{
			Keyword: "base-strict",
			Msg:     "Feat: use 'base-strict' as base abstraction",
		},
	})
}

func (b BaseStrict) Apply(opt *Option, profile string) (string, error) {
	profile = strings.ReplaceAll(profile,
		"include <abstractions/base>",
		"include <abstractions/base-strict>",
	)
	return profile, nil
}
