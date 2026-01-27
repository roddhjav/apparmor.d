// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2023-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package runtime

import (
	"os"
	"os/exec"
	"testing"

	"github.com/roddhjav/apparmor.d/pkg/configure"
	"github.com/roddhjav/apparmor.d/pkg/directive"
	"github.com/roddhjav/apparmor.d/pkg/paths"
	"github.com/roddhjav/apparmor.d/pkg/tasks"
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

func TestRunners_Build(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
		abi     int
		dist    string
	}{
		{
			name:    "Build for Archlinux",
			wantErr: false,
			abi:     4,
			dist:    "arch",
		},
		{
			name:    "Build for Ubuntu",
			wantErr: false,
			abi:     4,
			dist:    "ubuntu",
		},
		{
			name:    "Build for Debian",
			wantErr: false,
			abi:     4,
			dist:    "debian",
		},
		{
			name:    "Build for OpenSUSE Tumbleweed",
			wantErr: false,
			abi:     4,
			dist:    "opensuse",
		},
	}
	chdirGitRoot()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tasks.Distribution = tt.dist
			root := paths.New("/tmp/tests").Join(tt.dist)
			cfg := tasks.NewTaskConfig(root)
			cfg.ABI = tt.abi
			cfg.Test = true
			r := NewRunners(cfg)

			// Add required configure tasks
			r.Configures.
				Add(configure.NewSynchronise([]*paths.Path{paths.New("apparmor.d")})).
				Add(configure.NewMerge())

			// Register all directives
			r.Directives.
				Register(directive.NewDbus()).
				Register(directive.NewExec()).
				Register(directive.NewFilterOnly()).
				Register(directive.NewFilterExclude()).
				Register(directive.NewStack())

			if err := r.Configure(); (err != nil) != tt.wantErr {
				t.Errorf("Runners.Configure() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := r.Build(); (err != nil) != tt.wantErr {
				t.Errorf("Runners.Build() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
