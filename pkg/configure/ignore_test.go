// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package configure

import (
	"slices"
	"testing"

	"github.com/roddhjav/apparmor.d/pkg/paths"
	"github.com/roddhjav/apparmor.d/pkg/prebuild"
	"github.com/roddhjav/apparmor.d/pkg/tasks"
)

// setIgnoreDir points prebuild.IgnoreDir to a temporary directory seeded
// with main.conf, restoring the original on cleanup.
func setIgnoreDir(t *testing.T, mainContent string) {
	t.Helper()
	dir := paths.New(t.TempDir())
	if mainContent != "" {
		if err := dir.Join("main.conf").WriteFile([]byte(mainContent)); err != nil {
			t.Fatalf("write main.conf: %v", err)
		}
	}
	original := prebuild.IgnoreDir
	prebuild.IgnoreDir = dir
	t.Cleanup(func() { prebuild.IgnoreDir = original })
}

// userIgnoreDir returns a temporary directory seeded with a user.conf
// drop-in holding content.
func userIgnoreDir(t *testing.T, name string, content string) *paths.Path {
	t.Helper()
	dir := paths.New(t.TempDir())
	if content != "" {
		if err := dir.Join(name).WriteFile([]byte(content)); err != nil {
			t.Fatalf("write %s: %v", name, err)
		}
	}
	return dir
}

func TestIgnore_Apply(t *testing.T) {
	tests := []struct {
		name           string
		userOnly       bool              // apply NewUserIgnore instead of NewIgnore
		profiles       map[string]string // apparmor.d relative file -> content
		mainIgnore     string            // content of dists ignore.d main.conf
		vendorIgnore   string            // content of vendor drop-in 00-main.conf
		userIgnore     string            // content of the user drop-in file
		userIgnoreName string            // basename of the user drop-in, default 10-user.conf
		missingRoot    bool              // build root does not exist
		roDirs         []string          // apparmor.d relative dirs made read-only
		wantErr        bool
		wantEmpty      bool     // Apply is expected to return no result
		wantKept       []string // apparmor.d relative paths expected to remain
		wantRemoved    []string // apparmor.d relative paths expected to be removed
	}{
		{
			name:     "no user entries is noop",
			userOnly: true,
			profiles: map[string]string{
				"foo": "profile foo {\n}\n",
			},
			wantKept:  []string{"foo"},
			wantEmpty: true,
		},
		{
			name: "remove root entries and recursive names",
			profiles: map[string]string{
				"foo":           "profile foo {\n}\n",
				"groups/g1/bar": "profile bar {\n}\n",
				"groups/g1/baz": "profile baz {\n}\n",
			},
			mainIgnore:  "apparmor.d/foo\nbar\n",
			wantKept:    []string{"groups/g1/baz"},
			wantRemoved: []string{"foo", "groups/g1/bar"},
		},
		{
			name:     "user ignore entries removed",
			userOnly: true,
			profiles: map[string]string{
				"foo": "profile foo {\n}\n",
				"bar": "profile bar {\n}\n",
			},
			userIgnore:  "apparmor.d/foo\n",
			wantKept:    []string{"bar"},
			wantRemoved: []string{"foo"},
		},
		{
			name:     "vendor and user ignore entries combined",
			userOnly: true,
			profiles: map[string]string{
				"foo": "profile foo {\n}\n",
				"bar": "profile bar {\n}\n",
				"baz": "profile baz {\n}\n",
			},
			vendorIgnore: "apparmor.d/foo\n",
			userIgnore:   "apparmor.d/bar\n",
			wantKept:     []string{"baz"},
			wantRemoved:  []string{"foo", "bar"},
		},
		{
			name:     "same name user file replaces vendor file",
			userOnly: true,
			profiles: map[string]string{
				"foo": "profile foo {\n}\n",
				"bar": "profile bar {\n}\n",
			},
			vendorIgnore:   "apparmor.d/foo\n",
			userIgnore:     "apparmor.d/bar\n",
			userIgnoreName: "00-main.conf",
			wantKept:       []string{"foo"},
			wantRemoved:    []string{"bar"},
		},
		{
			name:        "recursive read error",
			mainIgnore:  "does-not-exist\n",
			missingRoot: true,
			wantErr:     true,
		},
		{
			name:        "user entries error",
			userOnly:    true,
			userIgnore:  "does-not-exist\n",
			missingRoot: true,
			wantErr:     true,
		},
		{
			name: "remove root entry error",
			profiles: map[string]string{
				"foo": "profile foo {\n}\n",
			},
			mainIgnore: "apparmor.d/foo\n",
			roDirs:     []string{"."},
			wantErr:    true,
		},
		{
			name: "remove recursive entry error",
			profiles: map[string]string{
				"groups/bar": "profile bar {\n}\n",
			},
			mainIgnore: "bar\n",
			roDirs:     []string{"groups"},
			wantErr:    true,
		},
	}

	setupTmp(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var c *tasks.TaskConfig
			if tt.missingRoot {
				c = tasks.NewTaskConfig(paths.New(t.TempDir()).Join("missing"))
			} else {
				c = newTaskConfigTmp(t)
			}
			seedFiles(t, c.RootApparmor, tt.profiles)
			setIgnoreDir(t, tt.mainIgnore)
			userName := tt.userIgnoreName
			if userName == "" {
				userName = "10-user.conf"
			}
			vendorDir := userIgnoreDir(t, "00-main.conf", tt.vendorIgnore)
			userDir := userIgnoreDir(t, userName, tt.userIgnore)
			for _, dir := range tt.roDirs {
				if dir == "." {
					chmodRO(t, c.RootApparmor)
				} else {
					chmodRO(t, c.RootApparmor.Join(dir))
				}
			}

			var task *Ignore
			if tt.userOnly {
				task = NewUserIgnore(paths.PathList{vendorDir, userDir})
			} else {
				task = NewIgnore()
			}
			task.SetConfig(c)
			got, err := task.Apply()
			if (err != nil) != tt.wantErr {
				t.Errorf("Apply() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if tt.wantEmpty {
				if len(got) != 0 {
					t.Errorf("Apply() = %v, want empty", got)
				}
			} else {
				want := prebuild.IgnoreDir.Join("main.conf").String()
				if tt.userOnly {
					want = userDir.String()
				}
				if !slices.Contains(got, want) {
					t.Errorf("Apply() = %v, want %v", got, want)
				}
			}
			for _, rel := range tt.wantKept {
				if c.RootApparmor.Join(rel).NotExist() {
					t.Errorf("expected %s kept, was removed", rel)
				}
			}
			for _, rel := range tt.wantRemoved {
				if c.RootApparmor.Join(rel).Exist() {
					t.Errorf("expected %s removed, still exists", rel)
				}
			}
		})
	}
}
