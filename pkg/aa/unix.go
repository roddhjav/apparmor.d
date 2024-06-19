// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"fmt"
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
		RuleBase:  newBaseFromLog(log),
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

func (r *Unix) Compare(other Rule) int {
	o, _ := other.(*Unix)
	if res := compare(r.Access, o.Access); res != 0 {
		return res
	}
	if res := compare(r.Type, o.Type); res != 0 {
		return res
	}
	if res := compare(r.Protocol, o.Protocol); res != 0 {
		return res
	}
	if res := compare(r.Address, o.Address); res != 0 {
		return res
	}
	if res := compare(r.Label, o.Label); res != 0 {
		return res
	}
	if res := compare(r.Attr, o.Attr); res != 0 {
		return res
	}
	if res := compare(r.Opt, o.Opt); res != 0 {
		return res
	}
	if res := compare(r.PeerLabel, o.PeerLabel); res != 0 {
		return res
	}
	if res := compare(r.PeerAddr, o.PeerAddr); res != 0 {
		return res
	}
	return r.Qualifier.Compare(o.Qualifier)
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
