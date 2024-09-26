// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package directive

import (
	"regexp"
	"slices"
	"strings"

	"github.com/roddhjav/apparmor.d/pkg/prebuild/cfg"
)

type FilterOnly struct {
	cfg.Base
}

type FilterExclude struct {
	cfg.Base
}

func init() {
	RegisterDirective(&FilterOnly{
		Base: cfg.Base{
			Keyword: "only",
			Msg:     "Only directive applied",
			Help:    []string{"filters..."},
		},
	})
	RegisterDirective(&FilterExclude{
		Base: cfg.Base{
			Keyword: "exclude",
			Msg:     "Exclude directive applied",
			Help:    []string{"filters..."},
		},
	})
}

func filterRuleForUs(opt *Option) bool {
	return slices.Contains(opt.ArgList, cfg.Distribution) || slices.Contains(opt.ArgList, cfg.Family)
}

func filter(only bool, opt *Option, profile string) (string, error) {
	if only && filterRuleForUs(opt) {
		return opt.Clean(profile), nil
	}
	if !only && !filterRuleForUs(opt) {
		return opt.Clean(profile), nil
	}

	inline := true
	tmp := strings.Split(opt.Raw, Keyword)
	if len(tmp) >= 1 {
		left := strings.TrimSpace(tmp[0])
		if len(left) == 0 {
			inline = false
		}
	}

	if inline {
		profile = strings.Replace(profile, opt.Raw, "", -1)
	} else {
		regRemoveParagraph := regexp.MustCompile(`(?s)` + opt.Raw + `\n.*?\n\n`)
		profile = regRemoveParagraph.ReplaceAllString(profile, "")
	}
	return profile, nil
}

func (d FilterOnly) Apply(opt *Option, profile string) (string, error) {
	return filter(true, opt, profile)
}

func (d FilterExclude) Apply(opt *Option, profile string) (string, error) {
	return filter(false, opt, profile)
}
