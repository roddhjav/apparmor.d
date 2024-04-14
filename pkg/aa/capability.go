// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

type Capability struct {
	Rule
	Qualifier
	Name string
}

func newCapabilityFromLog(log map[string]string) *Capability {
	return &Capability{
		Rule:      newRuleFromLog(log),
		Qualifier: newQualifierFromLog(log),
		Name:      log["capname"],
	}
}

func (r *Capability) Less(other any) bool {
	o, _ := other.(*Capability)
	if r.Name != o.Name {
		return r.Name < o.Name
	}
	return r.Qualifier.Less(o.Qualifier)
}

func (r *Capability) Equals(other any) bool {
	o, _ := other.(*Capability)
	return r.Name == o.Name && r.Qualifier.Equals(o.Qualifier)
}
