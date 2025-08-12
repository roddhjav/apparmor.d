// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"fmt"
	"slices"
)

const CHANGEPROFILE Kind = "change_profile"

func init() {
	requirements[CHANGEPROFILE] = requirement{
		"mode": []string{"safe", "unsafe"},
	}
}

type ChangeProfile struct {
	Base
	Qualifier
	ExecMode    string
	Exec        string
	ProfileName string
}

func newChangeProfile(q Qualifier, rule rule) (Rule, error) {
	mode, exec, target := "", "", ""
	if len(rule) > 0 {
		if slices.Contains(requirements[CHANGEPROFILE]["mode"], rule.Get(0)) {
			mode = rule.Get(0)
			rule = rule[1:]
		}
		if len(rule) > 0 {
			if rule.Get(0) != tokARROW {
				exec = rule.Get(0)
				if len(rule) > 2 {
					if rule.Get(1) != tokARROW {
						return nil, fmt.Errorf("missing '%s' in rule: %s", tokARROW, rule)
					}
					target = rule.Get(2)
				}
			} else {
				if len(rule) > 1 {
					target = rule.Get(1)
				}
			}
		}
	}
	return &ChangeProfile{
		Base:        newBase(rule),
		Qualifier:   q,
		ExecMode:    mode,
		Exec:        exec,
		ProfileName: target,
	}, nil
}

func newChangeProfileFromLog(log map[string]string) Rule {
	return &ChangeProfile{
		Base:        newBaseFromLog(log),
		Qualifier:   newQualifierFromLog(log),
		ExecMode:    log["mode"],
		Exec:        log["exec"],
		ProfileName: log["target"],
	}
}

func (r *ChangeProfile) Kind() Kind {
	return CHANGEPROFILE
}

func (r *ChangeProfile) Constraint() Constraint {
	return BlockRule
}

func (r *ChangeProfile) String() string {
	return renderTemplate(r.Kind(), r)
}

func (r *ChangeProfile) Validate() error {
	if err := validateValues(r.Kind(), "mode", []string{r.ExecMode}); err != nil {
		return fmt.Errorf("%s: %w", r, err)
	}
	return nil
}

func (r *ChangeProfile) Compare(other Rule) int {
	o, _ := other.(*ChangeProfile)
	if res := compare(r.ExecMode, o.ExecMode); res != 0 {
		return res
	}
	if res := compare(r.Exec, o.Exec); res != 0 {
		return res
	}
	if res := compare(r.ProfileName, o.ProfileName); res != 0 {
		return res
	}
	return r.Qualifier.Compare(o.Qualifier)
}

func (r *ChangeProfile) Merge(other Rule) bool {
	return false // Never merge change_profile
}

func (r *ChangeProfile) Lengths() []int {
	return []int{
		r.getLenAudit(),
		r.getLenAccess(),
		length("", r.ExecMode),
		length("", r.Exec),
		length("", r.ProfileName),
	}
}

func (r *ChangeProfile) setPaddings(max []int) {
	r.Paddings = append(r.Qualifier.setPaddings(max[:2]), setPaddings(
		max[2:], []string{"", "", ""},
		[]any{r.ExecMode, r.Exec, r.ProfileName})...,
	)
}
