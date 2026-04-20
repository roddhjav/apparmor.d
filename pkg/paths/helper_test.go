// This file is part of PathsHelper library.
// Copyright (C) 2018-2025 Arduino AG (http://www.arduino.cc/)
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package paths

import (
	"os"
	"testing"
)

func tempDir(t *testing.T) *Path {
	t.Helper()
	if err := os.MkdirAll("/tmp/tests", 0o755); err != nil {
		t.Fatalf("mkdir /tmp/tests: %v", err)
	}
	t.Setenv("TMPDIR", "/tmp/tests")
	return New(t.TempDir())
}

func TestFilter(t *testing.T) {
	tests := []struct {
		name string
		src  string
		want string
	}{
		{
			name: "comment",
			src:  "# comment",
			want: "",
		},
		{
			name: "comment with space",
			src:  " # comment",
			want: "",
		},
		{
			name: "no comment",
			src:  "no comment",
			want: "no comment",
		},
		{
			name: "no comment # comment",
			src:  "no comment # comment",
			want: "no comment",
		},
		{
			name: "empty",
			src: `

`,
			want: ``,
		},
		{
			name: "main",
			src: `
# Common profile flags definition for all distributions
# File format: one profile by line using the format: '<profile> <flags>'

bwrap attach_disconnected,mediate_deleted,complain
bwrap-app attach_disconnected,complain

akonadi_akonotes_resource complain # Dev
gnome-disks complain

`,
			want: `bwrap attach_disconnected,mediate_deleted,complain
bwrap-app attach_disconnected,complain
akonadi_akonotes_resource complain
gnome-disks complain
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLine := Filter(tt.src)
			if gotLine != tt.want {
				t.Errorf("Filter() got = |%v|, want |%v|", gotLine, tt.want)
			}
		})
	}
}

func TestPathListFromArgs(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(t *testing.T, dir *Path) []string
		wantLen int
		wantErr bool
	}{
		{
			name: "single file",
			setup: func(t *testing.T, dir *Path) []string {
				p := dir.Join("foo")
				if err := p.WriteFile([]byte("x")); err != nil {
					t.Fatalf("write: %v", err)
				}
				return []string{p.String()}
			},
			wantLen: 1,
		},
		{
			name: "directory filters README.md",
			setup: func(t *testing.T, dir *Path) []string {
				for _, name := range []string{"a", "sub/b", "README.md"} {
					p := dir.Join(name)
					if err := p.Parent().MkdirAll(); err != nil {
						t.Fatalf("mkdir: %v", err)
					}
					if err := p.WriteFile([]byte("x")); err != nil {
						t.Fatalf("write: %v", err)
					}
				}
				return []string{dir.String()}
			},
			wantLen: 2,
		},
		{
			name: "missing path errors",
			setup: func(t *testing.T, dir *Path) []string {
				return []string{dir.Join("missing").String()}
			},
			wantErr: true,
		},
		{
			name: "missing path resolves via magicRoot",
			setup: func(t *testing.T, dir *Path) []string {
				if err := dir.Join("named").WriteFile([]byte("x")); err != nil {
					t.Fatalf("write: %v", err)
				}
				return []string{"named"}
			},
			wantLen: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := tempDir(t)
			args := tt.setup(t, dir)
			got, err := PathListFromArgs(args, dir)
			if (err != nil) != tt.wantErr {
				t.Fatalf("PathListFromArgs() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if len(got) != tt.wantLen {
				t.Errorf("PathListFromArgs() returned %d entries, want %d: %v", len(got), tt.wantLen, got)
			}
			for _, p := range got {
				if p.Base() == "README.md" {
					t.Errorf("README.md was not filtered out: %v", got)
				}
			}
		})
	}
}
