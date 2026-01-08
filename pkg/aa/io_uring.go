// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"fmt"
)

const IOURING Kind = "io_uring"

func init() {
	requirements[IOURING] = requirement{
		"access": []string{"sqpoll", "override_creds"},
	}
}

type IOUring struct {
	Base
	Qualifier
	Access []string
	Label  string
}

func newIOUring(q Qualifier, rule rule) (Rule, error) {
	accesses, err := toAccess(IOURING, rule.GetString())
	if err != nil {
		return nil, err
	}
	return &IOUring{
		Base:      newBase(rule),
		Qualifier: q,
		Access:    accesses,
		Label:     rule.GetValuesAsString("label"),
	}, rule.ValidateMapKeys([]string{"label"})
}

func newIOUringFromLog(log map[string]string) Rule {
	return &IOUring{
		Base:      newBaseFromLog(log),
		Qualifier: newQualifierFromLog(log),
		Access:    Must(toAccess(IOURING, log["requested"])),
		Label:     log["label"],
	}
}

func (r *IOUring) Kind() Kind {
	return IOURING
}

func (r *IOUring) Constraint() Constraint {
	return BlockRule
}

func (r *IOUring) String() string {
	return renderTemplate(r.Kind(), r)
}

func (r *IOUring) Validate() error {
	if err := validateValues(r.Kind(), "access", r.Access); err != nil {
		return fmt.Errorf("%s: %w", r, err)
	}
	return nil
}

func (r *IOUring) Compare(other Rule) int {
	o, _ := other.(*IOUring)
	if res := compare(r.Access, o.Access); res != 0 {
		return res
	}
	if res := compare(r.Label, o.Label); res != 0 {
		return res
	}
	return r.Qualifier.Compare(o.Qualifier)
}

func (r *IOUring) Merge(other Rule) bool {
	o, _ := other.(*IOUring)

	if !r.Equal(o.Qualifier) {
		return false
	}
	if r.Label == o.Label {
		r.Access = merge(r.Kind(), "access", r.Access, o.Access)
		b := &r.Base
		return b.merge(o.Base)
	}
	return false
}

func (r *IOUring) Lengths() []int {
	return []int{
		r.getLenAudit(),
		r.getLenAccess(),
		length("", r.Access),
		length("label=", r.Label),
	}
}

func (r *IOUring) setPaddings(max []int) {
	r.Paddings = append(r.Qualifier.setPaddings(max[:2]), setPaddings(
		max[2:], []string{"", "label="},
		[]any{r.Access, r.Label})...,
	)
}
