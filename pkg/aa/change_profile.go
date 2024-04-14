// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

type ChangeProfile struct {
	Rule
	Qualifier
	ExecMode    string
	Exec        string
	ProfileName string
}

func newChangeProfileFromLog(log map[string]string) *ChangeProfile {
	return &ChangeProfile{
		Rule:        newRuleFromLog(log),
		Qualifier:   newQualifierFromLog(log),
		ExecMode:    log["mode"],
		Exec:        log["exec"],
		ProfileName: log["target"],
	}
}

func (r *ChangeProfile) Less(other any) bool {
	o, _ := other.(*ChangeProfile)
	if r.ExecMode != o.ExecMode {
		return r.ExecMode < o.ExecMode
	}
	if r.Exec != o.Exec {
		return r.Exec < o.Exec
	}
	if r.ProfileName != o.ProfileName {
		return r.ProfileName < o.ProfileName
	}
	return r.Qualifier.Less(o.Qualifier)
}

func (r *ChangeProfile) Equals(other any) bool {
	o, _ := other.(*ChangeProfile)
	return r.ExecMode == o.ExecMode && r.Exec == o.Exec &&
		r.ProfileName == o.ProfileName && r.Qualifier.Equals(o.Qualifier)
}
