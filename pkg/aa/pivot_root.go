// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

const tokPIVOTROOT = "pivot_root"

type PivotRoot struct {
	RuleBase
	Qualifier
	OldRoot       string
	NewRoot       string
	TargetProfile string
}

func newPivotRootFromLog(log map[string]string) Rule {
	return &PivotRoot{
		RuleBase:      newRuleFromLog(log),
		Qualifier:     newQualifierFromLog(log),
		OldRoot:       log["srcname"],
		NewRoot:       log["name"],
		TargetProfile: "",
	}
}

func (r *PivotRoot) Less(other any) bool {
	o, _ := other.(*PivotRoot)
	if r.OldRoot != o.OldRoot {
		return r.OldRoot < o.OldRoot
	}
	if r.NewRoot != o.NewRoot {
		return r.NewRoot < o.NewRoot
	}
	if r.TargetProfile != o.TargetProfile {
		return r.TargetProfile < o.TargetProfile
	}
	return r.Qualifier.Less(o.Qualifier)
}

func (r *PivotRoot) Equals(other any) bool {
	o, _ := other.(*PivotRoot)
	return r.OldRoot == o.OldRoot && r.NewRoot == o.NewRoot &&
		r.TargetProfile == o.TargetProfile &&
		r.Qualifier.Equals(o.Qualifier)
}

func (r *PivotRoot) String() string {
	return renderTemplate(tokPIVOTROOT, r)
}
