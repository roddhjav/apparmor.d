// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package main

import (
	"testing"

	"github.com/roddhjav/apparmor.d/pkg/paths"
)

// setVendorConfigDir relocates vendorConfigDir to a temporary directory,
// restoring the original on cleanup.
func setVendorConfigDir(t *testing.T) *paths.Path {
	t.Helper()
	original := vendorConfigDir
	vendorConfigDir = paths.New(t.TempDir())
	t.Cleanup(func() { vendorConfigDir = original })
	return vendorConfigDir
}

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name         string
		vendorModes  string // content of the vendor modes file, "" means not created
		modes        string // content of the config modes file, "" means not created
		flags        func() // set the -c/-e option globals
		want         string
		wantInclude  string // expected include mode, "" means default
		wantNoReload bool   // reload is expected unless set
		wantErr      bool
	}{
		{
			name:  "default mode",
			flags: func() {},
			want:  "complain",
		},
		{
			name:  "mode from config",
			modes: "default enforce\n",
			flags: func() {},
			want:  "enforce",
		},
		{
			name:        "mode from vendor config",
			vendorModes: "default enforce\n",
			flags:       func() {},
			want:        "enforce",
		},
		{
			name:        "config overrides vendor",
			vendorModes: "default enforce\n",
			modes:       "default complain\n",
			flags:       func() {},
			want:        "complain",
		},
		{
			name:  "complain overrides config",
			modes: "default enforce\n",
			flags: func() { complain = true },
			want:  "complain",
		},
		{
			name:  "enforce overrides config",
			modes: "default complain\n",
			flags: func() { enforce = true },
			want:  "enforce",
		},
		{
			name:    "invalid mode",
			modes:   "default unconfined\n",
			flags:   func() {},
			wantErr: true,
		},
		{
			name:        "include mode from config",
			modes:       "include full\n",
			flags:       func() {},
			want:        "complain",
			wantInclude: "full",
		},
		{
			name:        "include all from config",
			modes:       "include all\n",
			flags:       func() {},
			want:        "complain",
			wantInclude: "all",
		},
		{
			name:        "all overrides include config",
			modes:       "include full\n",
			flags:       func() { all = true },
			want:        "complain",
			wantInclude: "all",
		},
		{
			name:    "invalid include mode",
			modes:   "include bogus\n",
			flags:   func() {},
			wantErr: true,
		},
		{
			name:         "reload disabled from config",
			modes:        "reload no\n",
			flags:        func() {},
			want:         "complain",
			wantNoReload: true,
		},
		{
			name:         "no-reload overrides config",
			modes:        "reload yes\n",
			flags:        func() { noReload = true },
			want:         "complain",
			wantNoReload: true,
		},
		{
			name:    "invalid reload mode",
			modes:   "reload maybe\n",
			flags:   func() {},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vendorDir := setVendorConfigDir(t)
			configDir := paths.New(t.TempDir())
			if tt.vendorModes != "" {
				writeFile(t, vendorDir.Join("modes"), tt.vendorModes)
			}
			if tt.modes != "" {
				writeFile(t, configDir.Join("modes"), tt.modes)
			}
			t.Cleanup(func() { complain, enforce, noReload, all = false, false, false, false })
			tt.flags()

			cfg, err := loadConfig(configDir)
			if (err != nil) != tt.wantErr {
				t.Fatalf("loadConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if cfg.mode != tt.want {
				t.Errorf("loadConfig() mode = %q, want %q", cfg.mode, tt.want)
			}
			wantInclude := tt.wantInclude
			if wantInclude == "" {
				wantInclude = "default"
			}
			if cfg.include != wantInclude {
				t.Errorf("loadConfig() include = %q, want %q", cfg.include, wantInclude)
			}
			if cfg.reload == tt.wantNoReload {
				t.Errorf("loadConfig() reload = %v, want %v", cfg.reload, !tt.wantNoReload)
			}
			wantFlagDirs := paths.PathList{vendorDir.Join("flags.d"), configDir.Join("flags.d")}
			for i, dir := range wantFlagDirs {
				if cfg.flagDirs[i].String() != dir.String() {
					t.Errorf("loadConfig() flagDirs[%d] = %s, want %s", i, cfg.flagDirs[i], dir)
				}
			}
		})
	}
}

