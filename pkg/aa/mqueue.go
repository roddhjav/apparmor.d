// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2023 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

type Mqueue struct {
	Qualifier
	Access string
	Type   string
	Label  string
}

func MqueueFromLog(log map[string]string, noNewPrivs, fileInherit bool) ApparmorRule {
	return &Mqueue{
		Qualifier: NewQualifier(false, noNewPrivs, fileInherit),
		Access:    maskToAccess[log["requested_mask"]],
		Type:      log["type"],
		Label:     log["label"],
	}
}
