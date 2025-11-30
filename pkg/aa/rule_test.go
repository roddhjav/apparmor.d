// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"reflect"
	"testing"
)

func TestRule_FromLog(t *testing.T) {
	for _, tt := range testRule {
		if tt.fromLog == nil {
			continue
		}
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fromLog(tt.log); !reflect.DeepEqual(got, tt.rule) {
				t.Errorf("RuleFromLog() = %v, want %v", got, tt.rule)
			}
		})
	}
}

func TestRule_String(t *testing.T) {
	for _, tt := range testRule {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.rule.String(); got != tt.wString {
				t.Errorf("Rule.String() = %v, want %v", got, tt.wString)
			}
		})
	}
}

func TestRule_Validate(t *testing.T) {
	for _, tt := range testRule {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.rule.Validate(); (err != nil) != tt.wValidErr {
				t.Errorf("Rules.Validate() error = %v, wantErr %v", err, tt.wValidErr)
			}
		})
	}
}

func TestRule_Compare(t *testing.T) {
	for _, tt := range testRule {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.rule.Compare(tt.other); got != tt.wCompare {
				t.Errorf("Rule.Compare() = %v, want %v", got, tt.wCompare)
			}
		})
	}
}

func TestRule_Merge(t *testing.T) {
	for _, tt := range testRule {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.rule.Merge(tt.other); got != tt.wMerge {
				t.Errorf("Rule.Merge() = %v, want %v", got, tt.wMerge)
			}
		})
	}
}

