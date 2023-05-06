// aa-log - Review AppArmor generated messages
// Copyright (C) 2021 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package main

import (
	"testing"
)

func Test_app(t *testing.T) {
	tests := []struct {
		name    string
		logger  string
		path    string
		profile string
		anon    bool
		wantErr bool
	}{
		{
			name:    "Test audit.log",
			logger:  "auditd",
			path:    "../../tests/audit.log",
			profile: "",
			anon:    true,
			wantErr: false,
		},
		{
			name:    "Test Dbus Session",
			logger:  "systemd",
			path:    "../../tests/systemd.log",
			profile: "",
			anon:    false,
			wantErr: false,
		},
		{
			name:    "No logfile",
			logger:  "auditd",
			path:    "../../tests/log",
			profile: "",
			anon:    false,
			wantErr: true,
		},
		{
			name:    "Logger not supported",
			logger:  "raw",
			path:    "../../tests/audit.log",
			profile: "",
			anon:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := aaLog(tt.logger, tt.path, tt.profile, tt.anon); (err != nil) != tt.wantErr {
				t.Errorf("aaLog() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
