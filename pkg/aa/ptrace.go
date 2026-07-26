// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"fmt"
	"strings"
)

const PTRACE Kind = "ptrace"

var ptraceAccessBits = map[string]AccessMask{
	"r":        1 << 0,
	"w":        1 << 1,
	"rw":       1 << 2,
	"read":     1 << 3,
	"write":    1 << 4,
	"readby":   1 << 5,
	"trace":    1 << 6,
	"tracedby": 1 << 7,
}

func init() {
	requirements[PTRACE] = requirement{
		"access": []string{
			"r", "w", "rw", "read", "write", "readby", "trace", "tracedby",
		},
	}
}

type Ptrace struct {
	Base
	Qualifier
	Access AccessMask
	Peer   string
}

func newPtrace(q Qualifier, rule rule) (Rule, error) {
	accesses, err := toAccess(PTRACE, rule.GetString())
	if err != nil {
		return nil, err
	}
	peers := rule.GetValuesAsSlice("peer")
	if len(peers) > 1 {
		return nil, fmt.Errorf("multiple 'peer' not allowed in rule: %s", rule)
	}
	peer := ""
	if len(peers) == 1 {
		peer = peers[0]
	}
	return &Ptrace{
		Base:      newBase(rule),
		Qualifier: q,
		Access:    accesses,
		Peer:      peer,
	}, rule.ValidateMapKeys([]string{"peer"})
}

func newPtraceFromLog(log map[string]string) Rule {
	return &Ptrace{
		Base:      newBaseFromLog(log),
		Qualifier: newQualifierFromLog(log),
		Access:    Must(toAccess(PTRACE, log["requested_mask"])),
		Peer:      log["peer"],
	}
}

func (r *Ptrace) mergeKey(b *strings.Builder) {
	writeQualifierKey(b, r.Qualifier)
	b.WriteString(r.Peer)
}

func (r *Ptrace) Kind() Kind {
	return PTRACE
}

func (r *Ptrace) AccessStrings() []string {
	return r.Access.Strings(PTRACE)
}

func (r *Ptrace) Constraint() Constraint {
	return BlockRule
}

func (r *Ptrace) String() string {
	return renderTemplate(r.Kind(), r)
}

func (r *Ptrace) Validate() error {
	if err := validateAccess(r.Kind(), r.Access); err != nil {
		return fmt.Errorf("%s: %w", r, err)
	}
	return nil
}

func (r *Ptrace) Compare(other Rule) int {
	o, _ := other.(*Ptrace)
	if res := compareAccessMask(r.Access, o.Access, PTRACE); res != 0 {
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
		r.Access |= o.Access
		b := &r.Base
		return b.merge(o.Base)
	}
	return false
}

func (r *Ptrace) Lengths() []int {
	return []int{
		r.getLenAudit(),
		r.getLenAccess(),
		length("", r.Access.Strings(r.Kind())),
		length("peer=", r.Peer),
	}
}

func (r *Ptrace) setPaddings(max []int) {
	r.Paddings = append(r.Qualifier.setPaddings(max[:2]), setPaddings(
		max[2:], []string{"", "peer="},
		[]any{r.Access.Strings(r.Kind()), r.Peer})...,
	)
}
