// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package builder

import (
	"regexp"
	"slices"
	"strings"

	"github.com/roddhjav/apparmor.d/pkg/prebuild/cfg"
)

var (
	regFlags         = regexp.MustCompile(`flags=\(([^)]+)\)`)
	regProfileHeader = regexp.MustCompile(` {`)
)

type Complain struct {
	cfg.Base
}

func init() {
	RegisterBuilder(&Complain{
		Base: cfg.Base{
			Keyword: "complain",
			Msg:     "Set complain flag on all profiles",
		},
	})
}

func (b Complain) Apply(profile string) string {
	flags := []string{}
	matches := regFlags.FindStringSubmatch(profile)
	if len(matches) != 0 {
		flags = strings.Split(matches[1], ",")
		if slices.Contains(flags, "complain") {
			return profile
		}
	}
	flags = append(flags, "complain")
	strFlags := " flags=(" + strings.Join(flags, ",") + ") {"

	// Remove all flags definition, then set manifest' flags
	profile = regFlags.ReplaceAllLiteralString(profile, "")
	return regProfileHeader.ReplaceAllLiteralString(profile, strFlags)
}
