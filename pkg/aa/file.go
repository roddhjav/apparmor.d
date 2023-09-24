// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2023 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

type File struct {
	Qualifier
	Path   string
	Access string
	Target string
}

func FileFromLog(log map[string]string, noNewPrivs, fileInherit bool) ApparmorRule {
	owner := false
	if log["fsuid"] == log["ouid"] && log["OUID"] != "root" {
		owner = true
	}
	return &File{
		Qualifier: NewQualifier(owner, noNewPrivs, fileInherit),
		Path:      log["name"],
		Access:    maskToAccess[log["requested_mask"]],
		Target:    log["target"],
	}
}
