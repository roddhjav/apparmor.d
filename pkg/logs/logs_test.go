// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package logs

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/roddhjav/apparmor.d/pkg/aa"
)

var (
	refKmod = AppArmorLogs{
		{
			"apparmor":       "ALLOWED",
			"profile":        "kmod",
			"operation":      "file_inherit",
			"comm":           "modprobe",
			"family":         "unix",
			"sock_type":      "stream",
			"protocol":       "0",
			"requested_mask": "send receive",
			"class":          "net",
		},
	}
	refMan = AppArmorLogs{
		{
			"apparmor":       "ALLOWED",
			"profile":        "man",
			"operation":      "exec",
			"name":           "@{bin}/preconv",
			"target":         "man_groff",
			"info":           "no new privs",
			"comm":           "man",
			"requested_mask": "x",
			"denied_mask":    "x",
			"error":          "-1",
			"fsuid":          "1000",
			"ouid":           "1000",
			"FSUID":          "user",
			"OUID":           "user",
		},
	}
	refPowerProfiles = AppArmorLogs{
		{
			"apparmor":   "ALLOWED",
			"label":      "power-profiles-daemon",
			"operation":  "dbus_method_call",
			"name":       "org.freedesktop.DBus",
			"mask":       "send",
			"bus":        "system",
			"path":       "/org/freedesktop/DBus",
			"interface":  "org.freedesktop.DBus",
			"member":     "AddMatch",
			"peer_label": "dbus-daemon",
			"exe":        "/usr/bin/dbus-daemon",
			"sauid":      "102",
			"hostname":   "?",
			"addr":       "?",
			"terminal":   "?",
			"UID":        "messagebus",
			"AUID":       "unset",
			"SAUID":      "messagebus",
		},
	}
)

