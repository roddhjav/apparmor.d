// aa-log - Review AppArmor generated messages
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package main

import (
	"path/filepath"
	"testing"
)

var (
	testdata = "../../tests/testdata/logs"
)

func Test_app(t *testing.T) {
	tests := []struct {
		name      string
		logger    string
		path      string
		profile   string
		namespace string
		rules     bool
		raw       bool
		load      bool
		wantErr   bool
	}{
		{
			name:    "Test audit.log",
			logger:  "auditd",
			path:    filepath.Join(testdata, "audit.log"),
			rules:   false,
			raw:     false,
			load:    false,
			wantErr: false,
		},
		{
			name:    "Test audit.log to rules",
			logger:  "auditd",
			path:    filepath.Join(testdata, "audit.log"),
			rules:   true,
			raw:     false,
			load:    false,
			wantErr: false,
		},
		{
			name:    "Test Dbus Session",
			logger:  "systemd",
			path:    filepath.Join(testdata, "systemd.log"),
			rules:   false,
			raw:     true,
			load:    false,
			wantErr: false,
		},
		{
			name:    "No logfile",
			logger:  "auditd",
			path:    filepath.Join(testdata, "log"),
			rules:   false,
			raw:     false,
			load:    false,
			wantErr: true,
		},
		{
			name:    "Logger not supported",
			logger:  "raw",
			path:    filepath.Join(testdata, "audit.log"),
			rules:   false,
			raw:     true,
			load:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rules = tt.rules
			raw = tt.raw
			load = tt.load
			if err := aaLog(tt.logger, tt.path, tt.profile, tt.namespace); (err != nil) != tt.wantErr {
				t.Errorf("aaLog() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
