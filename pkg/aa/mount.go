// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"fmt"
)

const (
	MOUNT   Kind = "mount"
	REMOUNT Kind = "remount"
	UMOUNT  Kind = "umount"
)

func init() {
	requirements[MOUNT] = requirement{
		"flags": {
			"acl", "async", "atime", "ro", "rw", "bind", "rbind", "dev",
			"diratime", "dirsync", "exec", "iversion", "loud", "mand", "move",
			"noacl", "noatime", "nodev", "nodiratime", "noexec", "noiversion",
			"nomand", "norelatime", "nosuid", "nouser", "private", "relatime",
			"remount", "rprivate", "rshared", "rslave", "runbindable", "shared",
			"silent", "slave", "strictatime", "suid", "sync", "unbindable",
			"user", "verbose",
		},
	}
}

type MountConditions struct {
	FsType  string
	Options []string
}

func newMountConditionsFromLog(log map[string]string) MountConditions {
	if _, present := log["flags"]; present {
		return MountConditions{
			FsType:  log["fstype"],
			Options: Must(toValues(MOUNT, "flags", log["flags"])),
		}
	}
	return MountConditions{FsType: log["fstype"]}
}

func (m MountConditions) Validate() error {
	return validateValues(MOUNT, "flags", m.Options)
}

func (m MountConditions) Compare(other MountConditions) int {
	if res := compare(m.FsType, other.FsType); res != 0 {
		return res
	}
	return compare(m.Options, other.Options)
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

func (r *Mount) Validate() error {
	if err := r.MountConditions.Validate(); err != nil {
		return fmt.Errorf("%s: %w", r, err)
	}
	return nil
}

func (r *Mount) Compare(other Rule) int {
	o, _ := other.(*Mount)
	if res := compare(r.Source, o.Source); res != 0 {
		return res
	}
	if res := compare(r.MountPoint, o.MountPoint); res != 0 {
		return res
	}
	if res := r.MountConditions.Compare(o.MountConditions); res != 0 {
		return res
	}
	return r.Qualifier.Compare(o.Qualifier)
}

func (r *Mount) String() string {
	return renderTemplate(r.Kind(), r)
}

func (r *Mount) Constraint() constraint {
	return blockKind
}

func (r *Mount) Kind() Kind {
	return MOUNT
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

func (r *Umount) Validate() error {
	if err := r.MountConditions.Validate(); err != nil {
		return fmt.Errorf("%s: %w", r, err)
	}
	return nil
}

func (r *Umount) Compare(other Rule) int {
	o, _ := other.(*Umount)
	if res := compare(r.MountPoint, o.MountPoint); res != 0 {
		return res
	}
	if res := r.MountConditions.Compare(o.MountConditions); res != 0 {
		return res
	}
	return r.Qualifier.Compare(o.Qualifier)
}

func (r *Umount) String() string {
	return renderTemplate(r.Kind(), r)
}

func (r *Umount) Constraint() constraint {
	return blockKind
}

func (r *Umount) Kind() Kind {
	return UMOUNT
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

func (r *Remount) Validate() error {
	if err := r.MountConditions.Validate(); err != nil {
		return fmt.Errorf("%s: %w", r, err)
	}
	return nil
}

func (r *Remount) Compare(other Rule) int {
	o, _ := other.(*Remount)
	if res := compare(r.MountPoint, o.MountPoint); res != 0 {
		return res
	}
	if res := r.MountConditions.Compare(o.MountConditions); res != 0 {
		return res
	}
	return r.Qualifier.Compare(o.Qualifier)
}

func (r *Remount) String() string {
	return renderTemplate(r.Kind(), r)
}

func (r *Remount) Constraint() constraint {
	return blockKind
}

func (r *Remount) Kind() Kind {
	return REMOUNT
}
