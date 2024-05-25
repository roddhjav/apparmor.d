// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"reflect"
	"testing"
)

func TestRules_FromLog(t *testing.T) {
	tests := []struct {
		name    string
		fromLog func(map[string]string) Rule
		log     map[string]string
		want    Rule
	}{
		{
			name:    "capbability",
			fromLog: newCapabilityFromLog,
			log:     capability1Log,
			want:    capability1,
		},
		{
			name:    "network",
			fromLog: newNetworkFromLog,
			log:     network1Log,
			want:    network1,
		},
		{
			name:    "mount",
			fromLog: newMountFromLog,
			log:     mount1Log,
			want:    mount1,
		},
		{
			name:    "umount",
			fromLog: newUmountFromLog,
			log:     umount1Log,
			want:    umount1,
		},
		{
			name:    "pivotroot",
			fromLog: newPivotRootFromLog,
			log:     pivotroot1Log,
			want:    pivotroot1,
		},
		{
			name:    "changeprofile",
			fromLog: newChangeProfileFromLog,
			log:     changeprofile1Log,
			want:    changeprofile1,
		},
		{
			name:    "signal",
			fromLog: newSignalFromLog,
			log:     signal1Log,
			want:    signal1,
		},
		{
			name:    "ptrace/xdg-document-portal",
			fromLog: newPtraceFromLog,
			log:     ptrace1Log,
			want:    ptrace1,
		},
		{
			name:    "ptrace/snap-update-ns.firefox",
			fromLog: newPtraceFromLog,
			log:     ptrace2Log,
			want:    ptrace2,
		},
		{
			name:    "unix",
			fromLog: newUnixFromLog,
			log:     unix1Log,
			want:    unix1,
		},
		{
			name:    "dbus",
			fromLog: newDbusFromLog,
			log:     dbus1Log,
			want:    dbus1,
		},
		{
			name:    "file",
			fromLog: newFileFromLog,
			log:     file1Log,
			want:    file1,
		},
		{
			name:    "link",
			fromLog: newLinkFromLog,
			log:     link1Log,
			want:    link1,
		},
		{
			name:    "link",
			fromLog: newFileFromLog,
			log:     link3Log,
			want:    link3,
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

func TestRules_Less(t *testing.T) {
	tests := []struct {
		name  string
		rule  Rule
		other Rule
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
			want:  false,
		},
		{
			name:  "ptrace/less",
			rule:  ptrace1,
			other: ptrace2,
			want:  false,
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
			rule:  &File{Path: "/usr/share/poppler/cMap/Identity-H", Access: []string{"r"}},
			other: &File{Path: "/usr/share/poppler/cMap/Identity-H", Access: []string{"w"}},
			want:  false,
		},
		{
			name:  "file/close",
			rule:  &File{Path: "/usr/share/poppler/cMap/"},
			other: &File{Path: "/usr/share/poppler/cMap/Identity-H"},
			want:  true,
		},
		{
			name:  "link",
			rule:  link1,
			other: link2,
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

func TestRules_Equals(t *testing.T) {
	tests := []struct {
		name  string
		rule  Rule
		other Rule
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
		{
			name:  "link",
			rule:  link1,
			other: link3,
			want:  false,
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

func TestRules_String(t *testing.T) {
	tests := []struct {
		name string
		rule Rule
		want string
	}{
		{
			name: "include1",
			rule: include1,
			want: "include <abstraction/base>",
		},
		{
			name: "include-local",
			rule: includeLocal1,
			want: "include if exists <local/foo>",
		},
		{
			name: "include-abs",
			rule: &Include{Path: "/usr/share/apparmor.d/", IsMagic: false},
			want: `include "/usr/share/apparmor.d/"`,
		},
		{
			name: "rlimit",
			rule: rlimit1,
			want: "set rlimit nproc <= 200,",
		},
		{
			name: "capability",
			rule: capability1,
			want: "capability net_admin,",
		},
		{
			name: "capability/multi",
			rule: &Capability{Names: []string{"dac_override", "dac_read_search"}},
			want: "capability dac_override dac_read_search,",
		},
		{
			name: "capability/all",
			rule: &Capability{},
			want: "capability,",
		},
		{
			name: "network",
			rule: network1,
			want: "network netlink raw,",
		},
		{
			name: "mount",
			rule: mount1,
			want: "mount fstype=overlay overlay -> /var/lib/docker/overlay2/opaque-bug-check1209538631/merged/,  # failed perms check",
		},
		{
			name: "pivot_root",
			rule: pivotroot1,
			want: "pivot_root oldroot=@{run}/systemd/mount-rootfs/ @{run}/systemd/mount-rootfs/,",
		},
		{
			name: "change_profile",
			rule: changeprofile1,
			want: "change_profile -> systemd-user,",
		},
		{
			name: "signal",
			rule: signal1,
			want: "signal receive set=kill peer=firefox//&firejail-default,",
		},
		{
			name: "ptrace",
			rule: ptrace1,
			want: "ptrace read peer=nautilus,",
		},
		{
			name: "unix",
			rule: unix1,
			want: "unix (send receive) type=stream protocol=0 addr=none peer=(label=dbus-daemon, addr=@/tmp/dbus-AaKMpxzC4k),",
		},
		{
			name: "dbus",
			rule: dbus1,
			want: `dbus receive bus=session path=/org/gtk/vfs/metadata
       interface=org.gtk.vfs.Metadata
       member=Remove
       peer=(name=:1.15, label=tracker-extract),`,
		},
		{
			name: "dbus-bind",
			rule: &Dbus{Access: []string{"bind"}, Bus: "session", Name: "org.gnome.*"},
			want: `dbus bind bus=session name=org.gnome.*,`,
		},
		{
			name: "dbus-full",
			rule: &Dbus{Bus: "accessibility"},
			want: `dbus bus=accessibility,`,
		},
		{
			name: "file",
			rule: file1,
			want: "/usr/share/poppler/cMap/Identity-H r,",
		},
		{
			name: "link",
			rule: link3,
			want: "owner link @{user_config_dirs}/kiorc -> @{user_config_dirs}/#3954,",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.rule
			if got := r.String(); got != tt.want {
				t.Errorf("Rule.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
