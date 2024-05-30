// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2023-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package main

import (
	"os"
	"os/exec"
	"testing"

	"github.com/roddhjav/apparmor.d/pkg/prebuild/builder"
	"github.com/roddhjav/apparmor.d/pkg/prebuild/cfg"
	"github.com/roddhjav/apparmor.d/pkg/prebuild/prepare"
)

func chdirGitRoot() {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	root := string(out[0 : len(out)-1])
	if err := os.Chdir(root); err != nil {
		panic(err)
	}
}

func Test_AAPrebuild(t *testing.T) {
	tests := []struct {
		name     string
		wantErr  bool
		full     bool
		complain bool
		dist     string
	}{
		{
			name:     "Build for Archlinux",
			wantErr:  false,
			full:     false,
			complain: true,
			dist:     "arch",
		},
		{
			name:     "Build for Ubuntu",
			wantErr:  false,
			full:     true,
			complain: false,
			dist:     "ubuntu",
		},
		{
			name:     "Build for Debian",
			wantErr:  false,
			full:     true,
			complain: false,
			dist:     "debian",
		},
		{
			name:     "Build for OpenSUSE Tumbleweed",
			wantErr:  false,
			full:     true,
			complain: true,
			dist:     "opensuse",
		},
		// {
		// 	name:     "Build for Fedora",
		// 	wantErr:  true,
		// 	full:     false,
		// 	complain: false,
		// 	dist:     "fedora",
		// },
	}
	chdirGitRoot()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg.Distribution = tt.dist
			if tt.full {
				prepare.Register("fsp")
				builder.Register("fsp")
			}
			if tt.complain {
				builder.Register("complain")
			}
			if err := aaPrebuild(); (err != nil) != tt.wantErr {
				t.Errorf("aaPrebuild() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
