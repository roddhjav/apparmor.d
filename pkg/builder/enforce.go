// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package builder

import (
	"slices"
	"strings"

	"github.com/roddhjav/apparmor.d/pkg/tasks"
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
	matches := regFlags.FindStringSubmatch(profile)
	if len(matches) == 0 {
		return profile, nil
	}

	flags := strings.Split(matches[1], ",")
	idx := slices.Index(flags, "complain")
	if idx == -1 {
		return profile, nil
	}
	flags = slices.Delete(flags, idx, idx+1)
	strFlags := "{\n"
	if len(flags) >= 1 {
		strFlags = " flags=(" + strings.Join(flags, ",") + ") {\n"
	}

	// Remove all flags definition, then set new flags
	profile = regFlags.ReplaceAllLiteralString(profile, "")
	return regProfileHeader.ReplaceAllLiteralString(profile, strFlags), nil
}
