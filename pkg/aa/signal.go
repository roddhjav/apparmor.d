// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"fmt"
)

const SIGNAL Kind = "signal"

func init() {
	requirements[SIGNAL] = requirement{
		"access": {
			"r", "w", "rw", "read", "write", "send", "receive",
		},
		"set": {
			"abrt", "alrm", "bus", "chld", "cont", "emt", "exists", "fpe", "hup",
			"ill", "int", "io", "kill", "pipe", "prof", "pwr", "quit", "segv",
			"stkflt", "stop", "stp", "sys", "term", "trap", "ttin", "ttou",
			"urg", "usr1", "usr2", "vtalrm", "winch", "xcpu", "xfsz",
			"rtmin+0", "rtmin+1", "rtmin+2", "rtmin+3", "rtmin+4",
			"rtmin+5", "rtmin+6", "rtmin+7", "rtmin+8", "rtmin+9", "rtmin+10",
			"rtmin+11", "rtmin+12", "rtmin+13", "rtmin+14", "rtmin+15",
			"rtmin+16", "rtmin+17", "rtmin+18", "rtmin+19", "rtmin+20",
			"rtmin+21", "rtmin+22", "rtmin+23", "rtmin+24", "rtmin+25",
			"rtmin+26", "rtmin+27", "rtmin+28", "rtmin+29", "rtmin+30",
			"rtmin+31", "rtmin+32",
		},
	}
}

type Signal struct {
	Base
	Qualifier
	Access []string
	Set    []string
	Peer   string
}

func newSignal(q Qualifier, rule rule) (Rule, error) {
	accesses, err := toAccess(SIGNAL, rule.GetString())
	if err != nil {
		return nil, err
	}
	set, err := toValues(SIGNAL, "set", rule.GetValuesAsString("set"))
	if err != nil {
		return nil, err
	}
	return &Signal{
		Base:      newBase(rule),
		Qualifier: q,
		Access:    accesses,
		Set:       set,
		Peer:      rule.GetValuesAsString("peer"),
	}, nil
}

func newSignalFromLog(log map[string]string) Rule {
	return &Signal{
		Base:      newBaseFromLog(log),
		Qualifier: newQualifierFromLog(log),
		Access:    Must(toAccess(SIGNAL, log["requested_mask"])),
		Set:       []string{log["signal"]},
		Peer:      log["peer"],
	}
}

func (r *Signal) Kind() Kind {
	return SIGNAL
}

func (r *Signal) Constraint() Constraint {
	return BlockRule
}

func (r *Signal) String() string {
	return renderTemplate(r.Kind(), r)
}

func (r *Signal) Validate() error {
	if err := validateValues(r.Kind(), "access", r.Access); err != nil {
		return fmt.Errorf("%s: %w", r, err)
	}
	if err := validateValues(r.Kind(), "set", r.Set); err != nil {
		return fmt.Errorf("%s: %w", r, err)
	}
	return nil
}

func (r *Signal) Compare(other Rule) int {
	o, _ := other.(*Signal)
	if res := compare(r.Access, o.Access); res != 0 {
		return res
	}
	if res := compare(r.Set, o.Set); res != 0 {
		return res
	}
	if res := compare(r.Peer, o.Peer); res != 0 {
		return res
	}
	return r.Qualifier.Compare(o.Qualifier)
}

func (r *Signal) Merge(other Rule) bool {
	o, _ := other.(*Signal)

	if !r.Equal(o.Qualifier) {
		return false
	}
	switch {
	case r.Peer == o.Peer && compare(r.Set, o.Set) == 0:
		r.Access = merge(r.Kind(), "access", r.Access, o.Access)
		b := &r.Base
		return b.merge(o.Base)
	case r.Peer == o.Peer && compare(r.Access, o.Access) == 0:
		r.Set = merge(r.Kind(), "set", r.Set, o.Set)
		b := &r.Base
		return b.merge(o.Base)
	}
	return false
}

func (r *Signal) Lengths() []int {
	return []int{
		r.getLenAudit(),
		r.getLenAccess(),
		length("", r.Access),
		length("set=", r.Set),
		length("peer=", r.Peer),
	}
}

func (r *Signal) setPaddings(max []int) {
	r.Paddings = append(r.Qualifier.setPaddings(max[:2]), setPaddings(
		max[2:], []string{"", "set=", "peer="},
		[]any{r.Access, r.Set, r.Peer})...,
	)
}
