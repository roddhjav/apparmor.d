// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package logs

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

var (
	testdata = "../../tests/testdata/logs"
)

func TestGetJournalctlLogs(t *testing.T) {
	tests := []struct {
		name      string
		namespace string
		path      string
		useFile   bool
		want      AppArmorLogs
	}{
		{
			name:    "gsd-xsettings",
			useFile: true,
			path:    filepath.Join(testdata, "systemd.log"),
			want: AppArmorLogs{
				{
					"apparmor":   "ALLOWED",
					"label":      "gsd-xsettings",
					"operation":  "dbus_method_call",
					"name":       "@{busname}",
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
			reader, _ := GetJournalctlLogs(tt.path, "", "", tt.useFile)
			if got := New(reader, tt.name, tt.namespace); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSelectLogFile(t *testing.T) {
	canReadPath := func(path string) bool {
		if _, err := os.Stat(path); err == nil {
			if file, err := os.Open(path); err == nil {
				if err := file.Close(); err != nil {
					return false
				}
				return true
			}
		}
		return false
	}

	tests := []struct {
		name    string
		path    string
		want    string
		wantErr bool
	}{
		{
			name:    "Get audit.log",
			path:    filepath.Join(testdata, "audit.log"),
			want:    filepath.Join(testdata, "audit.log"),
			wantErr: false,
		},
		{
			name:    "Get /var/log/audit/audit.log.1",
			path:    "1",
			want:    "/var/log/audit/audit.log.1",
			wantErr: !canReadPath("/var/log/audit/audit.log.1"),
		},
		{
			name:    "Get default log file",
			path:    "",
			want:    "/var/log/audit/audit.log",
			wantErr: !canReadPath("/var/log/audit/audit.log.1"),
		},
		{
			name:    "File not found",
			path:    "/nonexistent/file",
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SelectLogFile(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("SelectLogFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SelectLogFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
