// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

type Userns struct {
	Rule
	Qualifier
	Create bool
}

func newUsernsFromLog(log map[string]string) *Userns {
	return &Userns{
		Rule:      newRuleFromLog(log),
		Qualifier: newQualifierFromLog(log),
		Create:    true,
	}
}

func (r *Userns) Less(other any) bool {
	o, _ := other.(*Userns)
	if r.Create != o.Create {
		return r.Create
	}
	return r.Qualifier.Less(o.Qualifier)
}

func (r *Userns) Equals(other any) bool {
	o, _ := other.(*Userns)
	return r.Create == o.Create && r.Qualifier.Equals(o.Qualifier)
}
