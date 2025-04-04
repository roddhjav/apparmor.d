// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

// TODO: Local variables in profile header need to be resolved

package directive

import (
	"fmt"
	"slices"
	"strings"

	"github.com/roddhjav/apparmor.d/pkg/aa"
	"github.com/roddhjav/apparmor.d/pkg/prebuild"
)

type Exec struct {
	prebuild.Base
}

func init() {
	RegisterDirective(&Exec{
		Base: prebuild.Base{
			Keyword: "exec",
			Msg:     "Exec directive applied",
			Help:    []string{"[P|U|p|u|PU|pu|] profiles..."},
		},
	})
}

func (d Exec) Apply(opt *Option, profileRaw string) (string, error) {
	if len(opt.ArgList) == 0 {
		return "", fmt.Errorf("no profile to exec")
	}
	transition := "Px"
	transitions := []string{"P", "U", "p", "u", "PU", "pu"}
	t := opt.ArgList[0]
	if slices.Contains(transitions, t) {
		transition = t + "x"
		delete(opt.ArgMap, t)
	}

	rules := aa.Rules{}
	for name := range opt.ArgMap {
		profiletoTransition := prebuild.RootApparmord.Join(name).MustReadFileAsString()
		dstProfile := aa.DefaultTunables()
		if _, err := dstProfile.Parse(profiletoTransition); err != nil {
			return "", err
		}
		if err := dstProfile.Resolve(); err != nil {
			return "", err
		}
		for _, variable := range dstProfile.Preamble.GetVariables() {
			if variable.Name == "exec_path" {
				for _, v := range variable.Values {
					rules = append(rules, &aa.File{
						Path:   v,
						Access: []string{transition},
					})
				}
				break
			}
		}
	}

	aa.IndentationLevel = strings.Count(
		strings.SplitN(opt.Raw, Keyword, 1)[0], aa.Indentation,
	)
	rules = rules.Sort()
	new := rules.String()
	new = new[:len(new)-1]
	return strings.ReplaceAll(profileRaw, opt.Raw, new), nil
}
