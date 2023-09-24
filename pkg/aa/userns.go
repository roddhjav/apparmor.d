// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2023 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

type Userns struct {
	Qualifier
	Create bool
}

func UsernsFromLog(log map[string]string, noNewPrivs, fileInherit bool) ApparmorRule {
	return &Userns{
		Qualifier: NewQualifier(false, noNewPrivs, fileInherit),
		Create:    true,
	}
}

