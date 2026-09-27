// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package configure

import (
	"testing"

	"github.com/roddhjav/apparmor.d/pkg/paths"
)

func TestSystemdFSP_Apply(t *testing.T) {
	tests := []struct {
		name    string
		vendor  string
		admin   string
		want    map[string]string // drop-in path relative to the systemd root -> content
		absent  []string          // drop-in paths that must not exist
		wantErr bool
	}{
		{
			name:   "writes one drop-in per unit",
			vendor: "# comment\napt-news &apt_news\nsystemd-coredump@.service &systemd-coredump\n",
			want: map[string]string{
				"system/apt-news.service.d/apparmor.conf":          "[Service]\nAppArmorProfile=&apt_news\n",
				"system/systemd-coredump@.service.d/apparmor.conf": "[Service]\nAppArmorProfile=&systemd-coredump\n",
			},
		},
		{
			name:   "admin file overrides the vendor one",
			vendor: "apt-news &apt_news\ncolord &colord\n",
			admin:  "apt-news apt-news.service\n",
			want: map[string]string{
				"system/apt-news.service.d/apparmor.conf": "[Service]\nAppArmorProfile=apt-news.service\n",
			},
			absent: []string{"system/colord.service.d/apparmor.conf"},
		},
		{
			name:    "invalid entry",
			vendor:  "apt-news\n",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newTaskConfigTmp(t)
			vendor := paths.New(t.TempDir())
			admin := paths.New(t.TempDir())
			seedFiles(t, vendor, map[string]string{"00-main.conf": tt.vendor})
			if tt.admin != "" {
				seedFiles(t, admin, map[string]string{"00-main.conf": tt.admin})
			}
			seedFiles(t, c.Root.Join(SystemdFSPRel), map[string]string{
				"system/stale.service.d/apparmor.conf": "[Service]\n",
			})

			task := NewSystemdFSP(paths.PathList{vendor, admin})
			task.SetConfig(c)
			_, err := task.Apply()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Apply() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			root := c.Root.Join(SystemdFSPRel)
			for rel, want := range tt.want {
				got, err := root.Join(rel).ReadFileAsString()
				if err != nil {
					t.Fatalf("read %s: %v", rel, err)
				}
				if got != want {
					t.Errorf("Apply() %s = %q, want %q", rel, got, want)
				}
			}
			for _, rel := range append(tt.absent, "system/stale.service.d/apparmor.conf") {
				if root.Join(rel).Exist() {
					t.Errorf("Apply() %s = %v, want %v", rel, "present", "absent")
				}
			}
		})
	}
}
