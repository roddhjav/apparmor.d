// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2023-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package main

import (
	"os"
	"os/exec"
	"testing"

	"github.com/roddhjav/apparmor.d/pkg/prebuild"
	"github.com/roddhjav/apparmor.d/pkg/prebuild/builder"
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

func Test_main(t *testing.T) {
	tests := []struct {
		name string
		dist string
	}{
		{
			name: "Build for Archlinux",
			dist: "arch",
		},
		{
			name: "Build for Ubuntu",
			dist: "ubuntu",
		},
		{
			name: "Build for Debian",
			dist: "debian",
		},
		{
			name: "Build for OpenSUSE Tumbleweed",
			dist: "opensuse",
		},
	}
	chdirGitRoot()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prepare.Prepares = []prepare.Task{}
			builder.Builds = []builder.Builder{}
			prebuild.Distribution = tt.dist
			main()
		})
	}
}
