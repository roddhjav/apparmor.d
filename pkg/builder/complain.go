// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package builder

import (
	"regexp"
	"slices"
	"strings"

	"github.com/roddhjav/apparmor.d/pkg/tasks"
)

var (
	regFlags         = regexp.MustCompile(`flags=\(([^)]+)\)`)
	regProfileHeader = regexp.MustCompile(` {\n`)
)

type Complain struct {
	tasks.Base
}

func init() {
	RegisterBuilder(&Complain{
		Base: tasks.Base{
			Keyword: "complain",
			Msg:     "Build: set complain flag on all profiles",
		},
	})
}

func (b Complain) Apply(opt *Option, profile string) (string, error) {
	flags := []string{}
	matches := regFlags.FindStringSubmatch(profile)
	if len(matches) != 0 {
		flags = strings.Split(matches[1], ",")
		if slices.Contains(flags, "complain") {
			return profile, nil
		}
		if slices.Contains(flags, "unconfined") {
			return profile, nil
		}
	}
	flags = append(flags, "complain")
	strFlags := " flags=(" + strings.Join(flags, ",") + ") {\n"

	// Remove all flags definition, then set manifest' flags
	profile = regFlags.ReplaceAllLiteralString(profile, "")
	return regProfileHeader.ReplaceAllLiteralString(profile, strFlags), nil
}
