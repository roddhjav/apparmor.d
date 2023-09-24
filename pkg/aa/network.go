// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2023 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

type AddressExpr struct {
	Source      string
	Destination string
	Port        string
}


type Network struct {
	Qualifier
	Domain   string
	Type     string
	Protocol string
	AddressExpr
}

func NetworkFromLog(log map[string]string, noNewPrivs, fileInherit bool) ApparmorRule {
	return &Network{
		Qualifier: NewQualifier(false, noNewPrivs, fileInherit),
		AddressExpr: AddressExpr{
			Source:      log["laddr"],
			Destination: log["faddr"],
			Port:        log["lport"],
		},
		Domain:   log["family"],
		Type:     log["sock_type"],
		Protocol: log["protocol"],
	}
}
