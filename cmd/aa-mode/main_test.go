// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package main

import (
	"os"
	"testing"

	"github.com/roddhjav/apparmor.d/pkg/paths"
)

func tempDir(t *testing.T) *paths.Path {
	t.Helper()
	if err := os.MkdirAll("/tmp/tests", 0o755); err != nil {
		t.Fatalf("mkdir /tmp/tests: %v", err)
	}
	t.Setenv("TMPDIR", "/tmp/tests")
	return paths.New(t.TempDir())
}

func resetFlags(t *testing.T) {
	t.Helper()
	enforce, complain, kill = false, false, false
	defaultAllow, unconfined, prompt = false, false, false
	noReload = true
}

func TestSelectedMode(t *testing.T) {
	tests := []struct {
		name    string
		setup   func()
		want    string
		wantErr bool
	}{
		{
			name:  "enforce",
			setup: func() { enforce = true },
			want:  "enforce",
		},
		{
			name:  "complain",
			setup: func() { complain = true },
			want:  "complain",
		},
		{
			name:  "kill",
			setup: func() { kill = true },
			want:  "kill",
		},
		{
			name:  "default_allow",
			setup: func() { defaultAllow = true },
			want:  "default_allow",
		},
		{
			name:  "unconfined",
			setup: func() { unconfined = true },
			want:  "unconfined",
		},
		{
			name:  "prompt",
			setup: func() { prompt = true },
			want:  "prompt",
		},
		{
			name:    "no mode set",
			setup:   func() {},
			wantErr: true,
		},
		{
			name:    "two modes set",
			setup:   func() { enforce = true; complain = true },
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetFlags(t)
			tt.setup()
			got, err := selectedMode()
			if (err != nil) != tt.wantErr {
				t.Fatalf("selectedMode() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && got != tt.want {
				t.Errorf("selectedMode() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestAaSetMode(t *testing.T) {
	tests := []struct {
		name    string
		profile string
		mode    string
		want    string
		wantErr bool
	}{
		{
			name:    "add complain to unflagged profile",
			profile: "profile foo /usr/bin/foo {\n}\n",
			mode:    "complain",
			want:    "profile foo /usr/bin/foo flags=(complain) {\n}\n",
		},
		{
			name:    "replace complain with kill",
			profile: "profile foo /usr/bin/foo flags=(complain) {\n}\n",
			mode:    "kill",
			want:    "profile foo /usr/bin/foo flags=(kill) {\n}\n",
		},
		{
			name:    "enforce removes mode flag",
			profile: "profile foo /usr/bin/foo flags=(complain) {\n}\n",
			mode:    "enforce",
			want:    "profile foo /usr/bin/foo {\n}\n",
		},
		{
			name:    "unknown mode errors",
			profile: "profile foo /usr/bin/foo {\n}\n",
			mode:    "bogus",
			wantErr: true,
		},
		{
			name:    "unconfined profile is never modified",
			profile: "profile foo /usr/bin/foo flags=(unconfined) {\n}\n",
			mode:    "complain",
			want:    "profile foo /usr/bin/foo flags=(unconfined) {\n}\n",
		},
		{
			name:    "unconfined profile preserved even when setting unconfined",
			profile: "profile foo /usr/bin/foo flags=(attach_disconnected, unconfined) {\n}\n",
			mode:    "unconfined",
			want:    "profile foo /usr/bin/foo flags=(attach_disconnected, unconfined) {\n}\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetFlags(t)
			path := tempDir(t).Join("foo")
			if err := path.WriteFile([]byte(tt.profile)); err != nil {
				t.Fatalf("write profile: %v", err)
			}
			err := aaSetMode(paths.PathList{path}, tt.mode)
			if (err != nil) != tt.wantErr {
				t.Fatalf("aaSetMode() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			got, err := path.ReadFileAsString()
			if err != nil {
				t.Fatalf("read profile: %v", err)
			}
			if got != tt.want {
				t.Errorf("profile content = %q, want %q", got, tt.want)
			}
		})
	}
}

