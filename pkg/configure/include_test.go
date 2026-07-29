// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package configure

import (
	"slices"
	"testing"

	"github.com/roddhjav/apparmor.d/pkg/paths"
	"github.com/roddhjav/apparmor.d/pkg/tasks"
)

func TestInclude_Active(t *testing.T) {
	tests := []struct {
		name    string
		entries map[string]string // include filename -> content
		want    bool
	}{
		{
			name:    "no user dir entries",
			entries: map[string]string{},
			want:    false,
		},
		{
			name: "one entry",
			entries: map[string]string{
				"a.conf": "foo\n",
			},
			want: true,
		},
		{
			name: "non-include file ignored",
			entries: map[string]string{
				"README.md": "foo\n",
			},
			want: false,
		},
	}

	setupTmp(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userDir := paths.New(t.TempDir())
			for name, content := range tt.entries {
				if err := userDir.Join(name).WriteFile([]byte(content)); err != nil {
					t.Fatalf("write include file: %v", err)
				}
			}
			if got := NewInclude(paths.PathList{userDir}).Active(); got != tt.want {
				t.Errorf("Active() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInclude_Apply(t *testing.T) {
	tests := []struct {
		name         string
		vendorFiles  map[string]string // vendor include filename -> content
		includeFiles map[string]string // user include filename -> content
		profiles     []string          // apparmor.d/<path> files to seed
		srcFiles     []string          // source <path> files to seed (restore mode)
		missingRoot  bool              // build root does not exist
		roDirs       []string          // apparmor.d relative dirs made read-only
		wantErr      bool
		wantKept     []string // apparmor.d-relative paths expected to remain
		wantRemoved  []string // apparmor.d-relative paths expected to be removed
		wantEmpty    bool
	}{
		{
			name:         "no entries is noop",
			includeFiles: map[string]string{},
			profiles:     []string{"foo", "bar"},
			wantKept:     []string{"foo", "bar"},
			wantEmpty:    true,
		},
		{
			name: "keep only listed profile",
			includeFiles: map[string]string{
				"user.conf": "foo\n",
			},
			profiles:    []string{"foo", "bar", "baz"},
			wantKept:    []string{"foo"},
			wantRemoved: []string{"bar", "baz"},
		},
		{
			name: "vendor and user entries combined",
			vendorFiles: map[string]string{
				"vendor.conf": "foo\n",
			},
			includeFiles: map[string]string{
				"user.conf": "bar\n",
			},
			profiles:    []string{"foo", "bar", "baz"},
			wantKept:    []string{"foo", "bar"},
			wantRemoved: []string{"baz"},
		},
		{
			name: "keep files under listed directory",
			includeFiles: map[string]string{
				"user.conf": "groups\n",
			},
			profiles: []string{
				"groups/a/x", "groups/a/y", "unrelated",
			},
			wantKept:    []string{"groups/a/x", "groups/a/y"},
			wantRemoved: []string{"unrelated"},
		},
		{
			name: "protect reserved directories",
			includeFiles: map[string]string{
				"user.conf": "foo\n",
			},
			profiles: []string{
				"foo", "bar",
				"abstractions/base",
				"tunables/global",
				"disable/something",
				"mappings/foo",
			},
			wantKept: []string{
				"foo",
				"abstractions/base",
				"tunables/global",
				"disable/something",
				"mappings/foo",
			},
			wantRemoved: []string{"bar"},
		},
		{
			name: "always kept profiles survive",
			includeFiles: map[string]string{
				"user.conf": "foo\n",
			},
			profiles:    []string{"foo", "bar", "child-open", "dbus-session", "namespaces/makepkg/gpg"},
			wantKept:    []string{"foo", "child-open", "dbus-session", "namespaces/makepkg/gpg"},
			wantRemoved: []string{"bar"},
		},
		{
			name: "restore ignored profile by name",
			includeFiles: map[string]string{
				"user.conf": "foo\n",
			},
			srcFiles: []string{"groups/a/foo", "bar"},
			profiles: []string{"bar"},
			wantKept: []string{"groups/a/foo", "bar"},
		},
		{
			name: "restore ignored directory",
			includeFiles: map[string]string{
				"user.conf": "groups/a\n",
			},
			srcFiles: []string{"groups/a/x", "groups/a/y", "bar"},
			profiles: []string{},
			wantKept: []string{"groups/a/x", "groups/a/y"},
		},
		{
			name: "restore keeps existing profiles untouched",
			includeFiles: map[string]string{
				"user.conf": "foo\n",
			},
			srcFiles: []string{"foo"},
			profiles: []string{"foo", "bar"},
			wantKept: []string{"foo", "bar"},
		},
		{
			name: "read dir error",
			includeFiles: map[string]string{
				"user.conf": "foo\n",
			},
			missingRoot: true,
			wantErr:     true,
		},
		{
			name: "remove profile error",
			includeFiles: map[string]string{
				"user.conf": "foo\n",
			},
			profiles: []string{"sub/bar"},
			roDirs:   []string{"sub"},
			wantErr:  true,
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
			for _, rel := range tt.profiles {
				f := c.RootApparmor.Join(rel)
				if err := f.Parent().MkdirAll(); err != nil {
					t.Fatalf("mkdir: %v", err)
				}
				if err := f.WriteFile([]byte("profile\n")); err != nil {
					t.Fatalf("write profile: %v", err)
				}
			}
			vendorDir := paths.New(t.TempDir())
			for name, content := range tt.vendorFiles {
				if err := vendorDir.Join(name).WriteFile([]byte(content)); err != nil {
					t.Fatalf("write vendor include file: %v", err)
				}
			}
			userDir := paths.New(t.TempDir())
			for name, content := range tt.includeFiles {
				if err := userDir.Join(name).WriteFile([]byte(content)); err != nil {
					t.Fatalf("write include file: %v", err)
				}
			}
			for _, dir := range tt.roDirs {
				chmodRO(t, c.RootApparmor.Join(dir))
			}

			task := NewInclude(paths.PathList{vendorDir, userDir})
			if tt.srcFiles != nil {
				src := paths.New(t.TempDir())
				for _, rel := range tt.srcFiles {
					f := src.Join(rel)
					if err := f.Parent().MkdirAll(); err != nil {
						t.Fatalf("mkdir: %v", err)
					}
					if err := f.WriteFile([]byte("profile\n")); err != nil {
						t.Fatalf("write source profile: %v", err)
					}
				}
				task = NewRestoreInclude(paths.PathList{vendorDir, userDir}, src)
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
			if tt.wantEmpty && len(got) != 0 {
				t.Errorf("Apply() = %v, want empty", got)
			}
			if !tt.wantEmpty && !slices.Contains(got, userDir.String()) {
				t.Errorf("Apply() = %v, want contains %q", got, userDir.String())
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
