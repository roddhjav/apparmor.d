// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

type ChangeProfile struct {
	Qualifier
	ExecMode    string
	Exec        string
	ProfileName string
}

func ChangeProfileFromLog(log map[string]string) ApparmorRule {
	return &ChangeProfile{
		Qualifier:   NewQualifierFromLog(log),
		ExecMode:    log["mode"],
		Exec:        log["exec"],
		ProfileName: log["target"],
	}
}

func (r *ChangeProfile) Less(other any) bool {
	o, _ := other.(*ChangeProfile)
	if r.ExecMode == o.ExecMode {
		if r.Exec == o.Exec {
			return r.ProfileName < o.ProfileName
		}
		return r.Exec < o.Exec
	}
	return r.ExecMode < o.ExecMode
}

func (r *ChangeProfile) Equals(other any) bool {
	o, _ := other.(*ChangeProfile)
	return r.ExecMode == o.ExecMode && r.Exec == o.Exec && r.ProfileName == o.ProfileName
}
