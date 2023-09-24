// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2023 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

type Capability struct {
	Qualifier
	Name string
}

func CapabilityFromLog(log map[string]string, noNewPrivs, fileInherit bool) ApparmorRule {
	return &Capability{
		Qualifier: NewQualifier(false, noNewPrivs, fileInherit),
		Name:      log["capname"],
	}
}