func TestReadModeConfig(t *testing.T) {
	tests := []struct {
		name     string
		content  string // "" means the file is not created
		override string // content of a second file overriding the first
		key      string
		want     string
	}{
		{
			name:    "default mode",
			content: "default enforce\n",
			key:     "default",
			want:    "enforce",
		},
		{
			name:    "comments and blanks ignored",
			content: "# deploy mode\n\ndefault complain\n",
			key:     "default",
			want:    "complain",
		},
		{
			name:    "missing file is empty",
			content: "",
			key:     "default",
			want:    "",
		},
		{
			name:    "unknown key absent",
			content: "default enforce\n",
			key:     "other",
			want:    "",
		},
		{
			name:     "later file overrides earlier",
			content:  "default enforce\n",
			override: "default complain\n",
			key:      "default",
			want:     "complain",
		},
		{
			name:     "later file fully replaces earlier",
			content:  "default enforce\n",
			override: "other value\n",
			key:      "default",
			want:     "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tempPath(t, "modes")
			if tt.content != "" {
				if err := path.WriteFile([]byte(tt.content)); err != nil {
					t.Fatalf("write modes: %v", err)
				}
			}
			override := tempPath(t, "modes")
			if tt.override != "" {
				if err := override.WriteFile([]byte(tt.override)); err != nil {
					t.Fatalf("write override modes: %v", err)
				}
			}
			if got := readModeConfig(path, override)[tt.key]; got != tt.want {
				t.Errorf("readModeConfig()[%q] = %q, want %q", tt.key, got, tt.want)
			}
		})
	}
}

func TestUserModeOverrides(t *testing.T) {
	tests := []struct {
		name        string
		vendorFiles map[string]string
		flagFiles   map[string]string
		want        []string // profiles expected in the override set
		wantNot     []string // profiles expected absent
	}{
		{
			name:      "mode entry overrides",
			flagFiles: map[string]string{"core.conf": "ps enforce\ngjs complain\n"},
			want:      []string{"ps", "gjs"},
		},
		{
			name:      "non-mode-only entry does not override",
			flagFiles: map[string]string{"core.conf": "foo attach_disconnected\n"},
			wantNot:   []string{"foo"},
		},
		{
			name:      "no flags dir",
			flagFiles: nil,
			wantNot:   []string{"ps"},
		},
		{
			name:        "vendor mode entry overrides",
			vendorFiles: map[string]string{"00-core.conf": "ps complain\n"},
			flagFiles:   map[string]string{"10-user.conf": "gjs attach_disconnected\n"},
			want:        []string{"ps"},
			wantNot:     []string{"gjs"},
		},
		{
			name:        "same name user file replaces vendor file",
			vendorFiles: map[string]string{"core.conf": "ps complain\n"},
			flagFiles:   map[string]string{"core.conf": "gjs attach_disconnected\n"},
			wantNot:     []string{"ps", "gjs"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writeDir := func(files map[string]string) *paths.Path {
				dir := paths.New(t.TempDir())
				for name, content := range files {
					if err := dir.Join(name).WriteFile([]byte(content)); err != nil {
						t.Fatalf("write flag file: %v", err)
					}
				}
				return dir
			}
			dirs := paths.PathList{writeDir(tt.vendorFiles), writeDir(tt.flagFiles)}

			got := userModeOverrides(dirs)
			for _, p := range tt.want {
				if !got[p] {
					t.Errorf("userModeOverrides() missing %q, got %v", p, got)
				}
			}
			for _, p := range tt.wantNot {
				if got[p] {
					t.Errorf("userModeOverrides() contains %q, want absent", p)
				}
			}
		})
	}
}
