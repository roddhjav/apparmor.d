// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"slices"
	"strings"
)

const (
	tokMOUNT   = "mount"
	tokREMOUNT = "remount"
	tokUMOUNT  = "umount"
)

)

type MountConditions struct {
	FsType  string
	Options []string
}

func newMountConditionsFromLog(log map[string]string) MountConditions {
	if _, present := log["flags"]; present {
		return MountConditions{
			FsType:  log["fstype"],
			Options: strings.Split(log["flags"], ", "),
		}
	}
	return MountConditions{FsType: log["fstype"]}
}

func (m MountConditions) Less(other MountConditions) bool {
	if m.FsType != other.FsType {
		return m.FsType < other.FsType
	}
	return len(m.Options) < len(other.Options)
}

func (m MountConditions) Equals(other MountConditions) bool {
	return m.FsType == other.FsType && slices.Equal(m.Options, other.Options)
}

type Mount struct {
	RuleBase
	Qualifier
	MountConditions
	Source     string
	MountPoint string
}

func newMountFromLog(log map[string]string) Rule {
	return &Mount{
		RuleBase:        newRuleFromLog(log),
		Qualifier:       newQualifierFromLog(log),
		MountConditions: newMountConditionsFromLog(log),
		Source:          log["srcname"],
		MountPoint:      log["name"],
	}
}

func (r *Mount) Less(other any) bool {
	o, _ := other.(*Mount)
	if r.Source != o.Source {
		return r.Source < o.Source
	}
	if r.MountPoint != o.MountPoint {
		return r.MountPoint < o.MountPoint
	}
	if r.MountConditions.Equals(o.MountConditions) {
		return r.MountConditions.Less(o.MountConditions)
	}
	return r.Qualifier.Less(o.Qualifier)
}

func (r *Mount) Equals(other any) bool {
	o, _ := other.(*Mount)
	return r.Source == o.Source && r.MountPoint == o.MountPoint &&
		r.MountConditions.Equals(o.MountConditions) &&
		r.Qualifier.Equals(o.Qualifier)
}

func (r *Mount) String() string {
	return renderTemplate(r.Kind(), r)
}

func (r *Mount) Constraint() constraint {
	return blockKind
}

func (r *Mount) Kind() string {
	return tokMOUNT
}

type Umount struct {
	RuleBase
	Qualifier
	MountConditions
	MountPoint string
}

func newUmountFromLog(log map[string]string) Rule {
	return &Umount{
		RuleBase:        newRuleFromLog(log),
		Qualifier:       newQualifierFromLog(log),
		MountConditions: newMountConditionsFromLog(log),
		MountPoint:      log["name"],
	}
}

func (r *Umount) Less(other any) bool {
	o, _ := other.(*Umount)
	if r.MountPoint != o.MountPoint {
		return r.MountPoint < o.MountPoint
	}
	if r.MountConditions.Equals(o.MountConditions) {
		return r.MountConditions.Less(o.MountConditions)
	}
	return r.Qualifier.Less(o.Qualifier)
}

func (r *Umount) Equals(other any) bool {
	o, _ := other.(*Umount)
	return r.MountPoint == o.MountPoint &&
		r.MountConditions.Equals(o.MountConditions) &&
		r.Qualifier.Equals(o.Qualifier)
}

func (r *Umount) String() string {
	return renderTemplate(r.Kind(), r)
}

func (r *Umount) Constraint() constraint {
	return blockKind
}

func (r *Umount) Kind() string {
	return tokUMOUNT
}

type Remount struct {
	RuleBase
	Qualifier
	MountConditions
	MountPoint string
}

func newRemountFromLog(log map[string]string) Rule {
	return &Remount{
		RuleBase:        newRuleFromLog(log),
		Qualifier:       newQualifierFromLog(log),
		MountConditions: newMountConditionsFromLog(log),
		MountPoint:      log["name"],
	}
}

func (r *Remount) Less(other any) bool {
	o, _ := other.(*Remount)
	if r.MountPoint != o.MountPoint {
		return r.MountPoint < o.MountPoint
	}
	if r.MountConditions.Equals(o.MountConditions) {
		return r.MountConditions.Less(o.MountConditions)
	}
	return r.Qualifier.Less(o.Qualifier)
}

func (r *Remount) Equals(other any) bool {
	o, _ := other.(*Remount)
	return r.MountPoint == o.MountPoint &&
		r.MountConditions.Equals(o.MountConditions) &&
		r.Qualifier.Equals(o.Qualifier)
}

func (r *Remount) String() string {
	return renderTemplate(r.Kind(), r)
}

func (r *Remount) Constraint() constraint {
	return blockKind
}

func (r *Remount) Kind() string {
	return tokREMOUNT
}
