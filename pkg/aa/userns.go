// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import "fmt"

const USERNS Kind = "userns"

type Userns struct {
	RuleBase
	Qualifier
	Create bool
}

func newUsernsFromLog(log map[string]string) Rule {
	return &Userns{
		RuleBase:  newRuleFromLog(log),
		Qualifier: newQualifierFromLog(log),
		Create:    true,
	}
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

func (r *Userns) String() string {
	return renderTemplate(r.Kind(), r)
}

func (r *Userns) Constraint() constraint {
	return blockKind
}

func (r *Userns) Kind() Kind {
	return USERNS
}
