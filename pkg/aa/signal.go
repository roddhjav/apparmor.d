// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2023 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

type Signal struct {
	Qualifier
	Access string
	Set    string
	Peer   string
}

func SignalFromLog(log map[string]string, noNewPrivs, fileInherit bool) ApparmorRule {
	return &Signal{
		Qualifier: NewQualifier(false, noNewPrivs, fileInherit),
		Access:    maskToAccess[log["requested_mask"]],
		Set:       log["signal"],
		Peer:      log["peer"],
	}
}

