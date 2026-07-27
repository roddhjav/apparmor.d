// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package configure

import (
	"testing"

	"github.com/roddhjav/apparmor.d/pkg/paths"
	"github.com/roddhjav/apparmor.d/pkg/tasks"
)

func TestMerge_Apply(t *testing.T) {
	tests := []struct {
		name        string
		dirs        []string          // apparmor.d relative directories to create
		files       map[string]string // apparmor.d relative file -> content
		roDirs      []string          // apparmor.d relative dirs made read-only
		noReadDirs  []string          // apparmor.d relative dirs made unreadable
		badRoot     bool              // build root name is an invalid glob pattern
		wantErr     bool
		wantFiles   []string // apparmor.d relative paths expected to exist
		wantNoFiles []string // apparmor.d relative paths expected to be gone
	}{
		{
			name:    "glob error",
			badRoot: true,
			wantErr: true,
		},
		{
			name: "merge groups profiles and namespaces",
			files: map[string]string{
				"groups/g1/prof1":       "profile prof1 {\n}\n",
				"profiles-a-f/prof2":    "profile prof2 {\n}\n",
				"namespaces/ns1/prof3":  "profile prof3 {\n}\n",
				"namespaces/ns2/prof4":  "profile prof4 {\n}\n",
				"groups/g2/local/prof5": "included\n",
			},
			wantFiles: []string{
				"prof1", "prof2", ":ns1:prof3", ":ns2:prof4", "local/prof5",
			},
			wantNoFiles: []string{"groups", "profiles-a-f", "namespaces"},
		},
		{
			name: "no namespaces directory",
			files: map[string]string{
				"groups/g1/prof1": "profile prof1 {\n}\n",
			},
			wantFiles:   []string{"prof1"},
			wantNoFiles: []string{"groups"},
		},
		{
			name: "rename conflict error",
			dirs: []string{"prof1"},
			files: map[string]string{
				"groups/g1/prof1": "profile prof1 {\n}\n",
			},
			wantErr: true,
		},
		{
			name: "remove groups error",
			files: map[string]string{
				"groups/readme": "not a group dir\n",
			},
			roDirs:  []string{"groups"},
			wantErr: true,
		},
		{
			name: "namespaces not a directory",
			files: map[string]string{
				"namespaces": "not a dir\n",
			},
			wantErr: true,
		},
		{
			name:       "namespace read error",
			dirs:       []string{"namespaces/ns1"},
			noReadDirs: []string{"namespaces/ns1"},
			wantErr:    true,
		},
		{
			name: "namespace rename conflict error",
			dirs: []string{":ns1:prof3"},
			files: map[string]string{
				"namespaces/ns1/prof3": "profile prof3 {\n}\n",
			},
			wantErr: true,
		},
		{
			name:    "remove namespaces error",
			dirs:    []string{"namespaces/ns1"},
			roDirs:  []string{"namespaces"},
			wantErr: true,
		},
	}

	setupTmp(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var c *tasks.TaskConfig
			if tt.badRoot {
				c = tasks.NewTaskConfig(paths.New(t.TempDir()).Join("bad["))
			} else {
				c = newTaskConfigTmp(t)
			}
			for _, dir := range tt.dirs {
				if err := c.RootApparmor.Join(dir).MkdirAll(); err != nil {
					t.Fatalf("mkdir: %v", err)
				}
			}
			seedFiles(t, c.RootApparmor, tt.files)
			for _, dir := range tt.roDirs {
				chmodRO(t, c.RootApparmor.Join(dir))
			}
			for _, dir := range tt.noReadDirs {
				chmodNoAccess(t, c.RootApparmor.Join(dir))
			}

			task := NewMerge()
			task.SetConfig(c)
			_, err := task.Apply()
			if (err != nil) != tt.wantErr {
				t.Errorf("Apply() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for _, rel := range tt.wantFiles {
				if c.RootApparmor.Join(rel).NotExist() {
					t.Errorf("Apply() = %v, want %v", rel, "exist")
				}
			}
			for _, rel := range tt.wantNoFiles {
				if c.RootApparmor.Join(rel).Exist() {
					t.Errorf("Apply() = %v, want %v", rel, "removed")
				}
			}
		})
	}
}
