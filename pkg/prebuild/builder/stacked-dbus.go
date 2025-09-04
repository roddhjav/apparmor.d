// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package builder

import (
	"slices"
	"strings"

	"github.com/roddhjav/apparmor.d/pkg/aa"
	"github.com/roddhjav/apparmor.d/pkg/prebuild"
)

var (
	resolve = map[string][]string{
		`"@{p_dbus_system}"`:  {"dbus-system", "dbus-system//&unconfined"},
		`"@{p_dbus_session}"`: {"dbus-session", "dbus-session//&unconfined"},
	}
)

// StackedDbus is a fix for https://gitlab.com/apparmor/apparmor/-/issues/537#note_2699570190
type StackedDbus struct {
	prebuild.Base
}

func init() {
	RegisterBuilder(&StackedDbus{
		Base: prebuild.Base{
			Keyword: "stacked-dbus",
			Msg:     "Fix: resolve peer label variable in dbus rules",
		},
	})
}

func parse(kind aa.FileKind, profile string) (aa.ParaRules, []string, error) {
	var raw string
	paragraphs := []string{}
	rulesByParagraph := aa.ParaRules{}

	switch kind {
	case aa.ProfileKind:
		f := &aa.AppArmorProfileFile{}
		nb, err := f.Parse(profile)
		if err != nil {
			return nil, nil, err
		}
		lines := strings.Split(profile, "\n")
		raw = strings.Join(lines[nb:], "\n")

	case aa.AbstractionKind, aa.TunableKind:
		raw = profile
	}

	r, par, err := aa.ParseRules(raw)
	if err != nil {
		return nil, nil, err
	}
	rulesByParagraph = append(rulesByParagraph, r...)
	paragraphs = append(paragraphs, par...)
	return rulesByParagraph, paragraphs, nil
}

func (b StackedDbus) Apply(opt *Option, profile string) (string, error) {
	kind := aa.KindFromPath(opt.File)
	if kind == aa.TunableKind {
		return profile, nil
	}

	toResolve := []string{}
	for k := range resolve {
		toResolve = append(toResolve, k)
	}

	rulesByParagraph, paragraphs, err := parse(kind, profile) //
	if err != nil {
		return "", err
	}
	for idx, rules := range rulesByParagraph {
		changed := false
		newRules := aa.Rules{}
		for _, rule := range rules {
			switch rule := rule.(type) {
			case *aa.Dbus:
				if slices.Contains(toResolve, rule.PeerLabel) {
					changed = true
					for _, label := range resolve[rule.PeerLabel] {
						newRule := *rule
						newRule.PeerLabel = label
						newRules = append(newRules, &newRule)
					}
				} else {
					newRules = append(newRules, rule)
				}
			default:
				newRules = append(newRules, rule)
			}
		}
		if changed {
			profile = strings.ReplaceAll(profile, paragraphs[idx], newRules.String()+"\n")
		}
	}
	return profile, nil
}
