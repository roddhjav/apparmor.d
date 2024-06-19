// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"fmt"
	"slices"
)

const NETWORK Kind = "network"

func init() {
	requirements[NETWORK] = requirement{
		"access": []string{
			"create", "bind", "listen", "accept", "connect", "shutdown",
			"getattr", "setattr", "getopt", "setopt", "send", "receive",
			"r", "w", "rw",
		},
		"domains": []string{
			"unix", "inet", "ax25", "ipx", "appletalk", "netrom", "bridge",
			"atmpvc", "x25", "inet6", "rose", "netbeui", "security", "key",
			"netlink", "packet", "ash", "econet", "atmsvc", "rds", "sna", "irda",
			"pppox", "wanpipe", "llc", "ib", "mpls", "can", "tipc", "bluetooth",
			"iucv", "rxrpc", "isdn", "phonet", "ieee802154", "caif", "alg",
			"nfc", "vsock", "kcm", "qipcrtr", "smc", "xdp", "mctp",
		},
		"type": []string{
			"stream", "dgram", "seqpacket", "rdm", "raw", "packet",
		},
		"protocol": []string{"tcp", "udp", "icmp"},
	}
}

type AddressExpr struct {
	Source      string
	Destination string
	Port        string
}

func newAddressExprFromLog(log map[string]string) AddressExpr {
	return AddressExpr{
		Source:      log["laddr"],
		Destination: log["faddr"],
		Port:        log["lport"],
	}
}

func (r AddressExpr) Compare(other AddressExpr) int {
	if res := compare(r.Source, other.Source); res != 0 {
		return res
	}
	if res := compare(r.Destination, other.Destination); res != 0 {
		return res
	}
	return compare(r.Port, other.Port)
}

func (r AddressExpr) Equals(other AddressExpr) bool {
	return r.Source == other.Source && r.Destination == other.Destination &&
		r.Port == other.Port
}

type Network struct {
	RuleBase
	Qualifier
	AddressExpr
	Domain   string
	Type     string
	Protocol string
}

func newNetwork(q Qualifier, rule rule) (Rule, error) {
	nType, protocol, domain := "", "", ""
	r := rule.GetSlice()
	if len(r) > 0 {
		domain = r[0]
	}
	if len(r) >= 2 {
		if slices.Contains(requirements[NETWORK]["type"], r[1]) {
			nType = r[1]
		} else if slices.Contains(requirements[NETWORK]["protocol"], r[1]) {
			protocol = r[1]
		}
	}
	return &Network{
		RuleBase:  newBase(rule),
		Qualifier: q,
		Domain:    domain,
		Type:      nType,
		Protocol:  protocol,
	}, nil
}

func newNetworkFromLog(log map[string]string) Rule {
	return &Network{
		RuleBase:    newBaseFromLog(log),
		Qualifier:   newQualifierFromLog(log),
		AddressExpr: newAddressExprFromLog(log),
		Domain:      log["family"],
		Type:        log["sock_type"],
		Protocol:    log["protocol"],
	}
}

func (r *Network) Validate() error {
	if err := validateValues(r.Kind(), "domains", []string{r.Domain}); err != nil {
		return fmt.Errorf("%s: %w", r, err)
	}
	if err := validateValues(r.Kind(), "type", []string{r.Type}); err != nil {
		return fmt.Errorf("%s: %w", r, err)
	}
	if err := validateValues(r.Kind(), "protocol", []string{r.Protocol}); err != nil {
		return fmt.Errorf("%s: %w", r, err)
	}
	return nil
}

func (r *Network) Compare(other Rule) int {
	o, _ := other.(*Network)
	if res := compare(r.Domain, o.Domain); res != 0 {
		return res
	}
	if res := compare(r.Type, o.Type); res != 0 {
		return res
	}
	if res := compare(r.Protocol, o.Protocol); res != 0 {
		return res
	}
	if res := r.AddressExpr.Compare(o.AddressExpr); res != 0 {
		return res
	}
	return r.Qualifier.Compare(o.Qualifier)
}

func (r *Network) String() string {
	return renderTemplate(r.Kind(), r)
}

func (r *Network) Constraint() constraint {
	return blockKind
}

func (r *Network) Kind() Kind {
	return NETWORK
}
