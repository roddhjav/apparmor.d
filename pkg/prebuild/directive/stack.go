// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package directive

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/roddhjav/apparmor.d/pkg/prebuild/cfg"
	"github.com/roddhjav/apparmor.d/pkg/util"
)

var (
	regRules            = regexp.MustCompile(`(?m)^profile.*{$((.|\n)*)}`)
	regEndOfRules       = regexp.MustCompile(`(?m)([\t ]*include if exists <.*>\n)+}`)
	regCleanStakedRules = util.ToRegexRepl([]string{
		`(?m)^.*include <abstractions/base>.*$`, ``, // Remove mandatory base abstraction
		`(?m)^.*@{exec_path}.*$`, ``, // Remove entry point
		`(?m)^.*(|P|p)(|U|u)(|i)x,.*$`, ``, // Remove transition rules
		`(?m)^(?:[\t ]*(?:\r?\n))+`, ``, // Remove empty lines
	})
)

type Stack struct {
	cfg.Base
}

func init() {
	RegisterDirective(&Stack{
		Base: cfg.Base{
			Keyword: "stack",
			Msg:     "Stack directive applied",
			Help:    Keyword + `stack profiles...`,
		},
	})
}

func (s Stack) Apply(opt *Option, profile string) string {
	res := ""
	for name := range opt.ArgMap {
		stackedProfile := util.MustReadFile(cfg.RootApparmord.Join(name))
		m := regRules.FindStringSubmatch(stackedProfile)
		if len(m) < 2 {
			panic(fmt.Sprintf("No profile found in %s", name))
		}
		stackedRules := m[1]
		stackedRules = regCleanStakedRules.Replace(stackedRules)
		res += "  # Stacked profile: " + name + "\n" + stackedRules + "\n"
	}

	// Insert the stacked profile at the end of the current profile, remove the stack directive
	m := regEndOfRules.FindStringSubmatch(profile)
	if len(m) <= 1 {
		panic(fmt.Sprintf("No end of rules found in %s", opt.File))
	}
	profile = strings.Replace(profile, m[0], res+m[0], -1)
	profile = strings.Replace(profile, opt.Raw, "", -1)
	return profile
}
