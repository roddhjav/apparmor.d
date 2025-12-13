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
	Base
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

func newUnix(q Qualifier, rule rule) (Rule, error) {
	accesses, err := toAccess(UNIX, rule.GetString())
	if err != nil {
		return nil, err
	}
	if err := rule.ValidateMapKeys([]string{"type", "protocol", "addr", "label", "attr", "opt", "peer"}); err != nil {
		return nil, err
	}
	if err := rule.GetValues("peer").ValidateMapKeys([]string{"label", "addr"}); err != nil {
		return nil, err
	}
	return &Unix{
		Base:      newBase(rule),
		Qualifier: q,
		Access:    accesses,
		Type:      rule.GetValuesAsString("type"),
		Protocol:  rule.GetValuesAsString("protocol"),
		Address:   rule.GetValuesAsString("addr"),
		Label:     rule.GetValuesAsString("label"),
		Attr:      rule.GetValuesAsString("attr"),
		Opt:       rule.GetValuesAsString("opt"),
		PeerLabel: rule.GetValues("peer").GetValuesAsString("label"),
		PeerAddr:  rule.GetValues("peer").GetValuesAsString("addr"),
	}, nil
}

func newUnixFromLog(log map[string]string) Rule {
	return &Unix{
		Base:      newBaseFromLog(log),
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

func (r *Unix) Kind() Kind {
	return UNIX
}

func (r *Unix) Constraint() Constraint {
	return BlockRule
}

func (r *Unix) String() string {
	return renderTemplate(r.Kind(), r)
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

func (r *Unix) Merge(other Rule) bool {
	o, _ := other.(*Unix)

	if !r.Equal(o.Qualifier) {
		return false
	}
	if r.Type == o.Type && r.Protocol == o.Protocol && r.Address == o.Address &&
		r.Label == o.Label && r.Attr == o.Attr && r.Opt == o.Opt &&
		r.PeerLabel == o.PeerLabel && r.PeerAddr == o.PeerAddr {
		r.Access = merge(r.Kind(), "access", r.Access, o.Access)
		b := &r.Base
		return b.merge(o.Base)
	}
	return false
}

func (r *Unix) Lengths() []int {
	return []int{
		r.getLenAudit(),
		r.getLenAccess(),
		length("", r.Access),
		length("type=", r.Type),
		length("addr=", r.Address),
		length("label=", r.Label),
	}
}

func (r *Unix) setPaddings(max []int) {
	r.Paddings = append(r.Qualifier.setPaddings(max[:2]), setPaddings(
		max[2:], []string{"", "type=", "addr=", "label="},
		[]any{r.Access, r.Type, r.Address, r.Label})...,
	)
}
