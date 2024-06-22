// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

const (
	ALL Kind = "all"
)

type All struct {
	RuleBase
}

func newAll(q Qualifier, rule rule) (Rule, error) {
	return &All{RuleBase: newBase(rule)}, nil
}

func (r *All) Validate() error {
	return nil
}

func (r *All) Compare(other Rule) int {
	return 0
}

func (r *All) Merge(other Rule) bool {
	o, _ := other.(*All)
	return r.RuleBase.merge(o.RuleBase)
}

func (r *All) String() string {
	return renderTemplate(r.Kind(), r)
}

func (r *All) Constraint() constraint {
	return blockKind
}

func (r *All) Kind() Kind {
	return ALL
}
