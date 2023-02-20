// aa-log - Review AppArmor generated messages
// Copyright (C) 2021 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package main

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

var refKmod = AppArmorLogs{
	{
		"apparmor":       "ALLOWED",
		"profile":        "kmod",
		"operation":      "file_inherit",
		"comm":           "modprobe",
		"family":         "unix",
		"sock_type":      "stream",
		"protocol":       "0",
		"requested_mask": "send receive",
	},
}

var refMan = AppArmorLogs{
	{
		"apparmor":       "ALLOWED",
		"profile":        "man",
		"operation":      "exec",
		"name":           "/usr/bin/preconv",
		"info":           "no new privs",
		"comm":           "man",
		"requested_mask": "x",
		"denied_mask":    "x",
		"error":          "-1",
	},
}

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
					"profile":        "/usr/sbin/httpd2-prefork//vhost_foo",
					"operation":      "rename_dest",
					"name":           "/home/www/foo.bar.in/httpdocs/apparmor/images/test/image 1.jpg",
					"comm":           "httpd2-prefork",
					"requested_mask": "wc",
					"denied_mask":    "wc",
					"parent":         "6974",
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
					"name":           "/home/foo/.bash_history",
					"comm":           "bash",
					"requested_mask": "rw",
					"denied_mask":    "rw",
					"parent":         "16001",
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
					"name":           "var/run/nscd/passwd",
					"comm":           "id",
					"info":           "Failed name lookup - disconnected path",
					"requested_mask": "r",
					"denied_mask":    "r",
					"error":          "-13",
				},
			},
		},
		{
			name:  "dbus_system",
			event: `type=USER_AVC msg=audit(1111111111.111:1111): pid=1780 uid=102 auid=4294967295 ses=4294967295 subj=? msg='apparmor="ALLOWED" operation="dbus_method_call" bus="system" path="/org/freedesktop/PolicyKit1/Authority" interface="org.freedesktop.PolicyKit1.Authority" member="CheckAuthorization" mask="send" name="org.freedesktop.PolicyKit1" pid=1794 label="snapd" peer_pid=1790 peer_label="polkitd"  exe="/usr/bin/dbus-daemon" sauid=102 hostname=? addr=? terminal=?'UID="messagebus" AUID="unset" SAUID="messagebus"`,
			want: AppArmorLogs{
				{
					"apparmor":   "ALLOWED",
					"profile":    "",
					"label":      "snapd",
					"operation":  "dbus_method_call",
					"name":       "org.freedesktop.PolicyKit1",
					"mask":       "send",
					"bus":        "system",
					"path":       "/org/freedesktop/PolicyKit1/Authority",
					"interface":  "org.freedesktop.PolicyKit1.Authority",
					"member":     "CheckAuthorization",
					"peer_label": "polkitd",
				},
			},
		},
		{
			name:  "dbus_session",
			event: `apparmor="ALLOWED" operation="dbus_bind"  bus="session" name="org.freedesktop.portal.Documents" mask="bind" pid=2174 label="xdg-document-portal"`,
			want: AppArmorLogs{
				{
					"apparmor":  "ALLOWED",
					"profile":   "",
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
			if got := NewApparmorLogs(file, ""); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewApparmorLogs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewApparmorLogs(t *testing.T) {
	tests := []struct {
		name string
		path string
		want AppArmorLogs
	}{
		{
			name: "dnsmasq",
			path: "../../tests/audit.log",
			want: AppArmorLogs{
				{
					"apparmor":       "DENIED",
					"profile":        "dnsmasq",
					"operation":      "open",
					"name":           "/proc/sys/kernel/osrelease",
					"comm":           "dnsmasq",
					"requested_mask": "r",
					"denied_mask":    "r",
				},
				{
					"apparmor":       "DENIED",
					"profile":        "dnsmasq",
					"operation":      "open",
					"name":           "/proc/1/environ",
					"comm":           "dnsmasq",
					"requested_mask": "r",
					"denied_mask":    "r",
				},
				{
					"apparmor":       "DENIED",
					"profile":        "dnsmasq",
					"operation":      "open",
					"name":           "/proc/cmdline",
					"comm":           "dnsmasq",
					"requested_mask": "r",
					"denied_mask":    "r",
				},
			},
		},
		{
			name: "kmod",
			path: "../../tests/audit.log",
			want: refKmod,
		},
		{
			name: "man",
			path: "../../tests/audit.log",
			want: refMan,
		},
		{
			name: "power-profiles-daemon",
			path: "../../tests/audit.log",
			want: AppArmorLogs{
				{
					"apparmor":   "ALLOWED",
					"profile":    "",
					"label":      "power-profiles-daemon",
					"operation":  "dbus_method_call",
					"name":       "org.freedesktop.DBus",
					"mask":       "send",
					"bus":        "system",
					"path":       "/org/freedesktop/DBus",
					"interface":  "org.freedesktop.DBus",
					"member":     "AddMatch",
					"peer_label": "dbus-daemon",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, _ := os.Open(tt.path)
			if got := NewApparmorLogs(file, tt.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewApparmorLogs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getJournalctlLogs(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		useFile bool
		want    AppArmorLogs
	}{
		{
			name:    "gsd-xsettings",
			useFile: true,
			path:    "../../tests/systemd.log",
			want: AppArmorLogs{
				{
					"apparmor":   "ALLOWED",
					"profile":    "",
					"label":      "gsd-xsettings",
					"operation":  "dbus_method_call",
					"name":       ":1.88",
					"mask":       "receive",
					"bus":        "session",
					"path":       "/org/gtk/Settings",
					"interface":  "org.freedesktop.DBus.Properties",
					"member":     "GetAll",
					"peer_label": "gnome-extension-ding",
				},
			},
		},
		{
			name:    "journalctl",
			useFile: false,
			path:    "",
			want:    AppArmorLogs{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader, _ := getJournalctlLogs(tt.path, tt.useFile)
			if got := NewApparmorLogs(reader, tt.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewApparmorLogs() = %v, want %v", got, tt.want)
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
			want:   "\033[1;32mALLOWED\033[0m \033[34mman\033[0m \033[33mexec\033[0m \033[35m/usr/bin/preconv\033[0m info=\"no new privs\" comm=man requested_mask=\033[1;31mx\033[0m denied_mask=\033[1;31mx\033[0m error=-1\n",
		},
		{
			name: "power-profiles-daemon",
			aaLogs: AppArmorLogs{
				{
					"apparmor":   "ALLOWED",
					"profile":    "",
					"label":      "power-profiles-daemon",
					"operation":  "dbus_method_call",
					"name":       "org.freedesktop.DBus",
					"mask":       "send",
					"bus":        "system",
					"path":       "/org/freedesktop/DBus",
					"interface":  "org.freedesktop.DBus",
					"member":     "AddMatch",
					"peer_label": "dbus-daemon",
				},
			},
			want: "\033[1;32mALLOWED\033[0m \033[34mpower-profiles-daemon\033[0m \033[33mdbus_method_call\033[0m \033[35morg.freedesktop.DBus\033[0m \033[1;31msend\033[0m \033[36mbus=system\033[0m path=\033[37m/org/freedesktop/DBus\033[0m interface=\033[37morg.freedesktop.DBus\033[0m member=\033[32mAddMatch\033[0m peer_label=dbus-daemon\n",
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

func Test_app(t *testing.T) {
	tests := []struct {
		name    string
		logger  string
		path    string
		profile string
		wantErr bool
	}{
		{
			name:    "Test audit.log",
			logger:  "auditd",
			path:    "../../tests/audit.log",
			profile: "",
			wantErr: false,
		},
		{
			name:    "Test Dbus Session",
			logger:  "systemd",
			path:    "../../tests/systemd.log",
			profile: "",
			wantErr: false,
		},
		{
			name:    "No logfile",
			logger:  "auditd",
			path:    "../../tests/log",
			profile: "",
			wantErr: true,
		},
		{
			name:    "Logger not supported",
			logger:  "raw",
			path:    "../../tests/audit.log",
			profile: "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := aaLog(tt.logger, tt.path, tt.profile); (err != nil) != tt.wantErr {
				t.Errorf("aaLog() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
