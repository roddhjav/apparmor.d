// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package directive

import (
	"strings"

	"github.com/roddhjav/apparmor.d/pkg/aa"
)

type Exec struct {
	DirectiveBase
}

func init() {
	Directives["exec"] = &Exec{
		DirectiveBase: DirectiveBase{
			message: "Exec directive applied",
			usage:   `#aa:exec [P|U|p|u|PU|pu|] profiles_name...`,
		},
	}
}

func (d Exec) Apply(opt *Option, profile string) string {
	transition := "Px"
	transitions := []string{"P", "U", "p", "u", "PU", "pu"}
	for _, t := range transitions {
		if _, present := opt.Args[t]; present {
			transition = t + "x"
			delete(opt.Args, t)
			break
		}
	}

	p := &aa.AppArmorProfile{}
	for name := range opt.Args {
		content, err := rootApparmord.Join(name).ReadFile()
		if err != nil {
			panic(err)
		}
		profiletoTransition := string(content)

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
