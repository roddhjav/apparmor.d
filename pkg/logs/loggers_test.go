// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2023 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package logs

import (
	"reflect"
	"testing"
)

func TestGetJournalctlLogs(t *testing.T) {
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
					"label":      "gsd-xsettings",
					"operation":  "dbus_method_call",
					"name":       ":*",
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
			reader, _ := GetJournalctlLogs(tt.path, tt.useFile)
			if got := NewApparmorLogs(reader, tt.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewApparmorLogs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSelectLogFile(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{
			name: "Get audit.log",
			path: "../../tests/audit.log",
			want: "../../tests/audit.log",
		},
		{
			name: "Get /var/log/audit/audit.log.1",
			path: "1",
			want: "/var/log/audit/audit.log.1",
		},
		{
			name: "Get default log file",
			path: "",
			want: "/var/log/audit/audit.log",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SelectLogFile(tt.path); got != tt.want {
				t.Errorf("SelectLogFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
