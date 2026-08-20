// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package configure

import (
	"testing"

	"github.com/roddhjav/apparmor.d/pkg/aa"
	"github.com/roddhjav/apparmor.d/pkg/paths"
)

func TestOverwrite_Apply(t *testing.T) {
	tests := []struct {
		name        string
		files       map[string]string // root relative file -> content
		entries     []string          // overwrite.d config entries
		upstream    []string          // upstream profiles present on the install target
		roRoot      bool
		want        []string
		wantErr     bool
		wantFiles   []string
		wantLinks   []string // present as symlinks
		wantNoFiles []string
	}{
		{
			name:        "upstream present renames profile and creates link",
			files:       map[string]string{"hostname": "profile hostname {}\n"},
			entries:     []string{"hostname"},
			upstream:    []string{"hostname"},
			want:        []string{"hostname"},
			wantFiles:   []string{"hostname.apparmor.d"},
			wantLinks:   []string{"disable/hostname"},
			wantNoFiles: []string{"hostname"},
		},
		{
			name:        "upstream absent still renames profile but skips link",
			files:       map[string]string{"hostname": "profile hostname {}\n"},
			entries:     []string{"hostname"},
			want:        []string{"hostname"},
			wantFiles:   []string{"hostname.apparmor.d"},
			wantNoFiles: []string{"hostname", "disable/hostname"},
		},
		{
			name:      "upstream without our profile only gets the link",
			files:     map[string]string{"other": "profile other {}\n"},
			entries:   []string{"usr.bin.foo"},
			upstream:  []string{"usr.bin.foo"},
			want:      []string{"usr.bin.foo"},
			wantFiles: []string{"other"},
			wantLinks: []string{"disable/usr.bin.foo"},
		},
		{
			name:        "no entries is noop",
			files:       map[string]string{"hostname": "profile hostname {}\n"},
			wantFiles:   []string{"hostname"},
			wantNoFiles: []string{"disable"},
		},
		{
			name:     "rename error",
			files:    map[string]string{"hostname": "profile hostname {}\n"},
			entries:  []string{"hostname"},
			upstream: []string{"hostname"},
			roRoot:   true,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newTaskConfigTmp(t)
			seedFiles(t, c.RootApparmor, tt.files)
			target := paths.New(t.TempDir())
			for _, name := range tt.upstream {
				if err := target.Join(name).WriteFile([]byte("profile x {\n}\n")); err != nil {
					t.Fatalf("seed target: %v", err)
				}
			}
			oldRoot := aa.MagicRoot
			aa.MagicRoot = target
			t.Cleanup(func() { aa.MagicRoot = oldRoot })
			if tt.roRoot {
				chmodRO(t, c.RootApparmor)
			}

			task := NewOverwrite(tt.entries)
			task.SetConfig(c)
			got, err := task.Apply()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Apply() error = %v, wantErr %v", err, tt.wantErr)
			}
			if len(got) != len(tt.want) {
				t.Errorf("Apply() = %v, want %v", got, tt.want)
			}
			for _, rel := range tt.wantFiles {
				if c.RootApparmor.Join(rel).NotExist() {
					t.Errorf("Apply() = %v, want %v", rel, "exist")
				}
			}
			for _, rel := range tt.wantLinks {
				if isLink, err := c.RootApparmor.Join(rel).IsSymlink(); err != nil || !isLink {
					t.Errorf("Apply() IsSymlink(%v) = %v, %v, want true", rel, isLink, err)
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
