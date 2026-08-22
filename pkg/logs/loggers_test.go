// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package logs

import (
	"path/filepath"
	"reflect"
	"syscall"
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
			name:    "gnome-clocks",
			useFile: true,
			path:    filepath.Join(testdata, "systemd.log"),
			want: AppArmorLogs{
				{
					"apparmor":   "DENIED",
					"label":      "gnome-clocks",
					"operation":  "dbus_method_call",
					"name":       "org.freedesktop.DBus",
					"mask":       "send",
					"bus":        "session",
					"path":       "/org/freedesktop/DBus",
					"interface":  "org.freedesktop.DBus",
					"member":     "ListActivatableNames",
					"peer_label": "dbus-session",
				},
			},
		},
		{
			name:    "gsd-media-keys",
			useFile: true,
			path:    filepath.Join(testdata, "systemd.log"),
			want: AppArmorLogs{
				{
					"apparmor":   "DENIED",
					"label":      "gsd-media-keys",
					"operation":  "dbus_method_call",
					"name":       "@{busname}",
					"mask":       "send",
					"bus":        "session",
					"path":       "/org/mpris/MediaPlayer2",
					"interface":  "org.mpris.MediaPlayer2.Player",
					"member":     "PlayPause",
					"peer_label": "spotify",
				},
			},
		},
		{
			name:    "org.gnome.NautilusPreviewer",
			useFile: true,
			path:    filepath.Join(testdata, "systemd.log"),
			want: AppArmorLogs{
				{
					"apparmor":  "DENIED",
					"label":     "org.gnome.NautilusPreviewer",
					"operation": "dbus_bind",
					"name":      "org.gnome.NautilusPreviewer",
					"mask":      "bind",
					"bus":       "session",
					"info":      "Failed to register: GDBus.Error:org.freedesktop.DBus.Error.AccessDenied: " + ownNameDenied,
				},
			},
		},
		// Skipping live journalctl test as it depends on system state
		// {
		// 	name:    "journalctl",
		// 	useFile: false,
		// 	path:    "",
		// 	want:    AppArmorLogs{},
		// },
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
	// canReadPath := func(path string) bool {
	// 	if _, err := os.Stat(path); err == nil {
	// 		if file, err := os.Open(path); err == nil {
	// 			if err := file.Close(); err != nil {
	// 				return false
	// 			}
	// 			return true
	// 		}
	// 	}
	// 	return false
	// }

	t.Setenv("TMPDIR", "/tmp/tests")
	fifo := filepath.Join(t.TempDir(), "logs.fifo")
	if err := syscall.Mkfifo(fifo, 0o600); err != nil {
		t.Fatalf("mkfifo: %v", err)
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
		// {
		// 	name:    "Get /var/log/audit/audit.log.1",
		// 	path:    "1",
		// 	want:    "/var/log/audit/audit.log.1",
		// 	wantErr: !canReadPath("/var/log/audit/audit.log.1"),
		// },
		// {
		// 	name:    "Get default log file",
		// 	path:    "",
		// 	want:    "/var/log/audit/audit.log",
		// 	wantErr: !canReadPath("/var/log/audit/audit.log.1"),
		// },
		{
			name:    "Named pipe (process substitution)",
			path:    fifo,
			want:    fifo,
			wantErr: false,
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
