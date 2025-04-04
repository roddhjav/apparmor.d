// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"fmt"
)

const PTRACE Kind = "ptrace"

func init() {
	requirements[PTRACE] = requirement{
		"access": []string{
			"r", "w", "rw", "read", "readby", "trace", "tracedby",
		},
	}
}

type Ptrace struct {
	Base
	Qualifier
	Access []string
	Peer   string
}

func newPtrace(q Qualifier, rule rule) (Rule, error) {
	accesses, err := toAccess(PTRACE, rule.GetString())
	if err != nil {
		return nil, err
	}
	return &Ptrace{
		Base:      newBase(rule),
		Qualifier: q,
		Access:    accesses,
		Peer:      rule.GetValuesAsString("peer"),
	}, nil
}

func newPtraceFromLog(log map[string]string) Rule {
	return &Ptrace{
		Base:      newBaseFromLog(log),
		Qualifier: newQualifierFromLog(log),
		Access:    Must(toAccess(PTRACE, log["requested_mask"])),
		Peer:      log["peer"],
	}
}

func (r *Ptrace) Kind() Kind {
	return PTRACE
}

func (r *Ptrace) Constraint() Constraint {
	return BlockRule
}

func (r *Ptrace) String() string {
	return renderTemplate(r.Kind(), r)
}

func (r *Ptrace) Validate() error {
	if err := validateValues(r.Kind(), "access", r.Access); err != nil {
		return fmt.Errorf("%s: %w", r, err)
	}
	return nil
}

func (r *Ptrace) Compare(other Rule) int {
	o, _ := other.(*Ptrace)
	if res := compare(r.Access, o.Access); res != 0 {
		return res
	}
	if res := compare(r.Peer, o.Peer); res != 0 {
		return res
	}
	return r.Qualifier.Compare(o.Qualifier)
}

func (r *Ptrace) Merge(other Rule) bool {
	o, _ := other.(*Ptrace)

	if !r.Equal(o.Qualifier) {
		return false
	}
	if r.Peer == o.Peer {
		r.Access = merge(r.Kind(), "access", r.Access, o.Access)
		b := &r.Base
		return b.merge(o.Base)
	}
	return false
}

func (r *Ptrace) Lengths() []int {
	return []int{
		r.getLenAudit(),
		r.getLenAccess(),
		length("", r.Access),
		length("peer=", r.Peer),
	}
}

func (r *Ptrace) setPaddings(max []int) {
	r.Paddings = append(r.Qualifier.setPaddings(max[:2]), setPaddings(
		max[2:], []string{"", "peer="},
		[]any{r.Access, r.Peer})...,
	)
}
