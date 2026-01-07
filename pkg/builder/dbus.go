// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package builder

import (
	"regexp"
	"slices"
	"strings"

	"github.com/roddhjav/apparmor.d/pkg/aa"
	"github.com/roddhjav/apparmor.d/pkg/tasks"
)

var (
	resolve = map[string][]string{
		`"@{p_dbus_system}"`:  {"dbus-system", "dbus-system//&unconfined"},
		`"@{p_dbus_session}"`: {"dbus-session", "dbus-session//&unconfined"},
	}
)

// StackedDbus is a fix for https://gitlab.com/apparmor/apparmor/-/issues/537#note_2699570190
type StackedDbus struct {
	tasks.Base
}

// DbusBroker is a fix for https://gitlab.com/apparmor/apparmor/-/issues/565
type DbusBroker struct {
	tasks.Base
}

func init() {
	RegisterBuilder(&StackedDbus{
		Base: tasks.Base{
			Keyword: "stacked-dbus",
			Msg:     "Fix: resolve peer label variable in dbus rules",
		},
	})
	RegisterBuilder(&DbusBroker{
		Base: tasks.Base{
			Keyword: "dbus-broker",
			Msg:     "Fix: ignore peer name in dbus rules",
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
	if opt.Kind == aa.TunableKind {
		return profile, nil
	}

	toResolve := []string{}
	for k := range resolve {
		toResolve = append(toResolve, k)
	}

	rulesByParagraph, paragraphs, err := parse(opt.Kind, profile)
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

func (b DbusBroker) Apply(opt *Option, profile string) (string, error) {
	// Remove peer name in two cases:
	// 1. peer=(name=..., label=...) -> peer=(label=...)
	// 2. peer=(name=...),           -> (keep only the comma)

	// First, handle peer name with other attributes (has attribute after comma)
	rePeerNameWithAttrs := regexp.MustCompile(`peer=\(\s*name\s*=\s*(?:"[^"]*"|'[^']*'|[^,)\s]+)\s*,\s*(\w+\s*=)`)
	profile = rePeerNameWithAttrs.ReplaceAllString(profile, "peer=($1")

	// Second, handle peer name alone (followed by closing paren and comma)
	rePeerNameAlone := regexp.MustCompile(`peer=\(\s*name\s*=\s*(?:"[^"]*"|'[^']*'|[^,)\s]+)\s*\)\s*,`)
	profile = rePeerNameAlone.ReplaceAllString(profile, ",")
	return profile, nil
}
