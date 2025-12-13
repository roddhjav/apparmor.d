// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

const (
	HAT  Kind = "hat"
	IF   Kind = "if"
	ELSE Kind = "else"
)

// Hat represents a single AppArmor hat.
type Hat struct {
	Base
	Name  string
	Flags []string
	Rules Rules
}

func newHat(rule rule) (*Hat, error) {
	name := ""
	if len(rule) > 0 {
		name = strings.TrimPrefix(rule.Get(0), HAT.Tok())
	}
	return &Hat{
		Name:  name,
		Flags: rule.GetValuesAsSlice(tokFLAGS),
	}, nil
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
