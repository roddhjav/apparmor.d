// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

const (
	ALL Kind = "all"
)

type All struct {
	Base
}

func newAll(q Qualifier, rule rule) (Rule, error) {
	return &All{Base: newBase(rule)}, rule.ValidateMapKeys([]string{})
}

func (r *All) Kind() Kind {
	return ALL
}

func (r *All) Constraint() Constraint {
	return BlockRule
}

func (r *All) String() string {
	return renderTemplate(r.Kind(), r)
}

func (r *All) Validate() error {
	return nil
}

func (r *All) Compare(other Rule) int {
	return 0
}

func (r *All) Merge(other Rule) bool {
	o, _ := other.(*All)
	b := &r.Base
	return b.merge(o.Base) // Always merge all rules
}

func (r *All) Lengths() []int {
	return []int{} // No len for all
}

func (r *All) setPaddings(max []int) {} // No paddings for all
