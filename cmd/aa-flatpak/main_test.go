// aa-flatpak - Confine flatpak applications
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package main

import (
	"os"
	"strings"
	"testing"

	"github.com/roddhjav/apparmor.d/pkg/aa"
	"github.com/roddhjav/apparmor.d/pkg/paths"
)

// testdata holds the metadata files of real flatpak applications.
var testdata = paths.New("../../tests/testdata/flatpak")

// setupTmp keeps the test temporary files under /tmp/tests.
func setupTmp(t *testing.T) {
	t.Helper()
	if err := os.MkdirAll("/tmp/tests", 0o755); err != nil {
		t.Fatalf("mkdir /tmp/tests: %v", err)
	}
	t.Setenv("TMPDIR", "/tmp/tests")
}

// setupTestEnv points the package directories at temporary ones. The apparmor
// magic root is emptied too: pkg/flatpak scans it to resolve dbus labels.
func setupTestEnv(t *testing.T) {
	t.Helper()
	setupTmp(t)
	origProfiles, origApp, origOverrides := profilesDir, appDir, overridesDir
	origMagic := aa.MagicRoot
	origEnforce := enforce
	t.Cleanup(func() {
		profilesDir, appDir, overridesDir = origProfiles, origApp, origOverrides
		aa.MagicRoot = origMagic
		enforce = origEnforce
	})

	profilesDir = paths.New(t.TempDir())
	appDir = paths.New(t.TempDir())
	overridesDir = t.TempDir()
	aa.MagicRoot = paths.New(t.TempDir())
}

// installApp copies a reference metadata file where an installed flatpak
// application would have it.
func installApp(t *testing.T, name string) {
	t.Helper()
	metadata := metadataPath(name)
	if err := metadata.Parent().MkdirAll(); err != nil {
		t.Fatalf("mkdir app dir: %v", err)
	}
	if err := testdata.Join(name).CopyTo(metadata); err != nil {
		t.Fatalf("copy metadata: %v", err)
	}
}

func TestMetadataPath(t *testing.T) {
	tests := []struct {
		name  string
		app   string
		want  string
		setup func(t *testing.T)
	}{
		{
			name: "chrome",
			app:  "com.google.Chrome",
			want: "com.google.Chrome/current/active/metadata",
		},
		{
			name: "steam",
			app:  "com.valvesoftware.Steam",
			want: "com.valvesoftware.Steam/current/active/metadata",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupTestEnv(t)
			want := appDir.Join(tt.want).String()
			if got := metadataPath(tt.app).String(); got != want {
				t.Errorf("metadataPath() = %v, want %v", got, want)
			}
		})
	}
}

func TestGenerateProfile(t *testing.T) {
	tests := []struct {
		name         string
		app          string
		install      bool
		enforce      bool
		override     string // content of the override file, "" means not created
		wantContains []string
		wantAbsent   []string
		wantErr      bool
	}{
		{
			name:    "chrome in complain mode",
			app:     "com.google.Chrome",
			install: true,
			wantContains: []string{
				"profile flatpak.com.google.Chrome",
				"complain",
				"include <abstractions/flatpak/baseapp/org.chromium.Chromium>",
				"include <abstractions/flatpak/sockets/x11>",
			},
		},
		{
			name:    "steam in enforce mode",
			app:     "com.valvesoftware.Steam",
			install: true,
			enforce: true,
			wantContains: []string{
				"profile flatpak.com.valvesoftware.Steam",
				"include <abstractions/flatpak/features/devel>",
			},
			wantAbsent: []string{"complain"},
		},
		{
			name:     "override is applied",
			app:      "com.google.Chrome",
			install:  true,
			override: "[Context]\nsockets=gpg-agent;\n",
			wantContains: []string{
				"include <abstractions/flatpak/sockets/gpg-agent>",
			},
		},
		{
			name:    "application not installed",
			app:     "com.google.Chrome",
			install: false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupTestEnv(t)
			enforce = tt.enforce
			if tt.install {
				installApp(t, tt.app)
			}
			if tt.override != "" {
				path := paths.New(overridesDir).Join(tt.app)
				if err := path.WriteFile([]byte(tt.override)); err != nil {
					t.Fatalf("write override: %v", err)
				}
			}

			out, err := generateProfile(tt.app)
			if (err != nil) != tt.wantErr {
				t.Fatalf("generateProfile() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			want := profilesDir.Join("flatpak." + tt.app).String()
			if out.String() != want {
				t.Errorf("generateProfile() = %v, want %v", out, want)
			}
			got, err := out.ReadFileAsString()
			if err != nil {
				t.Fatalf("read generated profile: %v", err)
			}
			for _, want := range tt.wantContains {
				if !strings.Contains(got, want) {
					t.Errorf("generateProfile() = missing %q in:\n%s", want, got)
				}
			}
			for _, absent := range tt.wantAbsent {
				if strings.Contains(got, absent) {
					t.Errorf("generateProfile() = unexpected %q in:\n%s", absent, got)
				}
			}
		})
	}
}

func TestRun(t *testing.T) {
	tests := []struct {
		name    string
		flags   func()
		wantErr bool
	}{
		{
			name:    "complain and enforce are mutually exclusive",
			flags:   func() { complain, enforce = true, true },
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(func() { complain, enforce = false, false })
			tt.flags()

			if err := run(); (err != nil) != tt.wantErr {
				t.Errorf("run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGenerateProfile_MissingAppDir(t *testing.T) {
	setupTestEnv(t)
	appDir = paths.New(t.TempDir()).Join("does-not-exist")

	if _, err := generateProfile("com.google.Chrome"); err == nil {
		t.Errorf("generateProfile() error = %v, wantErr %v", err, true)
	}
	if _, err := os.Stat(profilesDir.Join("flatpak.com.google.Chrome").String()); err == nil {
		t.Errorf("generateProfile() wrote a profile for a missing application")
	}
}
