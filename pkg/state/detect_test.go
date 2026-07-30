// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package state

import (
	"errors"
	"strings"
	"testing"

	"github.com/roddhjav/apparmor.d/pkg/paths"
	"github.com/roddhjav/apparmor.d/pkg/state/detector"
)

func TestEnabled(t *testing.T) {
	tests := []struct {
		name string
		abi  string // base-strict abi rule; "" means no file
		want bool
	}{
		{name: "abi 5", abi: "abi <abi/5.0>,\n", want: true},
		{name: "abi 4", abi: "abi <abi/4.0>,\n", want: false},
		{name: "no abi", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sys := &detector.System{
				Root:     paths.New(t.TempDir()),
				Apparmor: paths.New(t.TempDir()),
			}
			if tt.abi != "" {
				base := sys.Apparmor.Join("abstractions/base-strict")
				if err := base.Parent().MkdirAll(); err != nil {
					t.Fatalf("mkdir: %v", err)
				}
				if err := base.WriteFile([]byte(tt.abi)); err != nil {
					t.Fatalf("write base-strict: %v", err)
				}
			}
			if got := Enabled(sys); got != tt.want {
				t.Errorf("Enabled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDetect(t *testing.T) {
	tests := []struct {
		name       string
		state      string // src state file content; "" means missing
		wantErr    bool
		wantDropin string // dropin content after Save
	}{
		{
			name: "detected variables written to dropin",
			state: "# AppArmor version\n@{VERSION} = 5\n\n" +
				"# No detector for this one\n@{NO_DETECTOR_HERE} = x\n\n" +
				"include if exists <tunables/multiarch.d/state.d>\n",
			wantDropin: "@{VERSION} = 4\n",
		},
		{
			name:       "undetected value skipped",
			state:      "@{DS} = wayland\n",
			wantDropin: "",
		},
		{
			name:    "missing state file",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := paths.New(t.TempDir())
			src := root.Join("state")
			if tt.state != "" {
				if err := src.WriteFile([]byte(tt.state)); err != nil {
					t.Fatalf("write state: %v", err)
				}
			}
			sys := &detector.System{
				Root:     paths.New(t.TempDir()),
				Apparmor: paths.New(t.TempDir()),
				Run: func(name string, arg ...string) (string, error) {
					if name == "apparmor_parser" {
						return "AppArmor parser version 4.0.2\n", nil
					}
					return "", errors.New("stubbed: " + name)
				},
			}

			dst := root.Join("state.d/aa-install")
			file, err := Detect(sys, src, dst)
			if (err != nil) != tt.wantErr {
				t.Errorf("Detect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if err := file.Save(); err != nil {
				t.Fatalf("Save() error = %v", err)
			}
			got, err := dst.ReadFileAsString()
			if err != nil {
				t.Fatalf("read dropin: %v", err)
			}
			if strings.TrimSpace(got) != strings.TrimSpace(tt.wantDropin) {
				t.Errorf("Detect() dropin = %q, want %q", got, tt.wantDropin)
			}
		})
	}
}
