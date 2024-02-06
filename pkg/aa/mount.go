// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import "golang.org/x/exp/slices"

type MountConditions struct {
	Fs      string
	Op      string
	FsType  string
	Options []string
}

func (m MountConditions) Less(other MountConditions) bool {
	if m.Fs == other.Fs {
		if m.Op == other.Op {
			if m.FsType == other.FsType {
				return len(m.Options) < len(other.Options)
			}
			return m.FsType < other.FsType
		}
		return m.Op < other.Op
	}
	return m.Fs < other.Fs
}

func (m MountConditions) Equals(other MountConditions) bool {
	return m.Fs == other.Fs && m.Op == other.Op && m.FsType == other.FsType &&
		slices.Equal(m.Options, other.Options)
}

type Mount struct {
	Qualifier
	MountConditions
	Source     string
	MountPoint string
}

func MountFromLog(log map[string]string) ApparmorRule {
	return &Mount{
		Qualifier: NewQualifierFromLog(log),
		MountConditions: MountConditions{
			Fs:      "",
			Op:      "",
			FsType:  log["fstype"],
			Options: []string{},
		},
		Source:     log["srcname"],
		MountPoint: log["name"],
	}
}

func (r *Mount) Less(other any) bool {
	o, _ := other.(*Mount)
	if r.Qualifier.Equals(o.Qualifier) {
		if r.Source == o.Source {
			if r.MountPoint == o.MountPoint {
				return r.MountConditions.Less(o.MountConditions)
			}
			return r.MountPoint < o.MountPoint
		}
		return r.Source < o.Source
	}
	return r.Qualifier.Less(o.Qualifier)
}

func (r *Mount) Equals(other any) bool {
	o, _ := other.(*Mount)
	return r.Source == o.Source && r.MountPoint == o.MountPoint &&
		r.MountConditions.Equals(o.MountConditions) &&
		r.Qualifier.Equals(o.Qualifier)
}

type Umount struct {
	Qualifier
	MountConditions
	MountPoint string
}

func UmountFromLog(log map[string]string) ApparmorRule {
	return &Umount{
		Qualifier: NewQualifierFromLog(log),
		MountConditions: MountConditions{
			Fs:      "",
			Op:      "",
			FsType:  log["fstype"],
			Options: []string{},
		},
		MountPoint: log["name"],
	}
}

func (r *Umount) Less(other any) bool {
	o, _ := other.(*Umount)
	if r.Qualifier.Equals(o.Qualifier) {
		if r.MountPoint == o.MountPoint {
			return r.MountConditions.Less(o.MountConditions)
		}
		return r.MountPoint < o.MountPoint
	}
	return r.Qualifier.Less(o.Qualifier)
}

func (r *Umount) Equals(other any) bool {
	o, _ := other.(*Umount)
	return r.MountPoint == o.MountPoint &&
		r.MountConditions.Equals(o.MountConditions) &&
		r.Qualifier.Equals(o.Qualifier)
}

type Remount struct {
	Qualifier
	MountConditions
	MountPoint string
}

func RemountFromLog(log map[string]string) ApparmorRule {
	return &Remount{
		Qualifier: NewQualifierFromLog(log),
		MountConditions: MountConditions{
			Fs:      "",
			Op:      "",
			FsType:  log["fstype"],
			Options: []string{},
		},
		MountPoint: log["name"],
	}
}

func (r *Remount) Less(other any) bool {
	o, _ := other.(*Remount)
	if r.Qualifier.Equals(o.Qualifier) {
		if r.MountPoint == o.MountPoint {
			return r.MountConditions.Less(o.MountConditions)
		}
		return r.MountPoint < o.MountPoint
	}
	return r.Qualifier.Less(o.Qualifier)
}

func (r *Remount) Equals(other any) bool {
	o, _ := other.(*Remount)
	return r.MountPoint == o.MountPoint &&
		r.MountConditions.Equals(o.MountConditions) &&
		r.Qualifier.Equals(o.Qualifier)
}
