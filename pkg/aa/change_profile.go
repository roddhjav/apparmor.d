// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2023 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

type ChangeProfile struct {
	ExecMode    string
	Exec        string
	ProfileName string
}

func ChangeProfileFromLog(log map[string]string, noNewPrivs, fileInherit bool) ApparmorRule {
	return &ChangeProfile{
		ExecMode:    log["mode"],
		Exec:        log["exec"],
		ProfileName: log["name"],
	}
}
