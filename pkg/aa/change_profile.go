// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import "fmt"

const CHANGEPROFILE Kind = "change_profile"

func init() {
	requirements[CHANGEPROFILE] = requirement{
		"mode": []string{"safe", "unsafe"},
	}
}

type ChangeProfile struct {
	RuleBase
	Qualifier
	ExecMode    string
	Exec        string
	ProfileName string
}

func newChangeProfileFromLog(log map[string]string) Rule {
	return &ChangeProfile{
		RuleBase:    newBaseFromLog(log),
		Qualifier:   newQualifierFromLog(log),
		ExecMode:    log["mode"],
		Exec:        log["exec"],
		ProfileName: log["target"],
	}
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

func (r *ChangeProfile) String() string {
	return renderTemplate(r.Kind(), r)
}

func (r *ChangeProfile) Constraint() constraint {
	return blockKind
}

func (r *ChangeProfile) Kind() Kind {
	return CHANGEPROFILE
}
