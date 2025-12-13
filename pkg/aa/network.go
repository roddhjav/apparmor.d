// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"fmt"
	"net"
	"slices"
	"strconv"
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

type LocalAddress struct {
	IP   string
	Port string
}

func newLocalAddress(rule rule) (LocalAddress, error) {
	return LocalAddress{
		IP:   rule.GetValuesAsString("ip"),
		Port: rule.GetValuesAsString("port"),
	}, nil
}

func newLocalAddressFromLog(log map[string]string) LocalAddress {
	return LocalAddress{
		IP:   log["laddr"],
		Port: log["lport"],
	}
}

func (r LocalAddress) Validate() error {
	if r.IP != "" && r.IP != "none" && net.ParseIP(r.IP) == nil {
		return fmt.Errorf("invalid IP address: %s", r.IP)
	}
	if r.Port != "" {
		port, err := strconv.Atoi(r.Port)
		if err != nil || port < 0 || port > 65535 {
			return fmt.Errorf("invalid port: %s", r.Port)
		}
	}
	return nil
}

func (r LocalAddress) Compare(other LocalAddress) int {
	if res := compare(r.IP, other.IP); res != 0 {
		return res
	}
	return compare(r.Port, other.Port)
}

type PeerAddress struct {
	IP   string
	Port string
	Src  string
}

func newPeerAddress(rule rule) (PeerAddress, error) {
	return PeerAddress{
		IP:   rule.GetValues("peer").GetValuesAsString("ip"),
		Port: rule.GetValues("peer").GetValuesAsString("port"),
	}, nil
}

func newPeerAddressFromLog(log map[string]string) PeerAddress {
	return PeerAddress{
		IP:   log["faddr"],
		Port: log["fport"],
		Src:  log["saddr"],
	}
}

func (r PeerAddress) Validate() error {
	if r.IP != "" && r.IP != "none" && net.ParseIP(r.IP) == nil {
		return fmt.Errorf("invalid IP address: %s", r.IP)
	}
	if r.Port != "" {
		port, err := strconv.Atoi(r.Port)
		if err != nil || port < 0 || port > 65535 {
			return fmt.Errorf("invalid port: %s", r.Port)
		}
	}
	return nil
}

func (r PeerAddress) Compare(other PeerAddress) int {
	if res := compare(r.IP, other.IP); res != 0 {
		return res
	}
	if res := compare(r.Port, other.Port); res != 0 {
		return res
	}
	return compare(r.Src, other.Src)
}

type Network struct {
	Base
	Qualifier
	LocalAddress
	PeerAddress
	Access   []string
	Domain   string
	Type     string
	Protocol string
}

func newNetwork(q Qualifier, rule rule) (Rule, error) {
	var accesses []string
	nType, protocol, domain := "", "", ""

	// Classify each token as access, domain, type, or protocol
	for _, token := range rule.GetSlice() {
		switch {
		case slices.Contains(requirements[NETWORK]["access"], token):
			accesses = append(accesses, token)
		case slices.Contains(requirements[NETWORK]["domains"], token):
			domain = token
		case slices.Contains(requirements[NETWORK]["type"], token):
			nType = token
		case slices.Contains(requirements[NETWORK]["protocol"], token):
			protocol = token
		}
	}

	localAdress, err := newLocalAddress(rule)
	if err != nil {
		return nil, err
	}
	peerAddress, err := newPeerAddress(rule)
	if err != nil {
		return nil, err
	}
	return &Network{
		Base:         newBase(rule),
		Qualifier:    q,
		LocalAddress: localAdress,
		PeerAddress:  peerAddress,
		Access:       accesses,
		Domain:       domain,
		Type:         nType,
		Protocol:     protocol,
	}, nil
}

func newNetworkFromLog(log map[string]string) Rule {
	return &Network{
		Base:         newBaseFromLog(log),
		Qualifier:    newQualifierFromLog(log),
		LocalAddress: newLocalAddressFromLog(log),
		PeerAddress:  newPeerAddressFromLog(log),
		Access:       Must(toAccess(NETWORK, log["requested"])),
		Domain:       log["family"],
		Type:         log["sock_type"],
		Protocol:     log["protocol"],
	}
}

func (r *Network) Kind() Kind {
	return NETWORK
}

func (r *Network) Constraint() Constraint {
	return BlockRule
}

func (r *Network) String() string {
	return renderTemplate(r.Kind(), r)
}

func (r *Network) Validate() error {
	if err := validateValues(r.Kind(), "access", r.Access); err != nil {
		return fmt.Errorf("%s: %w", r, err)
	}
	if err := validateValues(r.Kind(), "domains", []string{r.Domain}); err != nil {
		return fmt.Errorf("%s: %w", r, err)
	}
	if err := validateValues(r.Kind(), "type", []string{r.Type}); err != nil {
		return fmt.Errorf("%s: %w", r, err)
	}
	if err := validateValues(r.Kind(), "protocol", []string{r.Protocol}); err != nil {
		return fmt.Errorf("%s: %w", r, err)
	}
	if err := r.LocalAddress.Validate(); err != nil {
		return fmt.Errorf("%s: %w", r, err)
	}
	if err := r.PeerAddress.Validate(); err != nil {
		return fmt.Errorf("%s: %w", r, err)
	}
	return nil
}

func (r *Network) Compare(other Rule) int {
	o, _ := other.(*Network)
	if res := compare(r.Domain, o.Domain); res != 0 {
		return res
	}
	if res := compare(r.Access, o.Access); res != 0 {
		return res
	}
	if res := compare(r.Type, o.Type); res != 0 {
		return res
	}
	if res := compare(r.Protocol, o.Protocol); res != 0 {
		return res
	}
	if res := r.LocalAddress.Compare(o.LocalAddress); res != 0 {
		return res
	}
	if res := r.PeerAddress.Compare(o.PeerAddress); res != 0 {
		return res
	}
	return r.Qualifier.Compare(o.Qualifier)
}

func (r *Network) Merge(other Rule) bool {
	return false // Never merge network
}

func (r *Network) Lengths() []int {
	return []int{
		r.getLenAudit(),
		r.getLenAccess(),
		length("", r.Domain),
		length("", r.Type),
		length("", r.Protocol),
	}
}

func (r *Network) setPaddings(max []int) {
	r.Paddings = append(r.Qualifier.setPaddings(max[:2]), setPaddings(
		max[2:], []string{"", "", ""},
		[]any{r.Domain, r.Type, r.Protocol})...,
	)
}