var (
	// Test cases for the Rule interface
	testRule = []struct {
		name      string
		fromLog   func(map[string]string) Rule
		log       map[string]string
		rule      Rule
		wValidErr bool
		other     Rule
		wCompare  int
		wMerge    bool
		wString   string
	}{
		{
			name:     "comment",
			rule:     comment1,
			other:    comment2,
			wCompare: 0,
			wMerge:   false,
			wString:  "#comment",
		},
		{
			name:     "abi",
			rule:     abi1,
			other:    abi2,
			wCompare: 1,
			wMerge:   false,
			wString:  "abi <abi/4.0>,",
		},
		{
			name:     "alias",
			rule:     alias1,
			other:    alias2,
			wCompare: -1,
			wMerge:   false,
			wString:  "alias /mnt/usr -> /usr,",
		},
		{
			name:     "include1",
			rule:     include1,
			other:    includeLocal1,
			wCompare: -11,
			wMerge:   false,
			wString:  "include <abstraction/base>",
		},
		{
			name:     "include2",
			rule:     include1,
			other:    include2,
			wCompare: 1,
			wMerge:   false,
			wString:  "include <abstraction/base>",
		},
		{
			name:     "include-local",
			rule:     includeLocal1,
			other:    include1,
			wCompare: 11,
			wMerge:   false,
			wString:  "include if exists <local/foo>",
		},
		{
			name:     "include-abs",
			rule:     &Include{Path: "/usr/share/apparmor.d/", IsMagic: false},
			other:    &Include{Path: "/usr/share/apparmor.d/", IsMagic: true},
			wCompare: -1,
			wMerge:   false,
			wString:  `include "/usr/share/apparmor.d/"`,
		},
		{
			name:     "variable",
			rule:     variable1,
			other:    variable2,
			wCompare: -3,
			wMerge:   false,
			wString:  "@{bin} = /{,usr/}{,s}bin",
		},
		{
			name:     "all",
			rule:     all1,
			other:    all2,
			wCompare: 0,
			wMerge:   true,
			wString:  "all,",
		},
		{
			name:     "rlimit",
			rule:     rlimit1,
			other:    rlimit2,
			wCompare: 11,
			wMerge:   false,
			wString:  "set rlimit nproc <= 200,",
		},
		{
			name:     "rlimit2",
			rule:     rlimit2,
			other:    rlimit2,
			wCompare: 0,
			wMerge:   false,
			wString:  "set rlimit cpu <= 2,",
		},
		{
			name:     "rlimit3",
			rule:     rlimit3,
			other:    rlimit1,
			wCompare: -1,
			wMerge:   false,
			wString:  "set rlimit nproc < 2,",
		},
		{
			name:     "userns",
			rule:     userns1,
			other:    userns2,
			wCompare: 1,
			wMerge:   true,
			wString:  "userns,",
		},
		{
			name:     "capbability",
			fromLog:  newCapabilityFromLog,
			log:      capability1Log,
			rule:     capability1,
			other:    capability2,
			wCompare: -5,
			wMerge:   false,
			wString:  "capability net_admin,",
		},
		{
			name:     "capability-multi",
			rule:     &Capability{Names: []string{"dac_override", "dac_read_search"}},
			other:    capability2,
			wCompare: -15,
			wMerge:   false,
			wString:  "capability dac_override dac_read_search,",
		},
		{
			name:     "capability-all",
			rule:     &Capability{},
			other:    capability2,
			wCompare: -1,
			wMerge:   false,
			wString:  "capability,",
		},
		{
			name:      "network",
			fromLog:   newNetworkFromLog,
			log:       network1Log,
			rule:      network1,
			wValidErr: true,
			other:     network2,
			wCompare:  5,
			wMerge:    false,
			wString:   "network netlink raw,",
		},
		{
			name:      "network3",
			fromLog:   newNetworkFromLog,
			log:       network3Log,
			rule:      network3,
			wValidErr: true,
			other:     network1,
			wCompare:  -7,
			wMerge:    false,
			wString:   "network dgram ip=127.0.0.1 port=57007 peer=(ip=127.0.0.53, port=53), # failed af match",
		},
		{
			name:     "mount",
			fromLog:  newMountFromLog,
			log:      mount1Log,
			rule:     mount1,
			other:    mount2,
			wCompare: 37,
			wMerge:   false,
			wString:  "mount fstype=overlay overlay -> /var/lib/docker/overlay2/opaque-bug-check1209538631/merged/, # failed perms check",
		},
		{
			name:     "remount",
			rule:     remount1,
			other:    remount2,
			wCompare: -6,
			wMerge:   false,
			wString:  "remount /,",
		},
		{
			name:     "umount",
			fromLog:  newUmountFromLog,
			log:      umount1Log,
			rule:     umount1,
			other:    umount2,
			wCompare: -8,
			wMerge:   false,
			wString:  "umount /,",
		},
		{
			name:     "pivot_root1",
			fromLog:  newPivotRootFromLog,
			log:      pivotroot1Log,
			rule:     pivotroot1,
			other:    pivotroot2,
			wCompare: -5,
			wMerge:   false,
			wString:  "pivot_root oldroot=@{run}/systemd/mount-rootfs/ @{run}/systemd/mount-rootfs/,",
		},
		{
			name:     "pivot_root2",
			rule:     pivotroot1,
			other:    pivotroot3,
			wCompare: 28,
			wMerge:   false,
			wString:  "pivot_root oldroot=@{run}/systemd/mount-rootfs/ @{run}/systemd/mount-rootfs/,",
		},
		{
			name:     "change_profile1",
			fromLog:  newChangeProfileFromLog,
			log:      changeprofile1Log,
			rule:     changeprofile1,
			other:    changeprofile2,
			wCompare: 17,
			wMerge:   false,
			wString:  "change_profile -> systemd-user,",
		},
		{
			name:     "change_profile2",
			rule:     changeprofile2,
			other:    changeprofile3,
			wCompare: -4,
			wMerge:   false,
			wString:  "change_profile -> brwap,",
		},
		{
			name:     "mqueue",
			rule:     mqueue1,
			other:    mqueue2,
			wCompare: -3,
			wMerge:   false,
			wString:  "mqueue r type=posix /,",
		},
		{
			name:     "iouring",
			rule:     iouring1,
			other:    iouring2,
			wCompare: 4,
			wMerge:   false,
			wString:  "io_uring sqpoll label=foo,",
		},
		{
			name:     "signal",
			fromLog:  newSignalFromLog,
			log:      signal1Log,
			rule:     signal1,
			other:    signal2,
			wCompare: -10,
			wMerge:   true,
			wString:  "signal receive set=kill peer=firefox//&firejail-default,",
		},
		{
			name:     "ptrace-xdg-document-portal",
			fromLog:  newPtraceFromLog,
			log:      ptrace1Log,
			rule:     ptrace1,
			other:    ptrace1,
			wCompare: 0,
			wMerge:   true,
			wString:  "ptrace read peer=nautilus,",
		},
		{
			name:     "ptrace-snap-update-ns.firefox",
			fromLog:  newPtraceFromLog,
			log:      ptrace2Log,
			rule:     ptrace2,
			other:    ptrace1,
			wCompare: 2,
			wMerge:   false,
			wString:  "ptrace readby peer=systemd-journald,",
		},
		{
			name:     "unix",
			fromLog:  newUnixFromLog,
			log:      unix1Log,
			rule:     unix1,
			other:    unix1,
			wCompare: 0,
			wMerge:   true,
			wString:  "unix (send receive) type=stream peer=(label=dbus-daemon, addr=@/tmp/dbus-AaKMpxzC4k),",
		},
		{
			name:     "dbus",
			fromLog:  newDbusFromLog,
			log:      dbus1Log,
			rule:     dbus1,
			other:    dbus1,
			wCompare: 0,
			wMerge:   true,
			wString:  "dbus receive bus=session path=/org/gtk/vfs/metadata\n       interface=org.gtk.vfs.Metadata\n       member=Remove\n       peer=(name=:1.15, label=tracker-extract),",
		},
		{
			name:     "dbus2",
			rule:     dbus2,
			other:    dbus3,
			wCompare: 9,
			wMerge:   false,
			wString:  "dbus bind bus=session name=org.gnome.evolution.dataserver.Sources5,",
		},
		{
			name:     "dbus-bind",
			rule:     &Dbus{Access: []string{"bind"}, Bus: "session", Name: "org.gnome.*"},
			other:    dbus2,
			wCompare: -39,
			wMerge:   false,
			wString:  `dbus bind bus=session name=org.gnome.*,`,
		},
		{
			name:     "dbus/full",
			rule:     &Dbus{Bus: "accessibility"},
			other:    dbus1,
			wCompare: -1,
			wMerge:   false,
			wString:  `dbus bus=accessibility,`,
		},
		{
			name:     "file",
			fromLog:  newFileFromLog,
			log:      file1Log,
			rule:     file1,
			other:    file2,
			wCompare: -14,
			wMerge:   false,
			wString:  "/usr/share/poppler/cMap/Identity-H r,",
		},
		{
			name:     "file-all",
			rule:     &File{},
			other:    &File{},
			wCompare: 0,
			wMerge:   true,
			wString:  " ,", // FIXME:
		},
		{
			name:      "file-equal",
			rule:      &File{Path: "/usr/share/poppler/cMap/Identity-H"},
			other:     &File{Path: "/usr/share/poppler/cMap/Identity-H"},
			wValidErr: true,
			wCompare:  0,
			wMerge:    true,
			wString:   "/usr/share/poppler/cMap/Identity-H ,",
		},
		{
			name:      "file-owner",
			rule:      &File{Path: "/usr/share/poppler/cMap/Identity-H", Owner: true},
			other:     &File{Path: "/usr/share/poppler/cMap/Identity-H"},
			wCompare:  1,
			wValidErr: true,
			wMerge:    false,
			wString:   "owner /usr/share/poppler/cMap/Identity-H ,",
		},
		{
			name:     "file-access",
			rule:     &File{Path: "/usr/share/poppler/cMap/Identity-H", Access: []string{"r"}},
			other:    &File{Path: "/usr/share/poppler/cMap/Identity-H", Access: []string{"w"}},
			wCompare: -5,
			wMerge:   true,
			wString:  "/usr/share/poppler/cMap/Identity-H r,",
		},
		{
			name:      "file-close",
			rule:      &File{Path: "/usr/share/poppler/cMap/"},
			other:     &File{Path: "/usr/share/poppler/cMap/Identity-H"},
			wCompare:  -10,
			wValidErr: true,
			wMerge:    false,
			wString:   "/usr/share/poppler/cMap/ ,",
		},
		{
			name:     "link1",
			fromLog:  newLinkFromLog,
			log:      link1Log,
			rule:     link1,
			other:    link2,
			wCompare: -1,
			wMerge:   false,
			wString:  "link /tmp/mkinitcpio.QDWtza/early@{lib}/firmware/i915/dg1_dmc_ver2_02.bin.zst -> /tmp/mkinitcpio.QDWtza/root@{lib}/firmware/i915/dg1_dmc_ver2_02.bin.zst,",
		},
		{
			name:     "link2",
			fromLog:  newFileFromLog,
			log:      link3Log,
			rule:     link3,
			other:    link1,
			wCompare: 1,
			wMerge:   false,
			wString:  "owner link @{user_config_dirs}/kiorc -> @{user_config_dirs}/#3954,",
		},
		{
			name:     "profile",
			rule:     profile1,
			other:    profile2,
			wCompare: -4,
			wMerge:   false,
			wString:  "profile sudo {\n}",
		},
		{
			name:     "hat",
			rule:     hat1,
			other:    hat2,
			wCompare: 3,
			wMerge:   false,
			wString:  "hat user {\n}",
		},
	}
)
