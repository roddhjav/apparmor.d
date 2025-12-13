// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import "fmt"

const PIVOTROOT Kind = "pivot_root"

type PivotRoot struct {
	Base
	Qualifier
	OldRoot       string
	NewRoot       string
	TargetProfile string
}

func newPivotRoot(q Qualifier, rule rule) (Rule, error) {
	newroot, target := "", ""
	r := rule.GetSlice()
	if len(r) > 0 {
		if r[0] != tokARROW {
			newroot = r[0]
			r = r[1:]
		}
		if len(r) == 2 {
			if r[0] != tokARROW {
				return nil, fmt.Errorf("missing '%s' in rule: %s", tokARROW, rule)
			}
			target = r[1]
		}
	}
	return &PivotRoot{
		Base:          newBase(rule),
		Qualifier:     q,
		OldRoot:       rule.GetValuesAsString("oldroot"),
		NewRoot:       newroot,
		TargetProfile: target,
	}, rule.ValidateMapKeys([]string{"oldroot"})
}

func newPivotRootFromLog(log map[string]string) Rule {
	return &PivotRoot{
		Base:          newBaseFromLog(log),
		Qualifier:     newQualifierFromLog(log),
		OldRoot:       log["srcname"],
		NewRoot:       log["name"],
		TargetProfile: "",
	}
}

func (r *PivotRoot) Kind() Kind {
	return PIVOTROOT
}

func (r *PivotRoot) Constraint() Constraint {
	return BlockRule
}

func (r *PivotRoot) String() string {
	return renderTemplate(r.Kind(), r)
}

func (r *PivotRoot) Validate() error {
	return nil
}

func (r *PivotRoot) Compare(other Rule) int {
	o, _ := other.(*PivotRoot)
	if res := compare(r.OldRoot, o.OldRoot); res != 0 {
		return res
	}
	if res := compare(r.NewRoot, o.NewRoot); res != 0 {
		return res
	}
	if res := compare(r.TargetProfile, o.TargetProfile); res != 0 {
		return res
	}
	return r.Qualifier.Compare(o.Qualifier)
}

func (r *PivotRoot) Merge(other Rule) bool {
	return false // Never merge pivot root
}

func (r *PivotRoot) Lengths() []int {
	return []int{
		r.getLenAudit(),
		r.getLenAccess(),
		length("oldroot=", r.OldRoot),
		length("", r.NewRoot),
		length("", r.TargetProfile),
	}
}

func (r *PivotRoot) setPaddings(max []int) {
	r.Paddings = append(r.Qualifier.setPaddings(max[:2]), setPaddings(
		max[2:], []string{"oldroot=", "", ""},
		[]any{r.OldRoot, r.NewRoot, r.TargetProfile})...,
	)
}
