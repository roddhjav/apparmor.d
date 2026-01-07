// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package configure

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/roddhjav/apparmor.d/pkg/prebuild"
	"github.com/roddhjav/apparmor.d/pkg/tasks"
)

var (
	regFlags         = regexp.MustCompile(`flags=\(([^)]+)\)`)
	regProfileHeader = regexp.MustCompile(` {\n`)
)

type SetFlags struct {
	tasks.Base
}

func init() {
	RegisterTask(&SetFlags{
		Base: tasks.Base{
			Keyword: "setflags",
			Msg:     "Set flags on some profiles",
		},
	})
}

func (p SetFlags) Apply() ([]string, error) {
	res := []string{}
	for _, name := range []string{"main", prebuild.Distribution} {
		for profile, flags := range prebuild.Flags.Read(name) {
			file := prebuild.RootApparmord.Join(profile)
			if !file.Exist() {
				res = append(res, fmt.Sprintf("Profile %s not found, ignoring", profile))
				continue
			}

			// Overwrite profile flags
			if len(flags) > 0 {
				flagsStr := " flags=(" + strings.Join(flags, ",") + ") {\n"
				out, err := file.ReadFileAsString()
				if err != nil {
					return res, err
				}

				// Remove all flags definition, then set manifest' flags
				out = regFlags.ReplaceAllLiteralString(out, "")
				out = regProfileHeader.ReplaceAllLiteralString(out, flagsStr)
				if err := file.WriteFile([]byte(out)); err != nil {
					return res, err
				}
			}
		}
		res = append(res, prebuild.FlagDir.Join(name+".flags").String())
	}
	return res, nil
}
