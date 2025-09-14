// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package directive

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/roddhjav/apparmor.d/pkg/prebuild"
)

type FilterOnly struct {
	prebuild.Base
}

type FilterExclude struct {
	prebuild.Base
}

func init() {
	RegisterDirective(&FilterOnly{
		Base: prebuild.Base{
			Keyword: "only",
			Msg:     "Only directive applied",
			Help:    []string{"filters..."},
		},
	})
	RegisterDirective(&FilterExclude{
		Base: prebuild.Base{
			Keyword: "exclude",
			Msg:     "Exclude directive applied",
			Help:    []string{"filters..."},
		},
	})
}

func filterRuleForUs(opt *Option) bool {
	if prebuild.RBAC && slices.Contains(opt.ArgList, "RBAC") {
		return true
	}

	if prebuild.Test && slices.Contains(opt.ArgList, "test") {
		return true
	}

	abiStr := fmt.Sprintf("abi%d", prebuild.ABI)
	if slices.Contains(opt.ArgList, abiStr) {
		return true
	}
	versionStr := fmt.Sprintf("apparmor%.1f", prebuild.Version)
	if slices.Contains(opt.ArgList, versionStr) {
		return true
	}
	return slices.Contains(opt.ArgList, prebuild.Distribution) || slices.Contains(opt.ArgList, prebuild.Family)
}

func filter(only bool, opt *Option, profile string) (string, error) {
	if only && filterRuleForUs(opt) {
		return opt.Clean(profile), nil
	}
	if !only && !filterRuleForUs(opt) {
		return opt.Clean(profile), nil
	}

	if opt.IsInline() {
		profile = strings.ReplaceAll(profile, opt.Raw, "")
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
