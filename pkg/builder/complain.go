// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package builder

import (
	"slices"

	"github.com/roddhjav/apparmor.d/pkg/tasks"
	"github.com/roddhjav/apparmor.d/pkg/util"
)

type Complain struct {
	tasks.BaseTask
}

// NewComplain creates a new Complain builder.
func NewComplain() *Complain {
	return &Complain{
		BaseTask: tasks.BaseTask{
			Keyword: "complain",
			Msg:     "Build: set complain flag on all profiles",
		},
	}
}

func (b Complain) Apply(opt *Option, profile string) (string, error) {
	flags := util.GetFlags(profile)
	if slices.Contains(flags, "complain") || slices.Contains(flags, "unconfined") {
		return profile, nil
	}
	flags = append(flags, "complain")
	return util.SetFlags(profile, flags), nil
}
