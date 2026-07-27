// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package builder

import (
	"testing"

	"github.com/roddhjav/apparmor.d/pkg/paths"
	"github.com/roddhjav/apparmor.d/pkg/prebuild"
	"github.com/roddhjav/apparmor.d/pkg/tasks"
)

// setupFlagDir points prebuild.FlagDir and tasks.Distribution at synthetic
// flag manifests, restored on test cleanup.
func setupFlagDir(t *testing.T) {
	t.Helper()
	flagDir := paths.New(t.TempDir())
	if err := flagDir.Join("main.conf").WriteFile([]byte(
		"foo complain\n" + // profile mode from main.conf
			"noflags\n" + // profile without any flag, must be skipped
			"baz not-a-mode\n", // invalid mode, Apply must error
	)); err != nil {
		t.Fatal(err)
	}
	if err := flagDir.Join("testdist.conf").WriteFile([]byte(
		"bar unconfined\n", // profile mode from the distribution flags
	)); err != nil {
		t.Fatal(err)
	}

	origFlagDir := prebuild.FlagDir
	origDist := tasks.Distribution
	prebuild.FlagDir = flagDir
	tasks.Distribution = "testdist"
	t.Cleanup(func() {
		prebuild.FlagDir = origFlagDir
		tasks.Distribution = origDist
	})
}

func TestProfileMode_Apply(t *testing.T) {
	setupFlagDir(t)
	b := NewProfileMode()
	b.SetConfig(cfg)

	tests := []struct {
		name    string
		profile string
		want    string
		wantErr bool
	}{
		{
			name:    "no-profile-name",
			profile: "# a file without profile header\n",
			want:    "# a file without profile header\n",
			wantErr: false,
		},
		{
			name:    "name-not-in-modes",
			profile: "profile other /usr/bin/other {\n}\n",
			want:    "profile other /usr/bin/other {\n}\n",
			wantErr: false,
		},
		{
			name:    "no-flags-entry",
			profile: "profile noflags /usr/bin/noflags {\n}\n",
			want:    "profile noflags /usr/bin/noflags {\n}\n",
			wantErr: false,
		},
		{
			name:    "set-mode-main",
			profile: "profile foo /usr/bin/foo {\n}\n",
			want:    "profile foo /usr/bin/foo flags=(complain) {\n}\n",
			wantErr: false,
		},
		{
			name:    "set-mode-dist",
			profile: "profile bar /usr/bin/bar flags=(attach_disconnected,complain) {\n}\n",
			want:    "profile bar /usr/bin/bar flags=(attach_disconnected,unconfined) {\n}\n",
			wantErr: false,
		},
		{
			name:    "invalid-mode",
			profile: "profile baz /usr/bin/baz {\n}\n",
			want:    "profile baz /usr/bin/baz {\n}\n",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := &Option{File: cfg.RootApparmor.Join(tt.name), Name: tt.name}
			got, err := b.Apply(opt, tt.profile)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProfileMode.Apply() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ProfileMode.Apply() = %v, want %v", got, tt.want)
			}
		})
	}
}
