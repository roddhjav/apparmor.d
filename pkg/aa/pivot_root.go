// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2023 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

type PivotRoot struct {
	Qualifier
	OldRoot       string
	NewRoot       string
	TargetProfile string
}

func PivotRootFromLog(log map[string]string, noNewPrivs, fileInherit bool) ApparmorRule {
	return &PivotRoot{
		Qualifier:     NewQualifier(false, noNewPrivs, fileInherit),
		OldRoot:       log["oldroot"],
		NewRoot:       log["root"],
		TargetProfile: log["name"],
	}
}
