// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package directive

import (
	"strings"

	"github.com/roddhjav/apparmor.d/pkg/aa"
	"github.com/roddhjav/apparmor.d/pkg/prebuild/cfg"
	"github.com/roddhjav/apparmor.d/pkg/util"
	"golang.org/x/exp/slices"
)

type Exec struct {
	cfg.Base
}

func init() {
	RegisterDirective(&Exec{
		Base: cfg.Base{
			Keyword: "exec",
			Msg:     "Exec directive applied",
			Help:    Keyword + `exec [P|U|p|u|PU|pu|] profiles...`,
		},
	})
}

func (d Exec) Apply(opt *Option, profile string) string {
	transition := "Px"
	transitions := []string{"P", "U", "p", "u", "PU", "pu"}
	t := opt.ArgList[0]
	if slices.Contains(transitions, t) {
		transition = t + "x"
		delete(opt.ArgMap, t)
	}

	p := &aa.AppArmorProfile{}
	for name := range opt.ArgMap {
		profiletoTransition := util.MustReadFile(cfg.RootApparmord.Join(name))
		dstProfile := aa.DefaultTunables()
		dstProfile.ParseVariables(profiletoTransition)
		for _, variable := range dstProfile.Variables {
			if variable.Name == "exec_path" {
				for _, v := range variable.Values {
					p.Rules = append(p.Rules, &aa.File{
						Path:   v,
						Access: transition,
					})
				}
				break
			}
		}
	}
	p.Sort()
	rules := p.String()
	lenRules := len(rules)
	rules = rules[:lenRules-1]
	return strings.Replace(profile, opt.Raw, rules, -1)
}
