// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2023 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import "golang.org/x/exp/slices"

type MountConditions struct {
	Fs      string
	Op      string
	FsType  string
	Options []string
}

type Mount struct {
	Qualifier
	MountConditions
	Source     string
	MountPoint string
}

func MountFromLog(log map[string]string, noNewPrivs, fileInherit bool) ApparmorRule {
	return &Mount{
		Qualifier: NewQualifier(false, noNewPrivs, fileInherit),
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
type Umount struct {
	Qualifier
	MountConditions
	MountPoint string
}

func UmountFromLog(log map[string]string, noNewPrivs, fileInherit bool) ApparmorRule {
	return &Umount{
		Qualifier: NewQualifier(false, noNewPrivs, fileInherit),
		MountConditions: MountConditions{
			Fs:      "",
			Op:      "",
			FsType:  log["fstype"],
			Options: []string{},
		},
		MountPoint: log["name"],
	}
}

type Remount struct {
	Qualifier
	MountConditions
	MountPoint string
}

func RemountFromLog(log map[string]string, noNewPrivs, fileInherit bool) ApparmorRule {
	return &Remount{
		Qualifier: NewQualifier(false, noNewPrivs, fileInherit),
		MountConditions: MountConditions{
			Fs:      "",
			Op:      "",
			FsType:  log["fstype"],
			Options: []string{},
		},
		MountPoint: log["name"],
	}
}
