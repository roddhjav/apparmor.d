// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2023-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package prebuild

import (
	"os"
	"os/exec"
	"testing"

	oss "github.com/roddhjav/apparmor.d/pkg/os"
)

func chdirGitRoot() {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	root := string(out)[0 : len(out)-1]
	if err := os.Chdir(root); err != nil {
		panic(err)
	}
}

func Test_PreBuild(t *testing.T) {
	tests := []struct {
		name     string
		wantErr  bool
		full     bool
		complain bool
		enforce  bool
		dist     string
	}{
		{
			name:     "Build for Archlinux",
			wantErr:  false,
			full:     false,
			complain: true,
			enforce:  false,
			dist:     "arch",
		},
		{
			name:     "Build for Ubuntu",
			wantErr:  false,
			full:     true,
			complain: false,
			enforce:  true,
			dist:     "ubuntu",
		},
		{
			name:     "Build for Debian",
			wantErr:  false,
			full:     true,
			complain: false,
			enforce:  false,
			dist:     "debian",
		},
		{
			name:     "Build for OpenSUSE Tumbleweed",
			wantErr:  false,
			full:     true,
			complain: true,
			enforce:  false,
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
			oss.Distribution = tt.dist
			if tt.full {
				Prepares = append(Prepares, SetFullSystemPolicy)
				Builds = append(Builds, BuildFullSystemPolicy)
			}
			if tt.complain {
				Builds = append(Builds, BuildComplain)
			}
			if tt.enforce {
				Builds = append(Builds, BuildEnforce)
			}
			if err := Prepare(); (err != nil) != tt.wantErr {
				t.Errorf("Prepare() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := Build(); (err != nil) != tt.wantErr {
				t.Errorf("Build() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
