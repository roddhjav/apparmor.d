// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package configure

import (
	"os"
	"slices"
	"testing"

	"github.com/roddhjav/apparmor.d/pkg/paths"
	"github.com/roddhjav/apparmor.d/pkg/tasks"
)

func setupTmp(t *testing.T) {
	t.Helper()
	if err := os.MkdirAll("/tmp/tests", 0o755); err != nil {
		t.Fatalf("mkdir /tmp/tests: %v", err)
	}
	t.Setenv("TMPDIR", "/tmp/tests")
}

func newTaskConfigTmp(t *testing.T) *tasks.TaskConfig {
	t.Helper()
	c := tasks.NewTaskConfig(paths.New(t.TempDir()))
	if err := c.RootApparmor.MkdirAll(); err != nil {
		t.Fatalf("mkdir apparmor.d: %v", err)
	}
	return c
}

// TestSetFlags_ApplyFileSources covers the prebuild use of the task: the
// sources are the dists/flags.d files of the target distribution, not
// directories, so other distributions' files must not be applied.
func TestSetFlags_ApplyFileSources(t *testing.T) {
	tests := []struct {
		name        string
		flagFiles   map[string]string // flag filename -> content
		sources     []string          // flag filenames passed to the task, in order
		want        string            // expected profile content after apply
		wantMissing string            // source filename expected absent from the result
	}{
		{
			name: "only given files applied",
			flagFiles: map[string]string{
				"main.conf":   "foo complain\n",
				"ubuntu.conf": "foo attach_disconnected\n",
			},
			sources: []string{"main.conf", "arch.conf"},
			want:    "profile foo /usr/bin/foo flags=(complain) {\n}\n",
		},
		{
			name: "later file overrides earlier",
			flagFiles: map[string]string{
				"main.conf": "foo complain\n",
				"arch.conf": "foo attach_disconnected\n",
			},
			sources: []string{"main.conf", "arch.conf"},
			want:    "profile foo /usr/bin/foo flags=(attach_disconnected) {\n}\n",
		},
	}

	setupTmp(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newTaskConfigTmp(t)
			profile := c.RootApparmor.Join("foo")
			if err := profile.WriteFile([]byte("profile foo /usr/bin/foo {\n}\n")); err != nil {
				t.Fatalf("write profile: %v", err)
			}
			flagDir := paths.New(t.TempDir())
			for name, content := range tt.flagFiles {
				if err := flagDir.Join(name).WriteFile([]byte(content)); err != nil {
					t.Fatalf("write flag file: %v", err)
				}
			}
			sources := paths.PathList{}
			for _, name := range tt.sources {
				sources = append(sources, flagDir.Join(name))
			}

			task := NewSetFlags(sources)
			task.SetConfig(c)
			if _, err := task.Apply(); err != nil {
				t.Fatalf("Apply() error = %v, wantErr false", err)
			}
			out, err := profile.ReadFileAsString()
			if err != nil {
				t.Fatalf("read profile: %v", err)
			}
			if out != tt.want {
				t.Errorf("Apply() profile = %q, want %q", out, tt.want)
			}
		})
	}
}

