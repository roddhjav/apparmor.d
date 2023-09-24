// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2023 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

type Unix struct {
	Qualifier
	Access   string
	Type     string
	Protocol string
	Address  string
	Label    string
	Attr     string
	Opt      string
	Peer     string
	PeerAddr string
}

func UnixFromLog(log map[string]string, noNewPrivs, fileInherit bool) ApparmorRule {
	return &Unix{
		Qualifier: NewQualifier(false, noNewPrivs, fileInherit),
		Access:    maskToAccess[log["requested_mask"]],
		Type:      log["sock_type"],
		Protocol:  log["protocol"],
		Address:   log["addr"],
		Label:     log["peer_label"],
		Attr:      log["attr"],
		Opt:       log["opt"],
		Peer:      log["peer"],
		PeerAddr:  log["peer_addr"],
	}
}
