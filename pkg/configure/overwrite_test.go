// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package configure

import (
	"os"
	"testing"
)

func TestOverwriteFromLinks_Apply(t *testing.T) {
	tests := []struct {
		name        string
		files       map[string]string // root relative file -> content
		links       map[string]string // root relative link -> target
		roRoot      bool
		want        []string
		wantErr     bool
		wantFiles   []string
		wantLinks   []string // still present as symlinks (Lstat, may dangle)
		wantNoFiles []string
	}{
		{
			name:        "rename profile with disable link",
			files:       map[string]string{"hostname": "profile hostname {}\n"},
			links:       map[string]string{"disable/hostname": "../hostname"},
			want:        []string{"hostname"},
			wantFiles:   []string{"hostname.apparmor.d"},
			wantLinks:   []string{"disable/hostname"},
			wantNoFiles: []string{"hostname"},
		},
		{
			name:      "link without matching profile",
			files:     map[string]string{"other": "profile other {}\n"},
			links:     map[string]string{"disable/hostname": "../hostname"},
			wantFiles: []string{"other"},
		},
		{
			name:      "no disable directory",
			files:     map[string]string{"hostname": "profile hostname {}\n"},
			wantFiles: []string{"hostname"},
		},
		{
			name:    "rename error",
			files:   map[string]string{"hostname": "profile hostname {}\n"},
			links:   map[string]string{"disable/hostname": "../hostname"},
			roRoot:  true,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newTaskConfigTmp(t)
			seedFiles(t, c.RootApparmor, tt.files)
			for rel, target := range tt.links {
				link := c.RootApparmor.Join(rel)
				if err := link.Parent().MkdirAll(); err != nil {
					t.Fatalf("mkdir %s: %v", link.Parent(), err)
				}
				if err := os.Symlink(target, link.String()); err != nil {
					t.Fatalf("symlink %s: %v", link, err)
				}
			}
			if tt.roRoot {
				chmodRO(t, c.RootApparmor)
			}

			task := NewOverwriteFromLinks()
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
