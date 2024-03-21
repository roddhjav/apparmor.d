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
			usage:   `#aa:exec [P|U|p|u|i|] profiles_name...`,
		},
	}
}

func (d Exec) Apply(opt *Option, profile string) string {
	res := ""
	transition := "Px"
	for name := range opt.Args {
		tmp, err := rootApparmord.Join(name).ReadFile()
		if err != nil {
			panic(err)
		}
		profiletoTransition := string(tmp)

		p := aa.DefaultTunables()
		p.ParseVariables(profiletoTransition)
		for _, variable := range p.Variables {
			if variable.Name == "exec_path" {
				for _, value := range variable.Values {
					res += "  " + value + " " + transition + ",\n"
				}
			}
		}
		profile = strings.Replace(profile, opt.Raw, res, -1)
	}
	return profile
}
