// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package prepare

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/roddhjav/apparmor.d/pkg/prebuild/cfg"
	"github.com/roddhjav/apparmor.d/pkg/util"
)

var (
	regFlags         = regexp.MustCompile(`flags=\(([^)]+)\)`)
	regProfileHeader = regexp.MustCompile(` {\n`)
)

type SetFlags struct {
	cfg.Base
}

func init() {
	RegisterTask(&SetFlags{
		Base: cfg.Base{
			Keyword: "setflags",
			Msg:     "Set flags on some profiles",
		},
	})
}

func (p SetFlags) Apply() ([]string, error) {
	res := []string{}
	for _, name := range []string{"main", cfg.Distribution} {
		for profile, flags := range cfg.Flags.Read(name) {
			file := cfg.RootApparmord.Join(profile)
			if !file.Exist() {
				res = append(res, fmt.Sprintf("Profile %s not found, ignoring", profile))
				continue
			}

			// Overwrite profile flags
			if len(flags) > 0 {
				flagsStr := " flags=(" + strings.Join(flags, ",") + ") {\n"
				out, err := util.ReadFile(file)
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
		res = append(res, cfg.FlagDir.Join(name+".flags").String())
	}
	return res, nil
}
