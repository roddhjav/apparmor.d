// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package configure

import (
	"strings"
	"testing"

	"github.com/roddhjav/apparmor.d/pkg/paths"
)

const profilesTunable = `@{p_sd}=unconfined
@{p_colord}=colord
`

func TestFullSystemPolicy_Apply(t *testing.T) {
	tests := []struct {
		name         string
		enabled      bool
		wantGroup    bool     // groups/_full is kept
		wantContains []string // lines of the written dropin, nil means no dropin
		wantErr      bool
	}{
		{
			name:      "disabled removes the group",
			enabled:   false,
			wantGroup: false,
		},
		{
			name:      "enabled keeps the group and overrides profile names",
			enabled:   true,
			wantGroup: true,
			wantContains: []string{
				"@{p_sd} := sd\n",
				"@{p_colord} := {colord,sd//&colord,colord//&sd}\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newTaskConfigTmp(t)
			seedFiles(t, c.RootApparmor, map[string]string{fspRel + "/sd": "profile sd {}\n"})
			tunables := paths.New(t.TempDir()).Join("profiles")
			if err := tunables.WriteFile([]byte(profilesTunable)); err != nil {
				t.Fatal(err)
			}

			task := NewFullSystemPolicy(tt.enabled, tunables)
			task.SetConfig(c)
			_, err := task.Apply()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Apply() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got := c.RootApparmor.Join(fspRel).Exist(); got != tt.wantGroup {
				t.Errorf("Apply() group kept = %v, want %v", got, tt.wantGroup)
			}
			dropin := c.RootApparmor.Join(fspDropinRel)
			if tt.wantContains == nil {
				if dropin.Exist() {
					t.Errorf("Apply() dropin = %v, want %v", "created", "absent")
				}
				return
			}
			content, err := dropin.ReadFileAsString()
			if err != nil {
				t.Fatalf("read dropin: %v", err)
			}
			for _, want := range tt.wantContains {
				if !strings.Contains(content, want) {
					t.Errorf("dropin = %q, want contains %q", content, want)
				}
			}
		})
	}
}
