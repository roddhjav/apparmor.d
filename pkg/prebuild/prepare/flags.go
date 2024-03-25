// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package prepare

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/roddhjav/apparmor.d/pkg/prebuild/cfg"
)

var (
	regFlags         = regexp.MustCompile(`flags=\(([^)]+)\)`)
	regProfileHeader = regexp.MustCompile(` {`)
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
	for _, name := range []string{"main.flags", cfg.Distribution + ".flags"} {
		path := cfg.FlagDir.Join(name)
		if !path.Exist() {
			continue
		}
		lines, _ := path.ReadFileAsLines()
		for _, line := range lines {
			if strings.HasPrefix(line, "#") || line == "" {
				continue
			}
			manifest := strings.Split(line, " ")
			profile := manifest[0]
			file := cfg.RootApparmord.Join(profile)
			if !file.Exist() {
				res = append(res, fmt.Sprintf("Profile %s not found, ignoring", profile))
				continue
			}

			// If flags is set, overwrite profile flag
			if len(manifest) > 1 {
				flags := " flags=(" + manifest[1] + ") {"
				content, err := file.ReadFile()
				if err != nil {
					return res, err
				}

				// Remove all flags definition, then set manifest' flags
				out := regFlags.ReplaceAllLiteralString(string(content), "")
				out = regProfileHeader.ReplaceAllLiteralString(out, flags)
				if err := file.WriteFile([]byte(out)); err != nil {
					return res, err
				}
			}
		}
		res = append(res, path.String())
	}
	return res, nil
}
