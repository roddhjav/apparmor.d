// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2023-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package runtime

import (
	"errors"
	"os"
	"os/exec"
	"testing"

	"github.com/roddhjav/apparmor.d/pkg/builder"
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

// fakeConfigureTask is a minimal configure.Task used to drive the
// Runners.Configure logging and error branches.
type fakeConfigureTask struct {
	tasks.BaseTask
	msg []string
	err error
}

func (f *fakeConfigureTask) Apply() ([]string, error) {
	return f.msg, f.err
}

// newTestRunners returns Runners rooted in a temporary directory with an
// empty apparmor.d tree.
func newTestRunners(t *testing.T) *Runners {
	t.Helper()
	cfg := tasks.NewTaskConfig(paths.New(t.TempDir()))
	if err := cfg.RootApparmor.MkdirAll(); err != nil {
		t.Fatal(err)
	}
	return NewRunners(cfg)
}

// writeProfile writes a profile file in the runner apparmor.d tree.
func writeProfile(t *testing.T, r *Runners, name string, content string) *paths.Path {
	t.Helper()
	file := r.RootApparmor.Join(name)
	if err := file.WriteFile([]byte(content)); err != nil {
		t.Fatal(err)
	}
	return file
}

// chmodProfile changes the profile file mode, restored on test cleanup.
// Skips the test when run as root as file permissions are not enforced.
func chmodProfile(t *testing.T, file *paths.Path, mode os.FileMode) {
	t.Helper()
	if os.Geteuid() == 0 {
		t.Skip("file permissions are not enforced for root")
	}
	if err := file.Chmod(mode); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = file.Chmod(0o644) })
}

func TestRunners_Configure(t *testing.T) {
	tests := []struct {
		name    string
		task    func(r *Runners) configure.Task
		wantErr bool
	}{
		{
			name: "log messages and warnings",
			task: func(r *Runners) configure.Task {
				return &fakeConfigureTask{
					BaseTask: tasks.BaseTask{Keyword: "fake", Msg: "fake configure task"},
					msg:      []string{"profile foo not found", "regular message"},
				}
			},
			wantErr: false,
		},
		{
			name: "task error",
			task: func(r *Runners) configure.Task {
				return &fakeConfigureTask{
					BaseTask: tasks.BaseTask{Keyword: "fake", Msg: "failing configure task"},
					err:      errors.New("apply failed"),
				}
			},
			wantErr: true,
		},
		{
			name: "synchronise error",
			task: func(r *Runners) configure.Task {
				return configure.NewSynchronise([]*paths.Path{r.Root.Join("does-not-exist")})
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newTestRunners(t)
			r.Configures.Add(tt.task(r))
			if err := r.Configure(); (err != nil) != tt.wantErr {
				t.Errorf("Runners.Configure() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRunners_BuildBranches(t *testing.T) {
	tests := []struct {
		name    string
		profile string
		setup   func(t *testing.T, r *Runners, file *paths.Path)
		want    string // expected profile content after Build, empty to skip the check
		wantErr bool
	}{
		{
			name:    "builder error",
			profile: "profile foo /usr/bin/foo {\n}\n",
			setup: func(t *testing.T, r *Runners, file *paths.Path) {
				r.Builders.Add(builder.NewUserspace())
			},
			wantErr: true,
		},
		{
			name:    "directive error",
			profile: "profile foo /usr/bin/foo {\n  #aa:unknown\n}\n",
			wantErr: true,
		},
		{
			name:    "unreadable file",
			profile: "profile foo /usr/bin/foo {\n}\n",
			setup: func(t *testing.T, r *Runners, file *paths.Path) {
				chmodProfile(t, file, 0o000)
			},
			wantErr: true,
		},
		{
			name:    "unwritable file",
			profile: "profile foo /usr/bin/foo {\n}\n",
			setup: func(t *testing.T, r *Runners, file *paths.Path) {
				chmodProfile(t, file, 0o444)
			},
			wantErr: true,
		},
		{
			name:    "skip missing file and list directives",
			profile: "profile foo /usr/bin/foo {\n}\n",
			setup: func(t *testing.T, r *Runners, file *paths.Path) {
				// A dangling symlink is listed but must be skipped by Build.
				dangling := r.RootApparmor.Join("dangling")
				if err := os.Symlink(r.RootApparmor.Join("missing").String(), dangling.String()); err != nil {
					t.Fatal(err)
				}
				r.Builders.Add(builder.NewComplain())
				r.Directives.Register(directive.NewDbus())
			},
			want:    "profile foo /usr/bin/foo flags=(complain) {\n}\n",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newTestRunners(t)
			file := writeProfile(t, r, "foo", tt.profile)
			if tt.setup != nil {
				tt.setup(t, r, file)
			}
			if err := r.Build(); (err != nil) != tt.wantErr {
				t.Errorf("Runners.Build() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want == "" {
				return
			}
			got, err := file.ReadFileAsString()
			if err != nil {
				t.Fatal(err)
			}
			if got != tt.want {
				t.Errorf("Runners.Build() profile = %v, want %v", got, tt.want)
			}
		})
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
