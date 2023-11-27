// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2023 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

type PivotRoot struct {
	Qualifier
	OldRoot       string
	NewRoot       string
	TargetProfile string
}

func PivotRootFromLog(log map[string]string) ApparmorRule {
	return &PivotRoot{
		Qualifier:     NewQualifierFromLog(log),
		OldRoot:       log["srcname"],
		NewRoot:       log["name"],
		TargetProfile: "",
	}
}

func (r *PivotRoot) Less(other any) bool {
	o, _ := other.(*PivotRoot)
	if r.Qualifier.Equals(o.Qualifier) {
		if r.OldRoot == o.OldRoot {
			if r.NewRoot == o.NewRoot {
				return r.TargetProfile < o.TargetProfile
			}
			return r.NewRoot < o.NewRoot
		}
		return r.OldRoot < o.OldRoot
	}
	return r.Qualifier.Less(o.Qualifier)
}

func (r *PivotRoot) Equals(other any) bool {
	o, _ := other.(*PivotRoot)
	return r.OldRoot == o.OldRoot && r.NewRoot == o.NewRoot &&
		r.TargetProfile == o.TargetProfile &&
		r.Qualifier.Equals(o.Qualifier)
}
