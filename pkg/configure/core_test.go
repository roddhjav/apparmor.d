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
	"github.com/roddhjav/apparmor.d/pkg/tasks"
)

var (
	cfg = tasks.NewTaskConfig(paths.New(".build"))
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
			task:      NewSynchronise([]*paths.Path{paths.New("apparmor.d"), paths.New("share")}),
			wantErr:   false,
			wantFiles: paths.PathList{cfg.RootApparmor.Join("/groups/_full/systemd")},
		},
		{
			name:    "ignore",
			task:    NewIgnore(),
			wantErr: false,
			want:    "dists/ignore/main.ignore",
		},
		{
			name:      "merge",
			task:      NewMerge(),
			wantErr:   false,
			wantFiles: paths.PathList{cfg.RootApparmor.Join("aa-log")},
		},
		{
			name:    "configure",
			task:    NewConfigure(),
			wantErr: false,
		},
		{
			name:    "setflags",
			task:    NewSetFlags(),
			wantErr: false,
			want:    "dists/flags/main.flags",
		},
		{
			name:      "overwrite",
			task:      NewOverwrite(false),
			wantErr:   false,
			wantFiles: paths.PathList{cfg.RootApparmor.Join("flatpak.apparmor.d")},
		},
		{
			name:      "systemd-default",
			task:      NewSystemdDefault(),
			wantErr:   false,
			wantFiles: paths.PathList{cfg.Root.Join("systemd/system/dbus.service")},
		},
		{
			name:      "fsp",
			task:      NewFullSystemPolicy(),
			wantErr:   false,
			wantFiles: paths.PathList{cfg.RootApparmor.Join("systemd")},
		},
	}
	chdirGitRoot()
	_ = cfg.Root.RemoveAll()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.task.SetConfig(cfg)
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

func TestConfigures_Add(t *testing.T) {
	tests := []struct {
		name  string
		tasks []Task
		want  []string
	}{
		{
			name:  "add-tasks",
			tasks: []Task{NewSynchronise(nil), NewIgnore()},
			want:  []string{"synchronise", "ignore"},
		},
	}
	c := tasks.NewTaskConfig(paths.New(".build"))
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRunner(c)
			for _, task := range tt.tasks {
				r.Add(task)
			}
			if len(r.Tasks) != len(tt.want) {
				t.Errorf("Configures.Add() len = %v, want %v", len(r.Tasks), len(tt.want))
			}
			for i, name := range tt.want {
				if r.Tasks[i].Name() != name {
					t.Errorf("Configures.Add() name = %v, want %v", r.Tasks[i].Name(), name)
				}
			}
		})
	}
}
