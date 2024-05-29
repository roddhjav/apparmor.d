// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"reflect"
	"testing"
)

func TestRules_FromLog(t *testing.T) {
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

func TestRules_Validate(t *testing.T) {
	for _, tt := range testRule {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.rule.Validate(); (err != nil) != tt.wValidErr {
				t.Errorf("Rules.Validate() error = %v, wantErr %v", err, tt.wValidErr)
			}
		})
	}
}

func TestRules_Less(t *testing.T) {
	for _, tt := range testRule {
		if tt.oLess == nil {
			continue
		}
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.rule.Less(tt.oLess); got != tt.wLessErr {
				t.Errorf("Rule.Less() = %v, want %v", got, tt.wLessErr)
			}
		})
	}
}

func TestRules_Equals(t *testing.T) {
	for _, tt := range testRule {
		if tt.oEqual == nil {
			continue
		}
		t.Run(tt.name, func(t *testing.T) {
			r := tt.rule
			if got := r.Equals(tt.oEqual); got != tt.wEqualErr {
				t.Errorf("Rule.Equals() = %v, want %v", got, tt.wEqualErr)
			}
		})
	}
}

func TestRules_String(t *testing.T) {
	for _, tt := range testRule {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.rule.String(); got != tt.wString {
				t.Errorf("Rule.String() = %v, want %v", got, tt.wString)
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
		oLess     Rule
		wLessErr  bool
		oEqual    Rule
		wEqualErr bool
		wString   string
	}{
		{
			name:      "comment",
			rule:      comment1,
			oLess:     comment2,
			wLessErr:  false,
			oEqual:    comment2,
			wEqualErr: false,
			wString:   "#comment",
		},
		{
			name:      "abi",
			rule:      abi1,
			oLess:     abi2,
			wLessErr:  false,
			oEqual:    abi1,
			wEqualErr: true,
			wString:   "abi <abi/4.0>,",
		},
		{
			name:      "alias",
			rule:      alias1,
			oLess:     alias2,
			wLessErr:  true,
			oEqual:    alias2,
			wEqualErr: false,
			wString:   "alias /mnt/usr -> /usr,",
		},
		{
			name:      "include1",
			rule:      include1,
			oLess:     includeLocal1,
			wLessErr:  false,
			oEqual:    includeLocal1,
			wEqualErr: false,
			wString:   "include <abstraction/base>",
		},
		{
			name:     "include2",
			rule:     include1,
			oLess:    include2,
			wLessErr: false,
			wString:  "include <abstraction/base>",
		},
		{
			name:     "include-local",
			rule:     includeLocal1,
			oLess:    include1,
			wLessErr: true,
			wString:  "include if exists <local/foo>",
		},
		{
			name:    "include/abs",
			rule:    &Include{Path: "/usr/share/apparmor.d/", IsMagic: false},
			wString: `include "/usr/share/apparmor.d/"`,
		},
		{
			name:      "variable",
			rule:      variable1,
			oLess:     variable2,
			wLessErr:  true,
			oEqual:    variable1,
			wEqualErr: true,
			wString:   "@{bin} = /{,usr/}{,s}bin",
		},
		{
			name:      "all",
			rule:      all1,
			oLess:     all2,
			wLessErr:  false,
			oEqual:    all2,
			wEqualErr: false,
			wString:   "all,",
		},
		{
			name:      "rlimit",
			rule:      rlimit1,
			oLess:     rlimit2,
			wLessErr:  false,
			oEqual:    rlimit1,
			wEqualErr: true,
			wString:   "set rlimit nproc <= 200,",
		},
		{
			name:     "rlimit2",
			rule:     rlimit2,
			oLess:    rlimit2,
			wLessErr: false,
			wString:  "set rlimit cpu <= 2,",
		},
		{
			name:     "rlimit3",
			rule:     rlimit3,
			oLess:    rlimit1,
			wLessErr: true,

			wString: "set rlimit nproc < 2,",
		},
		{
			name:      "userns",
			rule:      userns1,
			oLess:     userns2,
			wLessErr:  true,
			oEqual:    userns1,
			wEqualErr: true,
			wString:   "userns,",
		},
		{
			name:      "capbability",
			fromLog:   newCapabilityFromLog,
			log:       capability1Log,
			rule:      capability1,
			oLess:     capability2,
			wLessErr:  true,
			oEqual:    capability1,
			wEqualErr: true,
			wString:   "capability net_admin,",
		},
		{
			name:    "capability/multi",
			rule:    &Capability{Names: []string{"dac_override", "dac_read_search"}},
			wString: "capability dac_override dac_read_search,",
		},
		{
			name:    "capability/all",
			rule:    &Capability{},
			wString: "capability,",
		},
		{
			name:      "network",
			fromLog:   newNetworkFromLog,
			log:       network1Log,
			rule:      network1,
			wValidErr: true,
			oLess:     network2,
			wLessErr:  false,
			oEqual:    network1,
			wEqualErr: true,
			wString:   "network netlink raw,",
		},
		{
			name:      "mount",
			fromLog:   newMountFromLog,
			log:       mount1Log,
			rule:      mount1,
			oEqual:    mount2,
			wEqualErr: false,
			wString:   "mount fstype=overlay overlay -> /var/lib/docker/overlay2/opaque-bug-check1209538631/merged/,  # failed perms check",
		},
		{
			name:      "remount",
			rule:      remount1,
			oLess:     remount2,
			wLessErr:  true,
			oEqual:    remount1,
			wEqualErr: true,
			wString:   "remount /,",
		},
		{
			name:      "umount",
			fromLog:   newUmountFromLog,
			log:       umount1Log,
			rule:      umount1,
			oLess:     umount2,
			wLessErr:  true,
			oEqual:    umount1,
			wEqualErr: true,
			wString:   "umount /,",
		},
		{
			name:      "pivot_root1",
			fromLog:   newPivotRootFromLog,
			log:       pivotroot1Log,
			rule:      pivotroot1,
			oLess:     pivotroot2,
			wLessErr:  false,
			oEqual:    pivotroot2,
			wEqualErr: false,
			wString:   "pivot_root oldroot=@{run}/systemd/mount-rootfs/ @{run}/systemd/mount-rootfs/,",
		},
		{
			name:     "pivot_root2",
			rule:     pivotroot1,
			oLess:    pivotroot3,
			wLessErr: false,
			wString:  "pivot_root oldroot=@{run}/systemd/mount-rootfs/ @{run}/systemd/mount-rootfs/,",
		},
		{
			name:     "change_profile1",
			fromLog:  newChangeProfileFromLog,
			log:      changeprofile1Log,
			rule:     changeprofile1,
			oLess:    changeprofile2,
			wLessErr: false,
			wString:  "change_profile -> systemd-user,",
		},
		{
			name:      "change_profile2",
			rule:      changeprofile2,
			oLess:     changeprofile3,
			wLessErr:  true,
			oEqual:    changeprofile1,
			wEqualErr: false,
			wString:   "change_profile -> brwap,",
		},
		{
			name:      "mqueue",
			rule:      mqueue1,
			oLess:     mqueue2,
			wLessErr:  true,
			oEqual:    mqueue1,
			wEqualErr: true,
			wString:   "mqueue r type=posix /,",
		},
		{
			name:      "iouring",
			rule:      iouring1,
			oLess:     iouring2,
			wLessErr:  false,
			oEqual:    iouring2,
			wEqualErr: false,
			wString:   "io_uring sqpoll label=foo,",
		},
		{
			name:      "signal",
			fromLog:   newSignalFromLog,
			log:       signal1Log,
			rule:      signal1,
			oLess:     signal2,
			wLessErr:  false,
			oEqual:    signal1,
			wEqualErr: true,
			wString:   "signal receive set=kill peer=firefox//&firejail-default,",
		},
		{
			name:      "ptrace/xdg-document-portal",
			fromLog:   newPtraceFromLog,
			log:       ptrace1Log,
			rule:      ptrace1,
			oLess:     ptrace2,
			wLessErr:  false,
			oEqual:    ptrace1,
			wEqualErr: true,
			wString:   "ptrace read peer=nautilus,",
		},
		{
			name:      "ptrace/snap-update-ns.firefox",
			fromLog:   newPtraceFromLog,
			log:       ptrace2Log,
			rule:      ptrace2,
			oLess:     ptrace1,
			wLessErr:  false,
			oEqual:    ptrace1,
			wEqualErr: false,
			wString:   "ptrace readby peer=systemd-journald,",
		},
		{
			name:      "unix",
			fromLog:   newUnixFromLog,
			log:       unix1Log,
			rule:      unix1,
			oLess:     unix1,
			wLessErr:  false,
			oEqual:    unix1,
			wEqualErr: true,
			wString:   "unix (send receive) type=stream protocol=0 addr=none peer=(label=dbus-daemon, addr=@/tmp/dbus-AaKMpxzC4k),",
		},
		{
			name:      "dbus",
			fromLog:   newDbusFromLog,
			log:       dbus1Log,
			rule:      dbus1,
			oLess:     dbus1,
			wLessErr:  false,
			oEqual:    dbus2,
			wEqualErr: false,
			wString:   "dbus receive bus=session path=/org/gtk/vfs/metadata\n       interface=org.gtk.vfs.Metadata\n       member=Remove\n       peer=(name=:1.15, label=tracker-extract),",
		},
		{
			name:     "dbus2",
			rule:     dbus2,
			oLess:    dbus3,
			wLessErr: false,
			wString:  "dbus bind bus=session name=org.gnome.evolution.dataserver.Sources5,",
		},
		{
			name:    "dbus/bind",
			rule:    &Dbus{Access: []string{"bind"}, Bus: "session", Name: "org.gnome.*"},
			wString: `dbus bind bus=session name=org.gnome.*,`,
		},
		{
			name:    "dbus/full",
			rule:    &Dbus{Bus: "accessibility"},
			wString: `dbus bus=accessibility,`,
		},
		{
			name:      "file",
			fromLog:   newFileFromLog,
			log:       file1Log,
			rule:      file1,
			oLess:     file2,
			wLessErr:  true,
			oEqual:    file2,
			wEqualErr: false,
			wString:   "/usr/share/poppler/cMap/Identity-H r,",
		},
		{
			name:     "file/empty",
			rule:     &File{},
			oLess:    &File{},
			wLessErr: false,
			wString:  " ,",
		},
		{
			name:     "file/equal",
			rule:     &File{Path: "/usr/share/poppler/cMap/Identity-H"},
			oLess:    &File{Path: "/usr/share/poppler/cMap/Identity-H"},
			wLessErr: false,
			wString:  "/usr/share/poppler/cMap/Identity-H ,",
		},
		{
			name:     "file/owner",
			rule:     &File{Path: "/usr/share/poppler/cMap/Identity-H", Owner: true},
			oLess:    &File{Path: "/usr/share/poppler/cMap/Identity-H"},
			wLessErr: true,
			wString:  "owner /usr/share/poppler/cMap/Identity-H ,",
		},
		{
			name:     "file/access",
			rule:     &File{Path: "/usr/share/poppler/cMap/Identity-H", Access: []string{"r"}},
			oLess:    &File{Path: "/usr/share/poppler/cMap/Identity-H", Access: []string{"w"}},
			wLessErr: false,
			wString:  "/usr/share/poppler/cMap/Identity-H r,",
		},
		{
			name:     "file/close",
			rule:     &File{Path: "/usr/share/poppler/cMap/"},
			oLess:    &File{Path: "/usr/share/poppler/cMap/Identity-H"},
			wLessErr: true,
			wString:  "/usr/share/poppler/cMap/ ,",
		},
		{
			name:      "link",
			fromLog:   newLinkFromLog,
			log:       link1Log,
			rule:      link1,
			oLess:     link2,
			wLessErr:  true,
			oEqual:    link3,
			wEqualErr: false,
			wString:   "link /tmp/mkinitcpio.QDWtza/early@{lib}/firmware/i915/dg1_dmc_ver2_02.bin.zst -> /tmp/mkinitcpio.QDWtza/root@{lib}/firmware/i915/dg1_dmc_ver2_02.bin.zst,",
		},
		{
			name:    "link",
			fromLog: newFileFromLog,
			log:     link3Log,
			rule:    link3,
			wString: "owner link @{user_config_dirs}/kiorc -> @{user_config_dirs}/#3954,",
		},
		{
			name:      "profile",
			rule:      profile1,
			oLess:     profile2,
			wLessErr:  true,
			oEqual:    profile1,
			wEqualErr: true,
			wString:   "profile sudo {\n}",
		},
		{
			name:      "hat",
			rule:      hat1,
			oLess:     hat2,
			wLessErr:  false,
			oEqual:    hat1,
			wEqualErr: true,
			wString:   "hat user {\n}",
		},
	}
)
