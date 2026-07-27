// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package configure

import (
	"os"
	"testing"

	"github.com/roddhjav/apparmor.d/pkg/paths"
	"github.com/roddhjav/apparmor.d/pkg/tasks"
)

// chmodTemp temporarily changes the permission of a path, restoring a
// sane default on test cleanup. Skips the test when running as root as
// permissions are then not enforced.
func chmodTemp(t *testing.T, p *paths.Path, mode os.FileMode) {
	t.Helper()
	if os.Geteuid() == 0 {
		t.Skip("running as root, permission based error paths cannot be tested")
	}
	if err := p.Chmod(mode); err != nil {
		t.Fatalf("chmod %s: %v", p, err)
	}
	t.Cleanup(func() { _ = p.Chmod(0o755) })
}

// chmodRO makes a path read-only so that write operations inside it fail.
func chmodRO(t *testing.T, p *paths.Path) {
	t.Helper()
	chmodTemp(t, p, 0o555)
}

// chmodNoAccess removes all permissions from a path.
func chmodNoAccess(t *testing.T, p *paths.Path) {
	t.Helper()
	chmodTemp(t, p, 0o000)
}

// seedFiles writes the given root-relative files, creating parent
// directories as needed.
func seedFiles(t *testing.T, root *paths.Path, files map[string]string) {
	t.Helper()
	for rel, content := range files {
		f := root.Join(rel)
		if err := f.Parent().MkdirAll(); err != nil {
			t.Fatalf("mkdir %s: %v", f.Parent(), err)
		}
		if err := f.WriteFile([]byte(content)); err != nil {
			t.Fatalf("write %s: %v", f, err)
		}
	}
}

func TestSynchronise_Apply(t *testing.T) {
	tests := []struct {
		name        string
		srcDirs     []string          // src root relative directories to create
		srcFiles    map[string]string // src root relative file -> content
		sources     []string          // src root relative paths given to the task
		rootFiles   map[string]string // build root relative file -> content (pre-existing)
		roRoot      bool              // make the build root read-only
		missingRoot bool              // build root under a read-only parent
		wantErr     bool
		wantContent map[string]string // build root relative file -> expected content
	}{
		{
			name:        "copy directory source",
			srcFiles:    map[string]string{"apparmor.d/groups/foo": "profile foo {\n}\n"},
			sources:     []string{"apparmor.d"},
			wantContent: map[string]string{"apparmor.d/groups/foo": "profile foo {\n}\n"},
		},
		{
			name: "copy multiple directory sources",
			srcFiles: map[string]string{
				"apparmor.d/groups/foo": "profile foo {\n}\n",
				"share/apparmor.d/bar":  "profile bar {\n}\n",
			},
			sources: []string{"apparmor.d", "share"},
			wantContent: map[string]string{
				"apparmor.d/groups/foo": "profile foo {\n}\n",
				"share/apparmor.d/bar":  "profile bar {\n}\n",
			},
		},
		{
			name:        "copy file source over existing destination",
			srcFiles:    map[string]string{"single.conf": "new\n"},
			sources:     []string{"single.conf"},
			rootFiles:   map[string]string{"single.conf": "old\n"},
			wantContent: map[string]string{"single.conf": "new\n"},
		},
		{
			name:      "remove destination error",
			srcFiles:  map[string]string{"single.conf": "new\n"},
			sources:   []string{"single.conf"},
			rootFiles: map[string]string{"single.conf": "old\n"},
			roRoot:    true,
			wantErr:   true,
		},
		{
			name:    "copy directory error",
			srcDirs: []string{"data"},
			sources: []string{"data"},
			roRoot:  true,
			wantErr: true,
		},
		{
			name:        "mkdir parent error",
			srcFiles:    map[string]string{"single.conf": "new\n"},
			sources:     []string{"single.conf"},
			missingRoot: true,
			wantErr:     true,
		},
		{
			name:     "copy file error",
			srcFiles: map[string]string{"single.conf": "new\n"},
			sources:  []string{"single.conf"},
			roRoot:   true,
			wantErr:  true,
		},
	}

	setupTmp(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var c *tasks.TaskConfig
			if tt.missingRoot {
				parent := paths.New(t.TempDir())
				chmodRO(t, parent)
				c = tasks.NewTaskConfig(parent.Join("build"))
			} else {
				c = newTaskConfigTmp(t)
			}
			seedFiles(t, c.Root, tt.rootFiles)

			srcRoot := paths.New(t.TempDir())
			for _, dir := range tt.srcDirs {
				if err := srcRoot.Join(dir).MkdirAll(); err != nil {
					t.Fatalf("mkdir: %v", err)
				}
			}
			seedFiles(t, srcRoot, tt.srcFiles)
			if tt.roRoot {
				chmodRO(t, c.Root)
			}

			sources := make([]*paths.Path, 0, len(tt.sources))
			for _, src := range tt.sources {
				sources = append(sources, srcRoot.Join(src))
			}
			task := NewSynchronise(sources)
			task.SetConfig(c)
			got, err := task.Apply()
			if (err != nil) != tt.wantErr {
				t.Errorf("Apply() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(got) != len(tt.sources) {
				t.Errorf("Apply() = %v, want %v entries", got, len(tt.sources))
			}
			for rel, want := range tt.wantContent {
				out, err := c.Root.Join(rel).ReadFileAsString()
				if err != nil {
					t.Errorf("read %s: %v", rel, err)
					continue
				}
				if out != want {
					t.Errorf("Apply() %s = %q, want %q", rel, out, want)
				}
			}
		})
	}
}
