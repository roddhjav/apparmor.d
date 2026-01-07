// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package configure

import (
	"os"
	"os/exec"
	"slices"
	"testing"

	"github.com/roddhjav/apparmor.d/pkg/paths"
	"github.com/roddhjav/apparmor.d/pkg/prebuild"
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

func TestTask_Apply(t *testing.T) {
	tests := []struct {
		name        string
		task        Task
		want        string
		wantErr     bool
		wantFiles   paths.PathList
		wantNoFiles paths.PathList
	}{
		{
			name:      "synchronise",
			task:      Tasks["synchronise"],
			wantErr:   false,
			wantFiles: paths.PathList{prebuild.RootApparmord.Join("/groups/_full/systemd")},
		},
		{
			name:    "ignore",
			task:    Tasks["ignore"],
			wantErr: false,
			want:    "dists/ignore/main.ignore",
		},
		{
			name:      "merge",
			task:      Tasks["merge"],
			wantErr:   false,
			wantFiles: paths.PathList{prebuild.RootApparmord.Join("aa-log")},
		},
		{
			name:    "configure",
			task:    Tasks["configure"],
			wantErr: false,
		},
		{
			name:    "setflags",
			task:    Tasks["setflags"],
			wantErr: false,
			want:    "dists/flags/main.flags",
		},
		{
			name:      "overwrite",
			task:      Tasks["overwrite"],
			wantErr:   false,
			wantFiles: paths.PathList{prebuild.RootApparmord.Join("flatpak.apparmor.d")},
		},
		{
			name:      "systemd-default",
			task:      Tasks["systemd-default"],
			wantErr:   false,
			wantFiles: paths.PathList{prebuild.Root.Join("systemd/system/dbus.service")},
		},
		{
			name:      "fsp",
			task:      Tasks["fsp"],
			wantErr:   false,
			wantFiles: paths.PathList{prebuild.RootApparmord.Join("systemd")},
		},
	}
	chdirGitRoot()
	_ = prebuild.Root.RemoveAll()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.task.Apply()
			if (err != nil) != tt.wantErr {
				t.Errorf("Task.Apply() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != "" && !slices.Contains(got, tt.want) {
				t.Errorf("Task.Apply() = %v, want %v", got, tt.want)
			}
			for _, file := range tt.wantFiles {
				if file.NotExist() {
					t.Errorf("Task.Apply() = %v, want %v", file, "exist")
				}
			}
		})
	}
}

func TestRegister(t *testing.T) {
	tests := []struct {
		name        string
		names       []string
		wantSuccess bool
	}{
		{
			name:        "test",
			names:       []string{"synchronise", "ignore"},
			wantSuccess: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Register(tt.names...)
			for _, name := range tt.names {
				if got := slices.Contains(Prepares, Tasks[name]); got != tt.wantSuccess {
					t.Errorf("Register() = %v, want %v", got, tt.wantSuccess)
				}

			}
		})
	}
}
