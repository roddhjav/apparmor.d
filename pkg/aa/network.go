// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

const tokNETWORK = "network"

func init() {
	requirements[tokNETWORK] = requirement{
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

func (r AddressExpr) Less(other AddressExpr) bool {
	if r.Source != other.Source {
		return r.Source < other.Source
	}
	if r.Destination != other.Destination {
		return r.Destination < other.Destination
	}
	return r.Port < other.Port
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

func newNetworkFromLog(log map[string]string) Rule {
	return &Network{
		RuleBase:    newRuleFromLog(log),
		Qualifier:   newQualifierFromLog(log),
		AddressExpr: newAddressExprFromLog(log),
		Domain:      log["family"],
		Type:        log["sock_type"],
		Protocol:    log["protocol"],
	}
}

func (r *Network) Less(other any) bool {
	o, _ := other.(*Network)
	if r.Domain != o.Domain {
		return r.Domain < o.Domain
	}
	if r.Type != o.Type {
		return r.Type < o.Type
	}
	if r.Protocol != o.Protocol {
		return r.Protocol < o.Protocol
	}
	if r.AddressExpr.Less(o.AddressExpr) {
		return r.AddressExpr.Less(o.AddressExpr)
	}
	return r.Qualifier.Less(o.Qualifier)
}

func (r *Network) Equals(other any) bool {
	o, _ := other.(*Network)
	return r.Domain == o.Domain && r.Type == o.Type &&
		r.Protocol == o.Protocol && r.AddressExpr.Equals(o.AddressExpr) &&
		r.Qualifier.Equals(o.Qualifier)
}

func (r *Network) String() string {
	return renderTemplate(r.Kind(), r)
}

func (r *Network) Constraint() constraint {
	return blockKind
}

func (r *Network) Kind() string {
	return tokNETWORK
}
