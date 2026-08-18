// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package main

import (
	"os"
	"strings"
	"testing"

	"github.com/roddhjav/apparmor.d/pkg/paths"
	"github.com/roddhjav/apparmor.d/pkg/tasks"
)

const systemTunable = "@{pci_bus}=pci@{hex4}:@{hex2}\n@{multiarch}=*-linux-gnu*\n"

// writeFile writes content to path, creating the parent directories.
func writeFile(t *testing.T, path *paths.Path, content string) {
	t.Helper()
	if err := path.Parent().MkdirAll(); err != nil {
		t.Fatalf("mkdir %s: %v", path.Parent(), err)
	}
	if err := path.WriteFile([]byte(content)); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

// setupConfigure builds a minimal apparmor.d tree and a Configure task
// bound to it, with the distribution and version set for the test.
func setupConfigure(t *testing.T, distribution string, version float64) *Configure {
	t.Helper()
	if err := os.MkdirAll("/tmp/tests", 0o755); err != nil {
		t.Fatalf("mkdir /tmp/tests: %v", err)
	}
	t.Setenv("TMPDIR", "/tmp/tests")

	cfg := tasks.NewTaskConfig(paths.New(t.TempDir()))
	cfg.Version = version
	files := []string{
		"fapp", "fbwrap", "wg", "dig", "free", "nslookup", "su", "sudo",
		"tunables/multiarch.d/base", "abstractions/devices-usb",
		"abstractions/devices-usb-read", "abstractions/nameservice-strict",
	}
	for _, name := range files {
		writeFile(t, cfg.RootApparmor.Join(name), "# file\n")
	}
	writeFile(t, cfg.RootApparmor.Join("tunables/multiarch.d/system"), systemTunable)

	old := tasks.Distribution
	tasks.Distribution = distribution
	t.Cleanup(func() { tasks.Distribution = old })

	task := NewConfigure()
	task.SetConfig(cfg)
	return task
}

func TestConfigure_Apply(t *testing.T) {
	tests := []struct {
		name         string
		distribution string
		version      float64
		wantErr      bool
		wantFiles    []string
		wantNoFiles  []string
		wantPciBus   bool
	}{
		{
			name:         "arch 5.0 drops upstreamed profiles",
			distribution: "arch",
			version:      5.0,
			wantErr:      false,
			wantFiles:    []string{"fapp", "fbwrap", "su", "sudo"},
			wantNoFiles:  []string{"wg", "dig", "free", "nslookup", "tunables/multiarch.d/base"},
			wantPciBus:   false,
		},
		// {
		// 	name:         "ubuntu 5.0 drops su and sudo",
		// 	distribution: "ubuntu",
		// 	version:      5.0,
		// 	wantErr:      false,
		// 	wantFiles:    []string{"fapp"},
		// 	wantNoFiles:  []string{"su", "sudo"},
		// 	wantPciBus:   false,
		// },
		{
			name:         "debian 4.0 keeps the upstreamed files",
			distribution: "debian",
			version:      4.0,
			wantErr:      false,
			wantFiles:    []string{"wg", "dig", "tunables/multiarch.d/base"},
			wantNoFiles:  []string{"fapp", "fbwrap"},
			wantPciBus:   true,
		},
		{
			name:         "unsupported distribution",
			distribution: "gentoo",
			version:      5.0,
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := setupConfigure(t, tt.distribution, tt.version)

			_, err := task.Apply()
			if (err != nil) != tt.wantErr {
				t.Errorf("Configure.Apply() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for _, name := range tt.wantFiles {
				if task.RootApparmor.Join(name).NotExist() {
					t.Errorf("Configure.Apply() = %v removed, want kept", name)
				}
			}
			for _, name := range tt.wantNoFiles {
				if task.RootApparmor.Join(name).Exist() {
					t.Errorf("Configure.Apply() = %v kept, want removed", name)
				}
			}
			if !tt.wantErr {
				out, err := task.RootApparmor.Join("tunables/multiarch.d/system").ReadFileAsString()
				if err != nil {
					t.Fatalf("read system tunable: %v", err)
				}
				if got := strings.Contains(out, "@{pci_bus}="); got != tt.wantPciBus {
					t.Errorf("Configure.Apply() @{pci_bus} = %v, want %v", got, tt.wantPciBus)
				}
			}
		})
	}
}
