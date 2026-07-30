// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package configure

import (
	"errors"
	"strings"
	"testing"

	"github.com/roddhjav/apparmor.d/pkg/paths"
	"github.com/roddhjav/apparmor.d/pkg/state/detector"
)

const stateTunable = `# System states

# AppArmor version
@{VERSION} = 5

include if exists <tunables/multiarch.d/state.d>
`

func TestSetState_Apply(t *testing.T) {
	tests := []struct {
		name            string
		hasState        bool
		abi             string // abi version of abstractions/base-strict, "" means abi 5.0
		parserOut       string
		wantEmpty       bool   // Apply returns no result lines and no dropin
		wantContains    string // substring of the written dropin
		wantNotContains string // substring the dropin must not have
		wantErr         bool
	}{
		{
			name:         "writes detected values to the dropin",
			hasState:     true,
			parserOut:    "AppArmor parser version 4.0.2\n",
			wantContains: "@{VERSION} = 4",
		},
		{
			name:      "missing state file is a no-op",
			hasState:  false,
			wantEmpty: true,
		},
		{
			name:      "abi 4 tree is a no-op",
			hasState:  true,
			abi:       "4.0",
			parserOut: "AppArmor parser version 4.0.2\n",
			wantEmpty: true,
		},
		{
			name:            "failed detection skips the variable",
			hasState:        true,
			parserOut:       "garbage",
			wantNotContains: "@{VERSION}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newTaskConfigTmp(t)
			statePath := c.RootApparmor.Join(stateRel)
			if tt.hasState {
				abi := tt.abi
				if abi == "" {
					abi = "5.0"
				}
				seedFiles(t, c.RootApparmor, map[string]string{
					stateRel:                   stateTunable,
					"abstractions/base-strict": "abi <abi/" + abi + ">,\n",
				})
			}
			oldRoot, oldRun := detector.Root, detector.Run
			detector.Root = paths.New(t.TempDir())
			detector.Run = func(name string, arg ...string) (string, error) {
				if name == "apparmor_parser" {
					return tt.parserOut, nil
				}
				return "", errors.New("stubbed")
			}
			t.Cleanup(func() { detector.Root, detector.Run = oldRoot, oldRun })

			task := NewSetState()
			task.SetConfig(c)
			got, err := task.Apply()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Apply() error = %v, wantErr %v", err, tt.wantErr)
			}
			dropin := c.RootApparmor.Join(stateDropinRel)
			if tt.wantEmpty {
				if len(got) != 0 {
					t.Errorf("Apply() = %v, want empty", got)
				}
				if dropin.Exist() {
					t.Errorf("Apply() dropin = %v, want %v", "created", "absent")
				}
				return
			}
			content, err := dropin.ReadFileAsString()
			if err != nil {
				t.Fatalf("read dropin: %v", err)
			}
			if tt.wantContains != "" && !strings.Contains(content, tt.wantContains) {
				t.Errorf("dropin = %q, want contains %q", content, tt.wantContains)
			}
			if tt.wantNotContains != "" && strings.Contains(content, tt.wantNotContains) {
				t.Errorf("dropin = %q, want without %q", content, tt.wantNotContains)
			}
			state, err := statePath.ReadFileAsString()
			if err != nil {
				t.Fatalf("read state: %v", err)
			}
			if state != stateTunable {
				t.Errorf("state file = %q, want untouched", state)
			}
		})
	}
}
