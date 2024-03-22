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
	res := ""
	transition := "Px"
	transitions := []string{"P", "U", "p", "u", "PU", "pu"}
	for _, t := range transitions {
		if _, present := opt.Args[t]; present {
			transition = t + "x"
			delete(opt.Args, t)
			break
		}
	}

	for name := range opt.Args {
		content, err := rootApparmord.Join(name).ReadFile()
		if err != nil {
			panic(err)
		}
		profiletoTransition := string(content)

		p := &aa.AppArmorProfile{}
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
		res += p.String()
	}
	return strings.Replace(profile, opt.Raw, res, -1)
}
