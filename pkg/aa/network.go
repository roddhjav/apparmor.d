// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2023 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

type AddressExpr struct {
	Source      string
	Destination string
	Port        string
}

func (r AddressExpr) Equals(other AddressExpr) bool {
	return r.Source == other.Source && r.Destination == other.Destination &&
		r.Port == other.Port
}

func (r AddressExpr) Less(other AddressExpr) bool {
	if r.Source == other.Source {
		if r.Destination == other.Destination {
			return r.Port < other.Port
		}
		return r.Destination < other.Destination
	}
	return r.Source < other.Source
}

type Network struct {
	Qualifier
	Domain   string
	Type     string
	Protocol string
	AddressExpr
}

func NetworkFromLog(log map[string]string) ApparmorRule {
	return &Network{
		Qualifier: NewQualifierFromLog(log),
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

func (r *Network) Less(other any) bool {
	o, _ := other.(*Network)
	if r.Qualifier.Equals(o.Qualifier) {
		if r.Domain == o.Domain {
			if r.Type == o.Type {
				return r.Protocol < o.Protocol
			}
			return r.Type < o.Type
		}
		return r.Domain < o.Domain
	}
	return r.Qualifier.Less(o.Qualifier)
}

func (r *Network) Equals(other any) bool {
	o, _ := other.(*Network)
	return r.Domain == o.Domain && r.Type == o.Type &&
		r.Protocol == o.Protocol && r.AddressExpr.Equals(o.AddressExpr) &&
		r.Qualifier.Equals(o.Qualifier)
}
