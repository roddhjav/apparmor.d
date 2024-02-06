// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

type Userns struct {
	Qualifier
	Create bool
}

func UsernsFromLog(log map[string]string) ApparmorRule {
	return &Userns{
		Qualifier: NewQualifierFromLog(log),
		Create:    true,
	}
}

func (r *Userns) Less(other any) bool {
	o, _ := other.(*Userns)
	if r.Qualifier.Equals(o.Qualifier) {
		return r.Create
	}
	return r.Qualifier.Less(o.Qualifier)
}

func (r *Userns) Equals(other any) bool {
	o, _ := other.(*Userns)
	return r.Create == o.Create && r.Qualifier.Equals(o.Qualifier)
}
