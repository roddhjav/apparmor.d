// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2023 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

type Ptrace struct {
	Qualifier
	Access string
	Peer   string
}

func PtraceFromLog(log map[string]string, noNewPrivs, fileInherit bool) ApparmorRule {
	return &Ptrace{
		Qualifier: NewQualifier(false, noNewPrivs, fileInherit),
		Access:    maskToAccess[log["requested_mask"]],
		Peer:      log["peer"],
	}
}