func TestAppArmorEvents(t *testing.T) {
	tests := []struct {
		name  string
		event string
		want  AppArmorLogs
	}{
		{
			name:  "event_audit_1",
			event: `type=AVC msg=audit(1345027352.096:499): apparmor="ALLOWED" operation="rename_dest" parent=6974 profile="/usr/sbin/httpd2-prefork//vhost_foo" name=2F686F6D652F7777772F666F6F2E6261722E696E2F68747470646F63732F61707061726D6F722F696D616765732F746573742F696D61676520312E6A7067 pid=20143 comm="httpd2-prefork" requested_mask="wc" denied_mask="wc" fsuid=30 ouid=30`,
			want: AppArmorLogs{
				{
					"apparmor":       "ALLOWED",
					"profile":        "@{sbin}/httpd2-prefork//vhost_foo",
					"operation":      "rename_dest",
					"name":           "@{HOME}/foo.bar.in/httpdocs/apparmor/images/test/image 1.jpg",
					"comm":           "httpd2-prefork",
					"requested_mask": "wc",
					"denied_mask":    "wc",
					"parent":         "6974",
					"fsuid":          "30",
					"ouid":           "30",
				},
			},
		},
		{
			name:  "event_audit_2",
			event: `type=AVC msg=audit(1322614918.292:4376): apparmor="ALLOWED" operation="file_perm" parent=16001 profile=666F6F20626172 name="/home/foo/.bash_history" pid=17011 comm="bash" requested_mask="rw" denied_mask="rw" fsuid=0 ouid=1000`,
			want: AppArmorLogs{
				{
					"apparmor":       "ALLOWED",
					"profile":        "foo bar",
					"operation":      "file_perm",
					"name":           "@{HOME}/.bash_history",
					"comm":           "bash",
					"requested_mask": "rw",
					"denied_mask":    "rw",
					"parent":         "16001",
					"fsuid":          "0",
					"ouid":           "1000",
				},
			},
		},
		{
			name:  "disconnected_path",
			event: `type=AVC msg=audit(1424425690.883:716630): apparmor="ALLOWED" operation="file_mmap" info="Failed name lookup - disconnected path" error=-13 profile="/sbin/klogd" name="var/run/nscd/passwd" pid=25333 comm="id" requested_mask="r" denied_mask="r" fsuid=1002 ouid=0`,
			want: AppArmorLogs{
				{
					"apparmor":       "ALLOWED",
					"profile":        "/sbin/klogd",
					"operation":      "file_mmap",
					"name":           "var@{run}/nscd/passwd",
					"comm":           "id",
					"info":           "Failed name lookup - disconnected path",
					"requested_mask": "r",
					"denied_mask":    "r",
					"error":          "-13",
					"fsuid":          "1002",
					"ouid":           "0",
				},
			},
		},
		{
			name:  "dbus_system",
			event: `type=USER_AVC msg=audit(1111111111.111:1111): pid=1780 uid=102 auid=4294967295 ses=4294967295 subj=? msg='apparmor="ALLOWED" operation="dbus_method_call" bus="system" path="/org/freedesktop/PolicyKit1/Authority" interface="org.freedesktop.PolicyKit1.Authority" member="CheckAuthorization" mask="send" name="org.freedesktop.PolicyKit1" pid=1794 label="snapd" peer_pid=1790 peer_label="polkitd"  exe="/usr/bin/dbus-daemon" sauid=102 hostname=? addr=? terminal=? UID="messagebus" AUID="unset" SAUID="messagebus"`,
			want: AppArmorLogs{
				{
					"apparmor":   "ALLOWED",
					"label":      "snapd",
					"operation":  "dbus_method_call",
					"name":       "org.freedesktop.PolicyKit1",
					"mask":       "send",
					"bus":        "system",
					"path":       "/org/freedesktop/PolicyKit1/Authority",
					"interface":  "org.freedesktop.PolicyKit1.Authority",
					"member":     "CheckAuthorization",
					"peer_label": "polkitd",
					"exe":        "/usr/bin/dbus-daemon",
					"sauid":      "102",
					"hostname":   "?",
					"addr":       "?",
					"terminal":   "?",
					"UID":        "messagebus",
					"AUID":       "unset",
					"SAUID":      "messagebus",
				},
			},
		},
		{
			name:  "dbus_session",
			event: `apparmor="ALLOWED" operation="dbus_bind"  bus="session" name="org.freedesktop.portal.Documents" mask="bind" pid=2174 label="xdg-document-portal"`,
			want: AppArmorLogs{
				{
					"apparmor":  "ALLOWED",
					"label":     "xdg-document-portal",
					"operation": "dbus_bind",
					"name":      "org.freedesktop.portal.Documents",
					"mask":      "bind",
					"bus":       "session",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := strings.NewReader(tt.event)
			if got := New(file, "", ""); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name      string
		namespace string
		path      string
		want      AppArmorLogs
	}{
		{
			name: "dnsmasq",
			path: filepath.Join(testdata, "audit.log"),
			want: AppArmorLogs{
				{
					"apparmor":       "DENIED",
					"profile":        "dnsmasq",
					"operation":      "open",
					"name":           "@{PROC}/sys/kernel/osrelease",
					"comm":           "dnsmasq",
					"requested_mask": "r",
					"denied_mask":    "r",
					"fsuid":          "0",
					"ouid":           "0",
					"FSUID":          "root",
					"OUID":           "root",
				},
				{
					"apparmor":       "DENIED",
					"profile":        "dnsmasq",
					"operation":      "open",
					"name":           "@{PROC}/1/environ",
					"comm":           "dnsmasq",
					"requested_mask": "r",
					"denied_mask":    "r",
					"fsuid":          "0",
					"ouid":           "0",
					"FSUID":          "root",
					"OUID":           "root",
				},
				{
					"apparmor":       "DENIED",
					"profile":        "dnsmasq",
					"operation":      "open",
					"name":           "@{PROC}/cmdline",
					"comm":           "dnsmasq",
					"requested_mask": "r",
					"denied_mask":    "r",
					"fsuid":          "0",
					"ouid":           "0",
					"FSUID":          "root",
					"OUID":           "root",
				},
			},
		},
		{
			name: "kmod",
			path: filepath.Join(testdata, "audit.log"),
			want: refKmod,
		},
		{
			name: "man",
			path: filepath.Join(testdata, "audit.log"),
			want: refMan,
		},
		{
			name: "power-profiles-daemon",
			path: filepath.Join(testdata, "audit.log"),
			want: refPowerProfiles,
		},
		{
			name: "signal-desktop",
			path: filepath.Join(testdata, "audit.log"),
			want: AppArmorLogs{
				{
					"apparmor":       "ALLOWED",
					"profile":        "signal-desktop",
					"operation":      "open",
					"class":          "file",
					"name":           "@{sys}/devices/@{pci}/boot_vga",
					"comm":           "signal-desktop",
					"requested_mask": "r",
					"denied_mask":    "r",
					"fsuid":          "1000",
					"ouid":           "0",
					"FSUID":          "user",
					"OUID":           "root",
				},
			},
		},
		{
			name: "startplasma",
			path: filepath.Join(testdata, "audit.log"),
			want: AppArmorLogs{
				{
					"apparmor":       "ALLOWED",
					"operation":      "link",
					"class":          "file",
					"profile":        "startplasma",
					"name":           "@{user_cache_dirs}/ksycoca5_de_LQ6f0J2qZg4vOKgw2NbXuW7iuVU=.isNSBz",
					"target":         "@{user_cache_dirs}/#@{int}",
					"comm":           "startplasma-way",
					"denied_mask":    "k",
					"requested_mask": "k",
					"fsuid":          "1000",
					"ouid":           "1000",
					"FSUID":          "user",
					"OUID":           "user",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, _ := os.Open(tt.path)
			if got := New(file, tt.name, tt.namespace); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoad(t *testing.T) {
	tests := []struct {
		name      string
		namespace string
		path      string
		want      AppArmorLogs
	}{
		{
			name: "dnsmasq",
			path: filepath.Join(testdata, "aa-log"),
			want: AppArmorLogs{
				{
					"apparmor":       "DENIED",
					"profile":        "dnsmasq",
					"operation":      "open",
					"name":           "@{PROC}/sys/kernel/osrelease",
					"comm":           "dnsmasq",
					"requested_mask": "r",
					"denied_mask":    "r",
				},
				{
					"apparmor":       "DENIED",
					"profile":        "dnsmasq",
					"operation":      "open",
					"name":           "@{PROC}/1/environ",
					"comm":           "dnsmasq",
					"requested_mask": "r",
					"denied_mask":    "r",
				},
				{
					"apparmor":       "DENIED",
					"profile":        "dnsmasq",
					"operation":      "open",
					"name":           "@{PROC}/cmdline",
					"comm":           "dnsmasq",
					"requested_mask": "r",
					"denied_mask":    "r",
				},
			},
		},
		{
			name: "kmod",
			path: filepath.Join(testdata, "aa-log"),
			want: refKmod,
		},
		{
			name: "man",
			path: filepath.Join(testdata, "aa-log"),
			want: refMan,
		},
		{
			name: "power-profiles-daemon",
			path: filepath.Join(testdata, "aa-log"),
			want: AppArmorLogs{
				{
					"addr":       "?",
					"apparmor":   "ALLOWED",
					"bus":        "system",
					"interface":  "org.freedesktop.DBus",
					"mask":       "send",
					"member":     "AddMatch",
					"name":       "org.freedesktop.DBus",
					"operation":  "dbus_method_call",
					"path":       "/org/freedesktop/DBus",
					"peer_label": "dbus-daemon",
					"profile":    "power-profiles-daemon",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, _ := os.Open(tt.path)
			if got := Load(file, tt.name, tt.namespace); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Load() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppArmorLogs_String(t *testing.T) {
	tests := []struct {
		name   string
		aaLogs AppArmorLogs
		want   string
	}{
		{
			name:   "kmod",
			aaLogs: refKmod,
			want:   "\033[1;32mALLOWED\033[0m \033[34mkmod\033[0m \033[33mfile_inherit\033[0m comm=modprobe family=unix sock_type=stream protocol=0 requested_mask=\033[1;31m\"send receive\"\033[0m\n",
		},
		{
			name:   "man",
			aaLogs: refMan,
			want:   "\033[1;32mALLOWED\033[0m \033[34mman\033[0m \033[33mexec\033[0m\033[35m owner\033[0m \033[35m@{bin}/preconv\033[0m -> \033[35mman_groff\033[0m info=\"no new privs\" comm=man requested_mask=\033[1;31mx\033[0m denied_mask=\033[1;31mx\033[0m error=-1\n",
		},
		{
			name:   "power-profiles-daemon",
			aaLogs: refPowerProfiles,
			want:   "\033[1;32mALLOWED\033[0m \033[34mpower-profiles-daemon\033[0m \033[33mdbus_method_call\033[0m \033[35morg.freedesktop.DBus\033[0m \033[1;31msend\033[0m \033[36mbus=system\033[0m path=\033[37m/org/freedesktop/DBus\033[0m interface=\033[37morg.freedesktop.DBus\033[0m member=\033[32mAddMatch\033[0m peer_label=dbus-daemon addr=?\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.aaLogs.String(); got != tt.want {
				t.Errorf("AppArmorLogs.String() = %v, want %v len: %d - %d", got, tt.want, len(got), len(tt.want))
			}
		})
	}
}

func TestAppArmorLogs_ParseToProfiles(t *testing.T) {
	tests := []struct {
		name   string
		aaLogs AppArmorLogs
		want   map[string]*aa.Profile
	}{
		{
			name:   "",
			aaLogs: append(append(refKmod, refPowerProfiles...), refKmod...),
			want: map[string]*aa.Profile{
				"kmod": {
					Header: aa.Header{Name: "kmod"},
					Rules: aa.Rules{
						&aa.Unix{
							Base:     aa.Base{FileInherit: true},
							Access:   []string{"send", "receive"},
							Type:     "stream",
							Protocol: "0",
						},
						&aa.Unix{
							Base:     aa.Base{FileInherit: true},
							Access:   []string{"send", "receive"},
							Type:     "stream",
							Protocol: "0",
						},
					},
				},
				"power-profiles-daemon": {
					Header: aa.Header{Name: "power-profiles-daemon"},
					Rules: aa.Rules{
						&aa.Dbus{
							Access:    []string{"send"},
							Bus:       "system",
							Path:      "/org/freedesktop/DBus",
							Interface: "org.freedesktop.DBus",
							Member:    "AddMatch",
							PeerName:  "org.freedesktop.DBus",
							PeerLabel: "dbus-daemon",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.aaLogs.ParseToProfiles(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AppArmorLogs.ParseToProfiles() = %v, want %v", got, tt.want)
			}
		})
	}
}
