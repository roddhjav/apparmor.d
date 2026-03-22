// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package builder

import (
	"slices"

	"github.com/roddhjav/apparmor.d/pkg/tasks"
	"github.com/roddhjav/apparmor.d/pkg/util"
)

type Enforce struct {
	tasks.BaseTask
}

// NewEnforce creates a new Enforce builder.
func NewEnforce() *Enforce {
	return &Enforce{
		BaseTask: tasks.BaseTask{
			Keyword: "enforce",
			Msg:     "Build: all profiles have been enforced",
		},
	}
}

func (b Enforce) Apply(opt *Option, profile string) (string, error) {
	flags := util.GetFlags(profile)
	idx := slices.Index(flags, "complain")
	if idx == -1 {
		return profile, nil
	}
	flags = slices.Delete(flags, idx, idx+1)
	return util.SetFlags(profile, flags), nil
}
