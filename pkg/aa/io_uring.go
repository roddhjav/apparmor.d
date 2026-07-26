// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"fmt"
	"strings"
)

const IOURING Kind = "io_uring"

var ioUringAccessBits = map[string]AccessMask{
	"sqpoll":         1 << 0,
	"override_creds": 1 << 1,
}

func init() {
	requirements[IOURING] = requirement{
		"access": []string{"sqpoll", "override_creds"},
	}
}

type IOUring struct {
	Base
	Qualifier
	Access AccessMask
	Label  string
}

func newIOUring(q Qualifier, rule rule) (Rule, error) {
	accesses, err := toAccess(IOURING, rule.GetString())
	if err != nil {
		return nil, err
	}
	if err := rule.ValidateNonEmptyValues([]string{"label"}); err != nil {
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

func (r *IOUring) mergeKey(b *strings.Builder) {
	writeQualifierKey(b, r.Qualifier)
	b.WriteString(r.Label)
}

func (r *IOUring) Kind() Kind {
	return IOURING
}

func (r *IOUring) AccessStrings() []string {
	return r.Access.Strings(IOURING)
}

func (r *IOUring) Constraint() Constraint {
	return BlockRule
}

func (r *IOUring) String() string {
	return renderTemplate(r.Kind(), r)
}

func (r *IOUring) Validate() error {
	if err := validateAccess(r.Kind(), r.Access); err != nil {
		return fmt.Errorf("%s: %w", r, err)
	}
	return nil
}

func (r *IOUring) Compare(other Rule) int {
	o, _ := other.(*IOUring)
	if res := compareAccessMask(r.Access, o.Access, IOURING); res != 0 {
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
		r.Access |= o.Access
		b := &r.Base
		return b.merge(o.Base)
	}
	return false
}

func (r *IOUring) Lengths() []int {
	return []int{
		r.getLenAudit(),
		r.getLenAccess(),
		length("", r.Access.Strings(r.Kind())),
		length("label=", r.Label),
	}
}

func (r *IOUring) setPaddings(max []int) {
	r.Paddings = append(r.Qualifier.setPaddings(max[:2]), setPaddings(
		max[2:], []string{"", "label="},
		[]any{r.Access.Strings(r.Kind()), r.Label})...,
	)
}