func TestSetFlags_Apply(t *testing.T) {
	tests := []struct {
		name            string
		vendorFiles     map[string]string // vendor flag filename -> content
		flagFiles       map[string]string // user flag filename -> content
		profiles        map[string]string // profile basename -> content
		noReadProfiles  []string          // profile basenames made unreadable
		noWriteProfiles []string          // profile basenames made read-only
		wantErr         bool
		wantResult      []string          // substrings expected in result
		wantContent     map[string]string // profile basename -> expected content after apply
		wantEmpty       bool              // true when Apply is expected to return no result
	}{
		{
			name:      "empty dir returns empty",
			flagFiles: map[string]string{},
			wantEmpty: true,
		},
		{
			name: "missing profile is reported not fatal",
			flagFiles: map[string]string{
				"user.conf": "missing-profile complain\n",
			},
			profiles:   map[string]string{},
			wantResult: []string{"Profile missing-profile not found, ignoring"},
		},
		{
			name: "overwritten profile matched by pkgname suffix",
			flagFiles: map[string]string{
				"user.conf": "loupe complain\n",
			},
			profiles: map[string]string{
				"loupe.apparmor.d": "profile loupe /usr/bin/loupe {\n}\n",
			},
			wantContent: map[string]string{
				"loupe.apparmor.d": "profile loupe /usr/bin/loupe flags=(complain) {\n}\n",
			},
		},
		{
			name: "apply single flag",
			flagFiles: map[string]string{
				"user.conf": "foo complain\n",
			},
			profiles: map[string]string{
				"foo": "profile foo /usr/bin/foo {\n  include <abstractions/base>\n}\n",
			},
			wantContent: map[string]string{
				"foo": "profile foo /usr/bin/foo flags=(complain) {\n  include <abstractions/base>\n}\n",
			},
		},
		{
			name: "apply multiple flags",
			flagFiles: map[string]string{
				"user.conf": "bar attach_disconnected,complain\n",
			},
			profiles: map[string]string{
				"bar": "profile bar /usr/bin/bar {\n}\n",
			},
			wantContent: map[string]string{
				"bar": "profile bar /usr/bin/bar flags=(attach_disconnected,complain) {\n}\n",
			},
		},
		{
			name: "multiple flag files merge",
			flagFiles: map[string]string{
				"a.conf": "foo complain\n",
				"b.conf": "bar enforce\n",
			},
			profiles: map[string]string{
				"foo": "profile foo /usr/bin/foo {\n}\n",
				"bar": "profile bar /usr/bin/bar {\n}\n",
			},
			wantContent: map[string]string{
				"foo": "profile foo /usr/bin/foo flags=(complain) {\n}\n",
				// enforce is the default mode: no flags clause is emitted.
				"bar": "profile bar /usr/bin/bar {\n}\n",
			},
		},
		{
			name: "mode entry preserves existing profile flags",
			flagFiles: map[string]string{
				"core.conf": "accounts-daemon enforce\n",
			},
			profiles: map[string]string{
				"accounts-daemon": "profile accounts-daemon @{exec_path} flags=(attach_disconnected) {\n}\n",
			},
			wantContent: map[string]string{
				"accounts-daemon": "profile accounts-daemon @{exec_path} flags=(attach_disconnected) {\n}\n",
			},
		},
		{
			name: "later file overrides earlier alphabetically",
			flagFiles: map[string]string{
				"10-base.conf":     "foo complain\n",
				"20-override.conf": "foo enforce\n",
			},
			profiles: map[string]string{
				"foo": "profile foo /usr/bin/foo {\n}\n",
			},
			// enforce wins (later file) and, as the default mode, clears the flags.
			wantContent: map[string]string{
				"foo": "profile foo /usr/bin/foo {\n}\n",
			},
		},
		{
			name: "same name user file replaces vendor file",
			vendorFiles: map[string]string{
				"user.conf": "foo complain\nbar complain\n",
			},
			flagFiles: map[string]string{
				"user.conf": "foo attach_disconnected\n",
			},
			profiles: map[string]string{
				"foo": "profile foo /usr/bin/foo {\n}\n",
				"bar": "profile bar /usr/bin/bar {\n}\n",
			},
			wantContent: map[string]string{
				"foo": "profile foo /usr/bin/foo flags=(attach_disconnected) {\n}\n",
				"bar": "profile bar /usr/bin/bar {\n}\n",
			},
		},
		{
			name: "different names merge per profile",
			vendorFiles: map[string]string{
				"00-main.conf": "foo complain\nbar complain\n",
			},
			flagFiles: map[string]string{
				"10-user.conf": "foo attach_disconnected\n",
			},
			profiles: map[string]string{
				"foo": "profile foo /usr/bin/foo {\n}\n",
				"bar": "profile bar /usr/bin/bar {\n}\n",
			},
			wantContent: map[string]string{
				"foo": "profile foo /usr/bin/foo flags=(attach_disconnected) {\n}\n",
				"bar": "profile bar /usr/bin/bar flags=(complain) {\n}\n",
			},
		},
		{
			name: "non-conf file ignored",
			flagFiles: map[string]string{
				"README.md": "not a flag file\n",
			},
			wantEmpty: true,
		},
		{
			name: "empty flag list skipped",
			flagFiles: map[string]string{
				"user.conf": "foo\n",
			},
			profiles: map[string]string{
				"foo": "profile foo /usr/bin/foo {\n}\n",
			},
			wantContent: map[string]string{
				"foo": "profile foo /usr/bin/foo {\n}\n",
			},
		},
		{
			name: "unreadable profile error",
			flagFiles: map[string]string{
				"user.conf": "foo complain\n",
			},
			profiles: map[string]string{
				"foo": "profile foo /usr/bin/foo {\n}\n",
			},
			noReadProfiles: []string{"foo"},
			wantErr:        true,
		},
		{
			name: "unwritable profile error",
			flagFiles: map[string]string{
				"user.conf": "foo complain\n",
			},
			profiles: map[string]string{
				"foo": "profile foo /usr/bin/foo {\n}\n",
			},
			noWriteProfiles: []string{"foo"},
			wantErr:         true,
		},
	}

	setupTmp(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newTaskConfigTmp(t)
			for name, content := range tt.profiles {
				if err := c.RootApparmor.Join(name).WriteFile([]byte(content)); err != nil {
					t.Fatalf("write profile: %v", err)
				}
			}

			vendorDir := paths.New(t.TempDir())
			for name, content := range tt.vendorFiles {
				if err := vendorDir.Join(name).WriteFile([]byte(content)); err != nil {
					t.Fatalf("write vendor flag file: %v", err)
				}
			}
			userDir := paths.New(t.TempDir())
			for name, content := range tt.flagFiles {
				if err := userDir.Join(name).WriteFile([]byte(content)); err != nil {
					t.Fatalf("write flag file: %v", err)
				}
			}
			for _, name := range tt.noReadProfiles {
				chmodNoAccess(t, c.RootApparmor.Join(name))
			}
			for _, name := range tt.noWriteProfiles {
				chmodRO(t, c.RootApparmor.Join(name))
			}

			task := NewSetFlags(paths.PathList{vendorDir, userDir})
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
			for _, w := range tt.wantResult {
				if !slices.Contains(got, w) {
					t.Errorf("Apply() = %v, want contains %q", got, w)
				}
			}
			for name, want := range tt.wantContent {
				out, err := c.RootApparmor.Join(name).ReadFileAsString()
				if err != nil {
					t.Errorf("read profile %s: %v", name, err)
					continue
				}
				if out != want {
					t.Errorf("profile %s = %q, want %q", name, out, want)
				}
			}
		})
	}
}
