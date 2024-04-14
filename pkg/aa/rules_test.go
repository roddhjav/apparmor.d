// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"reflect"
	"testing"
)

func TestRule_FromLog(t *testing.T) {
	tests := []struct {
		name    string
		fromLog func(map[string]string) ApparmorRule
		log     map[string]string
		want    ApparmorRule
	}{
		{
			name: "capbability",
			fromLog: func(m map[string]string) ApparmorRule {
				return newCapabilityFromLog(m)
			},
			log:  capability1Log,
			want: capability1,
		},
		{
			name: "network",
			fromLog: func(m map[string]string) ApparmorRule {
				return newNetworkFromLog(m)
			},
			log:  network1Log,
			want: network1,
		},
		{
			name: "mount",
			fromLog: func(m map[string]string) ApparmorRule {
				return newMountFromLog(m)
			},
			log:  mount1Log,
			want: mount1,
		},
		{
			name: "umount",
			fromLog: func(m map[string]string) ApparmorRule {
				return newUmountFromLog(m)
			},
			log:  umount1Log,
			want: umount1,
		},
		{
			name: "pivotroot",
			fromLog: func(m map[string]string) ApparmorRule {
				return newPivotRootFromLog(m)
			},
			log:  pivotroot1Log,
			want: pivotroot1,
		},
		{
			name: "changeprofile",
			fromLog: func(m map[string]string) ApparmorRule {
				return newChangeProfileFromLog(m)
			},
			log:  changeprofile1Log,
			want: changeprofile1,
		},
		{
			name: "signal",
			fromLog: func(m map[string]string) ApparmorRule {
				return newSignalFromLog(m)
			},
			log:  signal1Log,
			want: signal1,
		},
		{
			name: "ptrace/xdg-document-portal",
			fromLog: func(m map[string]string) ApparmorRule {
				return newPtraceFromLog(m)
			},
			log:  ptrace1Log,
			want: ptrace1,
		},
		{
			name: "ptrace/snap-update-ns.firefox",
			fromLog: func(m map[string]string) ApparmorRule {
				return newPtraceFromLog(m)
			},
			log:  ptrace2Log,
			want: ptrace2,
		},
		{
			name: "unix",
			fromLog: func(m map[string]string) ApparmorRule {
				return newUnixFromLog(m)
			},
			log:  unix1Log,
			want: unix1,
		},
		{
			name: "dbus",
			fromLog: func(m map[string]string) ApparmorRule {
				return newDbusFromLog(m)
			},
			log:  dbus1Log,
			want: dbus1,
		},
		{
			name: "file",
			fromLog: func(m map[string]string) ApparmorRule {
				return newFileFromLog(m)
			},
			log:  file1Log,
			want: file1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fromLog(tt.log); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RuleFromLog() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRule_Less(t *testing.T) {
	tests := []struct {
		name  string
		rule  ApparmorRule
		other ApparmorRule
		want  bool
	}{
		{
			name:  "include1",
			rule:  include1,
			other: includeLocal1,
			want:  false,
		},
		{
			name:  "include2",
			rule:  include1,
			other: include2,
			want:  false,
		},
		{
			name:  "include3",
			rule:  include1,
			other: include3,
			want:  false,
		},
		{
			name:  "rlimit",
			rule:  rlimit1,
			other: rlimit2,
			want:  false,
		},
		{
			name:  "rlimit2",
			rule:  rlimit2,
			other: rlimit2,
			want:  false,
		},
		{
			name:  "rlimit3",
			rule:  rlimit1,
			other: rlimit3,
			want:  false,
		},
		{
			name:  "capability",
			rule:  capability1,
			other: capability2,
			want:  true,
		},
		{
			name:  "network",
			rule:  network1,
			other: network2,
			want:  false,
		},
		{
			name:  "mount",
			rule:  mount1,
			other: mount2,
			want:  false,
		},
		{
			name:  "umount",
			rule:  umount1,
			other: umount2,
			want:  true,
		},
		{
			name:  "pivot_root1",
			rule:  pivotroot2,
			other: pivotroot1,
			want:  true,
		},
		{
			name:  "pivot_root2",
			rule:  pivotroot1,
			other: pivotroot3,
			want:  false,
		},
		{
			name:  "change_profile1",
			rule:  changeprofile1,
			other: changeprofile2,
			want:  false,
		},
		{
			name:  "change_profile2",
			rule:  changeprofile1,
			other: changeprofile3,
			want:  true,
		},
		{
			name:  "signal",
			rule:  signal1,
			other: signal2,
			want:  true,
		},
		{
			name:  "ptrace/less",
			rule:  ptrace1,
			other: ptrace2,
			want:  true,
		},
		{
			name:  "ptrace/more",
			rule:  ptrace2,
			other: ptrace1,
			want:  false,
		},
		{
			name:  "unix",
			rule:  unix1,
			other: unix1,
			want:  false,
		},
		{
			name:  "dbus",
			rule:  dbus1,
			other: dbus1,
			want:  false,
		},
		{
			name:  "dbus2",
			rule:  dbus2,
			other: dbus3,
			want:  false,
		},
		{
			name:  "file",
			rule:  file1,
			other: file2,
			want:  true,
		},
		{
			name:  "file/empty",
			rule:  &File{},
			other: &File{},
			want:  false,
		},
		{
			name:  "file/equal",
			rule:  &File{Path: "/usr/share/poppler/cMap/Identity-H"},
			other: &File{Path: "/usr/share/poppler/cMap/Identity-H"},
			want:  false,
		},
		{
			name:  "file/owner",
			rule:  &File{Path: "/usr/share/poppler/cMap/Identity-H", Owner: true},
			other: &File{Path: "/usr/share/poppler/cMap/Identity-H"},
			want:  true,
		},
		{
			name:  "file/access",
			rule:  &File{Path: "/usr/share/poppler/cMap/Identity-H", Access: "r"},
			other: &File{Path: "/usr/share/poppler/cMap/Identity-H", Access: "w"},
			want:  true,
		},
		{
			name:  "file/close",
			rule:  &File{Path: "/usr/share/poppler/cMap/"},
			other: &File{Path: "/usr/share/poppler/cMap/Identity-H"},
			want:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.rule
			if got := r.Less(tt.other); got != tt.want {
				t.Errorf("Rule.Less() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRule_Equals(t *testing.T) {
	tests := []struct {
		name  string
		rule  ApparmorRule
		other ApparmorRule
		want  bool
	}{
		{
			name:  "include1",
			rule:  include1,
			other: includeLocal1,
			want:  false,
		},
		{
			name:  "rlimit",
			rule:  rlimit1,
			other: rlimit1,
			want:  true,
		},
		{
			name:  "capability/equal",
			rule:  capability1,
			other: capability1,
			want:  true,
		},
		{
			name:  "network/equal",
			rule:  network1,
			other: network1,
			want:  true,
		},
		{
			name:  "mount",
			rule:  mount1,
			other: mount1,
			want:  true,
		},
		{
			name:  "pivot_root",
			rule:  pivotroot1,
			other: pivotroot2,
			want:  false,
		},
		{
			name:  "change_profile",
			rule:  changeprofile1,
			other: changeprofile2,
			want:  false,
		},
		{
			name:  "signal1/equal",
			rule:  signal1,
			other: signal1,
			want:  true,
		},
		{
			name:  "ptrace/equal",
			rule:  ptrace1,
			other: ptrace1,
			want:  true,
		},
		{
			name:  "ptrace/not_equal",
			rule:  ptrace1,
			other: ptrace2,
			want:  false,
		},
		{
			name:  "unix",
			rule:  unix1,
			other: unix1,
			want:  true,
		},
		{
			name:  "dbus",
			rule:  dbus1,
			other: dbus2,
			want:  false,
		},
		{
			name:  "file",
			rule:  file2,
			other: file2,
			want:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.rule
			if got := r.Equals(tt.other); got != tt.want {
				t.Errorf("Rule.Equals() = %v, want %v", got, tt.want)
			}
		})
	}
}
