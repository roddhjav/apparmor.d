// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

const (
	HAT Kind = "hat"
)

// Hat represents a single AppArmor hat.
type Hat struct {
	Base
	Name  string
	Rules Rules
}

func (p *Hat) Kind() Kind {
	return HAT
}

func (p *Hat) Constraint() constraint {
	return blockKind
}

func (p *Hat) String() string {
	return renderTemplate(p.Kind(), p)
}

func (r *Hat) Validate() error {
	return nil
}

func (r *Hat) Compare(other Rule) int {
	o, _ := other.(*Hat)
	return compare(r.Name, o.Name)
}
