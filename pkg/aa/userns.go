// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import "fmt"

const USERNS Kind = "userns"

type Userns struct {
	Base
	Qualifier
	Create bool
}

func newUserns(q Qualifier, rule rule) (Rule, error) {
	var create bool
	switch len(rule) {
	case 0:
		create = true
	case 1:
		if rule.Get(0) != "create" {
			return nil, fmt.Errorf("invalid userns format: %s", rule)
		}
		create = true
	default:
		return nil, fmt.Errorf("invalid userns format: %s", rule)
	}
	return &Userns{
		Base:      newBase(rule),
		Qualifier: q,
		Create:    create,
	}, rule.ValidateMapKeys([]string{})
}

func newUsernsFromLog(log map[string]string) Rule {
	return &Userns{
		Base:      newBaseFromLog(log),
		Qualifier: newQualifierFromLog(log),
		Create:    true,
	}
}

func (r *Userns) Kind() Kind {
	return USERNS
}

func (r *Userns) Constraint() Constraint {
	return BlockRule
}

func (r *Userns) String() string {
	return renderTemplate(r.Kind(), r)
}

func (r *Userns) Validate() error {
	return nil
}

func (r *Userns) Compare(other Rule) int {
	o, _ := other.(*Userns)
	if res := compare(r.Create, o.Create); res != 0 {
		return res
	}
	return r.Qualifier.Compare(o.Qualifier)
}

func (r *Userns) Merge(other Rule) bool {
	o, _ := other.(*Userns)
	b := &r.Base
	return b.merge(o.Base) // Always merge userns rules
}

func (r *Userns) Lengths() []int {
	return []int{} // No len for userns
}

func (r *Userns) setPaddings(max []int) {} // No paddings for userns
