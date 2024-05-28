// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"fmt"
	"slices"
)

const UNIX Kind = "unix"

func init() {
	requirements[UNIX] = requirement{
		"access": []string{
			"create", "bind", "listen", "accept", "connect", "shutdown",
			"getattr", "setattr", "getopt", "setopt", "send", "receive",
			"r", "w", "rw",
		},
	}
}

type Unix struct {
	RuleBase
	Qualifier
	Access    []string
	Type      string
	Protocol  string
	Address   string
	Label     string
	Attr      string
	Opt       string
	PeerLabel string
	PeerAddr  string
}

func newUnixFromLog(log map[string]string) Rule {
	return &Unix{
		RuleBase:  newRuleFromLog(log),
		Qualifier: newQualifierFromLog(log),
		Access:    Must(toAccess(UNIX, log["requested_mask"])),
		Type:      log["sock_type"],
		Protocol:  log["protocol"],
		Address:   log["addr"],
		Label:     log["label"],
		Attr:      log["attr"],
		Opt:       log["opt"],
		PeerLabel: log["peer"],
		PeerAddr:  log["peer_addr"],
	}
}

func (r *Unix) Validate() error {
	if err := validateValues(r.Kind(), "access", r.Access); err != nil {
		return fmt.Errorf("%s: %w", r, err)
	}
	return nil
}

func (r *Unix) Less(other any) bool {
	o, _ := other.(*Unix)
	if len(r.Access) != len(o.Access) {
		return len(r.Access) < len(o.Access)
	}
	if r.Type != o.Type {
		return r.Type < o.Type
	}
	if r.Protocol != o.Protocol {
		return r.Protocol < o.Protocol
	}
	if r.Address != o.Address {
		return r.Address < o.Address
	}
	if r.Label != o.Label {
		return r.Label < o.Label
	}
	if r.Attr != o.Attr {
		return r.Attr < o.Attr
	}
	if r.Opt != o.Opt {
		return r.Opt < o.Opt
	}
	if r.PeerLabel != o.PeerLabel {
		return r.PeerLabel < o.PeerLabel
	}
	if r.PeerAddr != o.PeerAddr {
		return r.PeerAddr < o.PeerAddr
	}
	return r.Qualifier.Less(o.Qualifier)
}

func (r *Unix) Equals(other any) bool {
	o, _ := other.(*Unix)
	return slices.Equal(r.Access, o.Access) && r.Type == o.Type &&
		r.Protocol == o.Protocol && r.Address == o.Address &&
		r.Label == o.Label && r.Attr == o.Attr && r.Opt == o.Opt &&
		r.PeerLabel == o.PeerLabel && r.PeerAddr == o.PeerAddr &&
		r.Qualifier.Equals(o.Qualifier)
}

func (r *Unix) String() string {
	return renderTemplate(r.Kind(), r)
}

func (r *Unix) Constraint() constraint {
	return blockKind
}

func (r *Unix) Kind() Kind {
	return UNIX
}
