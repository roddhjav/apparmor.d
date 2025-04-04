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

func (p *Hat) Constraint() Constraint {
	return BlockRule
}

func (p *Hat) String() string {
	return renderTemplate(p.Kind(), p)
}

func (p *Hat) Validate() error {
	return nil
}

func (p *Hat) Compare(other Rule) int {
	o, _ := other.(*Hat)
	return compare(p.Name, o.Name)
}

func (p *Hat) Merge(other Rule) bool {
	return false // Never merge hat blocks
}

func (p *Hat) Lengths() []int {
	return []int{} // No len for hat
}

func (p *Hat) setPaddings(max []int) {} // No paddings for hat
