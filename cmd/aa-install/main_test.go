// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package main

import (
	"errors"
	"strings"
	"testing"

	"github.com/roddhjav/apparmor.d/pkg/aa"
	"github.com/roddhjav/apparmor.d/pkg/logging"
	"github.com/roddhjav/apparmor.d/pkg/paths"
	"github.com/roddhjav/apparmor.d/pkg/state/detector"
)

const keptProfile = `abi <abi/5.0>,

profile aa_test_kept @{bin}/true {
  include <abstractions/base>
}
`

const droppedProfile = `abi <abi/5.0>,

profile aa_test_dropped /usr/bin/aa-install-no-such-program {
  include <abstractions/base>
}
`

const overwriteProfile = `abi <abi/5.0>,

profile aa_test_over @{bin}/true {
  include <abstractions/base>
}
`

const upstreamProfile = "# upstream dummy profile\n"

// runEnv holds the directories and reload counter a run() test operates on.
type runEnv struct {
	configDir *paths.Path
	vendorDir *paths.Path
	srcDir    *paths.Path
	targetDir *paths.Path
	reloads   int
}

// setupRunEnv points the option globals at a fake source tree and stubs
// the AppArmor reload. All mutated globals are restored on cleanup.
func setupRunEnv(t *testing.T) *runEnv {
	t.Helper()
	env := &runEnv{}

	root := paths.New(t.TempDir())
	env.srcDir = root.Join("apparmor.d")
	writeFile(t, env.srcDir.Join("groups/apps/aa_test_kept"), keptProfile)
	writeFile(t, env.srcDir.Join("profiles-a-f/aa_test_dropped"), droppedProfile)
	writeFile(t, env.srcDir.Join("tunables/multiarch.d/state"),
		"@{VERSION} = 5\n\ninclude if exists <tunables/multiarch.d/state.d>\n")
	writeFile(t, env.srcDir.Join("abstractions/base-strict"), "abi <abi/5.0>,\n")

	env.configDir = paths.New(t.TempDir())
	env.targetDir = paths.New(t.TempDir())
	env.vendorDir = setVendorConfigDir(t)

	config = env.configDir.String()
	magic = env.targetDir.String()
	src = env.srcDir.String()

	oldMagic := aa.MagicRoot
	oldReload := reloadAppArmor
	reloadAppArmor = func() error {
		env.reloads++
		return nil
	}
	// Keep the state detectors off the real host: fake filesystem root,
	// no command execution.
	oldRoot, oldRun := detector.Root, detector.Run
	detector.Root = paths.New(t.TempDir())
	detector.Run = func(name string, arg ...string) (string, error) {
		return "", errors.New("stubbed")
	}
	t.Cleanup(func() {
		install, all, complain, enforce = false, false, false, false
		uninstall, status, list = false, false, false
		config, magic, src = nilConfig, nilMagic, nilSrc
		verbose, logging.Quiet = false, false
		aa.MagicRoot = oldMagic
		reloadAppArmor = oldReload
		detector.Root, detector.Run = oldRoot, oldRun
	})

	// The build directory is created with os.MkdirTemp and deliberately
	// not removed; keep test runs from littering the system /tmp.
	t.Setenv("TMPDIR", t.TempDir())
	return env
}

func TestRun(t *testing.T) {
	tests := []struct {
		name        string
		flags       func()                          // set the action option globals
		setup       func(t *testing.T, env *runEnv) // extra environment preparation
		wantErr     bool
		wantReloads int
		check       func(t *testing.T, env *runEnv)
	}{
		{
			name:    "conflicting modes",
			flags:   func() { complain, enforce = true, true },
			wantErr: true,
		},
		{
			name:        "install keeps installed programs only",
			flags:       func() { install = true },
			wantReloads: 1,
			check: func(t *testing.T, env *runEnv) {
				if !env.targetDir.Join("aa_test_kept").Exist() {
					t.Error("kept profile was not installed")
				}
				if env.targetDir.Join("aa_test_dropped").Exist() {
					t.Error("profile for a missing program was installed")
				}
			},
		},
		{
			name:        "install all keeps profiles of missing programs",
			flags:       func() { install, all = true, true },
			wantReloads: 1,
			check: func(t *testing.T, env *runEnv) {
				if !env.targetDir.Join("aa_test_kept").Exist() {
					t.Error("kept profile was not installed")
				}
				if !env.targetDir.Join("aa_test_dropped").Exist() {
					t.Error("profile for a missing program was not installed")
				}
			},
		},
		{
			name:  "install with include.d all mode",
			flags: func() { install = true },
			setup: func(t *testing.T, env *runEnv) {
				writeFile(t, env.configDir.Join("modes"), "include all\n")
			},
			wantReloads: 1,
			check: func(t *testing.T, env *runEnv) {
				if !env.targetDir.Join("aa_test_dropped").Exist() {
					t.Error("profile for a missing program was not installed")
				}
			},
		},
		{
			name:        "install complain",
			flags:       func() { install, complain = true, true },
			wantReloads: 1,
			check: func(t *testing.T, env *runEnv) {
				got, err := env.targetDir.Join("aa_test_kept").ReadFileAsString()
				if err != nil {
					t.Fatalf("read installed profile: %v", err)
				}
				if !strings.Contains(got, "flags=(complain)") {
					t.Errorf("installed profile = %q, want complain flag", got)
				}
			},
		},
		{
			name:        "install enforce",
			flags:       func() { install, enforce = true, true },
			wantReloads: 1,
			check: func(t *testing.T, env *runEnv) {
				got, err := env.targetDir.Join("aa_test_kept").ReadFileAsString()
				if err != nil {
					t.Fatalf("read installed profile: %v", err)
				}
				if strings.Contains(got, "complain") {
					t.Errorf("installed profile = %q, want no complain flag", got)
				}
			},
		},
		{
			name:  "install writes detected state",
			flags: func() { install = true },
			setup: func(t *testing.T, env *runEnv) {
				detector.Run = func(name string, arg ...string) (string, error) {
					if name == "apparmor_parser" {
						return "AppArmor parser version 4.0.2\n", nil
					}
					return "", errors.New("stubbed")
				}
			},
			wantReloads: 1,
			check: func(t *testing.T, env *runEnv) {
				got, err := env.targetDir.Join("tunables/multiarch.d/state.d/aa-install").ReadFileAsString()
				if err != nil {
					t.Fatalf("read installed state dropin: %v", err)
				}
				if !strings.Contains(got, "@{VERSION} = 4") {
					t.Errorf("installed state dropin = %q, want @{VERSION} = 4", got)
				}
			},
		},
		{
			name:  "install with include.d full mode",
			flags: func() { install = true },
			setup: func(t *testing.T, env *runEnv) {
				writeFile(t, env.configDir.Join("modes"), "include full\n")
				writeFile(t, env.srcDir.Join("profiles-a-f/aa_test_extra"),
					strings.ReplaceAll(keptProfile, "aa_test_kept", "aa_test_extra"))
				writeFile(t, env.configDir.Join("include.d/user.conf"), "aa_test_kept\n")
			},
			wantReloads: 1,
			check: func(t *testing.T, env *runEnv) {
				if !env.targetDir.Join("aa_test_kept").Exist() {
					t.Error("included profile was not installed")
				}
				if env.targetDir.Join("aa_test_extra").Exist() {
					t.Error("profile not listed in include.d was installed")
				}
			},
		},
		{
			name:  "install with include.d default mode",
			flags: func() { install = true },
			setup: func(t *testing.T, env *runEnv) {
				writeFile(t, env.configDir.Join("ignore.d/user.conf"), "aa_test_kept\n")
				writeFile(t, env.configDir.Join("include.d/user.conf"),
					"aa_test_kept\naa_test_dropped\n")
			},
			wantReloads: 1,
			check: func(t *testing.T, env *runEnv) {
				if !env.targetDir.Join("aa_test_kept").Exist() {
					t.Error("ignored included profile was not re-applied")
				}
				if !env.targetDir.Join("aa_test_dropped").Exist() {
					t.Error("included profile of a missing program was not installed")
				}
			},
		},
		{
			name:  "install invalid include mode config",
			flags: func() { install = true },
			setup: func(t *testing.T, env *runEnv) {
				writeFile(t, env.configDir.Join("modes"), "include bogus\n")
			},
			wantErr: true,
		},
		{
			name:  "install invalid mode config",
			flags: func() { install = true },
			setup: func(t *testing.T, env *runEnv) {
				writeFile(t, env.configDir.Join("modes"), "default bogus\n")
			},
			wantErr: true,
		},
		{
			name:  "install config overrides vendor mode",
			flags: func() { install = true },
			setup: func(t *testing.T, env *runEnv) {
				writeFile(t, env.vendorDir.Join("modes"), "default enforce\n")
				writeFile(t, env.configDir.Join("modes"), "default complain\n")
			},
			wantReloads: 1,
			check: func(t *testing.T, env *runEnv) {
				got, err := env.targetDir.Join("aa_test_kept").ReadFileAsString()
				if err != nil {
					t.Fatalf("read installed profile: %v", err)
				}
				if !strings.Contains(got, "flags=(complain)") {
					t.Errorf("installed profile = %q, want complain flag", got)
				}
			},
		},
		{
			name:  "install missing source",
			flags: func() { install = true },
			setup: func(t *testing.T, env *runEnv) {
				src = "/does/not/exist/apparmor.d"
			},
			wantErr: true,
		},
		{
			name:  "install broken tmpdir",
			flags: func() { install = true },
			setup: func(t *testing.T, env *runEnv) {
				t.Setenv("TMPDIR", "/does/not/exist")
			},
			wantErr: true,
		},
		{
			name:  "install unknown directive",
			flags: func() { install = true },
			setup: func(t *testing.T, env *runEnv) {
				writeFile(t, env.srcDir.Join("groups/apps/aa_test_bad"), `abi <abi/5.0>,

# #aa:bogus
profile aa_test_bad @{bin}/true {
  include <abstractions/base>
}
`)
			},
			wantErr: true,
		},
		{
			name:  "install reload error",
			flags: func() { install = true },
			setup: func(t *testing.T, env *runEnv) {
				reloadAppArmor = func() error { return errors.New("reload failed") }
			},
			wantErr: true,
		},
		{
			name:  "uninstall remove error",
			flags: func() { uninstall = true },
			setup: func(t *testing.T, env *runEnv) {
				install = true
				if err := run(); err != nil {
					t.Fatalf("run(install): %v", err)
				}
				install = false
				chmod(t, env.targetDir, 0o555)
			},
			wantErr:     true,
			wantReloads: 1, // from the install in setup
		},
		{
			name:  "install does not replace upstream profile",
			flags: func() { install = true },
			setup: func(t *testing.T, env *runEnv) {
				writeFile(t, env.srcDir.Join("profiles-m-r/aa_test_over"), overwriteProfile)
				writeFile(t, env.configDir.Join("overwrite.d/user.conf"), "aa_test_over\n")
				writeFile(t, env.targetDir.Join("aa_test_over"), upstreamProfile)
			},
			wantReloads: 1,
			check: func(t *testing.T, env *runEnv) {
				got, err := env.targetDir.Join("aa_test_over").ReadFileAsString()
				if err != nil {
					t.Fatalf("read upstream profile: %v", err)
				}
				if got != upstreamProfile {
					t.Errorf("upstream profile = %q, want untouched %q", got, upstreamProfile)
				}
				if !env.targetDir.Join("aa_test_over.apparmor.d").Exist() {
					t.Error("overwriting profile was not installed as aa_test_over.apparmor.d")
				}
				if isLink, err := env.targetDir.Join("disable/aa_test_over").IsSymlink(); err != nil || !isLink {
					t.Errorf("disable link IsSymlink() = %v, %v, want true", isLink, err)
				}
			},
		},
		{
			name:  "uninstall keeps upstream profile",
			flags: func() { uninstall = true },
			setup: func(t *testing.T, env *runEnv) {
				writeFile(t, env.srcDir.Join("profiles-m-r/aa_test_over"), overwriteProfile)
				writeFile(t, env.configDir.Join("overwrite.d/user.conf"), "aa_test_over\n")
				writeFile(t, env.targetDir.Join("aa_test_over"), upstreamProfile)
				install = true
				if err := run(); err != nil {
					t.Fatalf("run(install): %v", err)
				}
				install = false
			},
			wantReloads: 2,
			check: func(t *testing.T, env *runEnv) {
				got, err := env.targetDir.Join("aa_test_over").ReadFileAsString()
				if err != nil {
					t.Fatalf("read upstream profile: %v", err)
				}
				if got != upstreamProfile {
					t.Errorf("upstream profile = %q, want untouched %q", got, upstreamProfile)
				}
				if env.targetDir.Join("aa_test_over.apparmor.d").Exist() {
					t.Error("overwriting profile still installed after uninstall")
				}
				if _, err := env.targetDir.Join("disable/aa_test_over").Lstat(); err == nil {
					t.Error("disable link still installed after uninstall")
				}
			},
		},
		{
			name:  "uninstall removes dangling disable link",
			flags: func() { uninstall = true },
			setup: func(t *testing.T, env *runEnv) {
				writeFile(t, env.srcDir.Join("profiles-m-r/aa_test_over"), overwriteProfile)
				writeFile(t, env.configDir.Join("overwrite.d/user.conf"), "aa_test_over\n")
				install = true
				if err := run(); err != nil {
					t.Fatalf("run(install): %v", err)
				}
				install = false
			},
			wantReloads: 2,
			check: func(t *testing.T, env *runEnv) {
				if _, err := env.targetDir.Join("disable/aa_test_over").Lstat(); err == nil {
					t.Error("dangling disable link still installed after uninstall")
				}
			},
		},
		{
			name:  "status default action",
			flags: func() {},
		},
		{
			name:  "list",
			flags: func() { list = true },
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := setupRunEnv(t)
			if tt.setup != nil {
				tt.setup(t, env)
			}
			tt.flags()
			err := run()
			if (err != nil) != tt.wantErr {
				t.Fatalf("run() error = %v, wantErr %v", err, tt.wantErr)
			}
			if env.reloads != tt.wantReloads {
				t.Errorf("run() reloads = %d, want %d", env.reloads, tt.wantReloads)
			}
			if tt.check != nil {
				tt.check(t, env)
			}
		})
	}
}

// TestRun_Lifecycle exercises the full sequential user workflow:
// install, no-op reinstall, status, list, uninstall, empty uninstall.
func TestRun_Lifecycle(t *testing.T) {
	env := setupRunEnv(t)

	steps := []struct {
		name        string
		flags       func()
		wantReloads int
		check       func(t *testing.T)
	}{
		{
			name:        "install",
			flags:       func() { install = true },
			wantReloads: 1,
			check: func(t *testing.T) {
				if !env.targetDir.Join("aa_test_kept").Exist() {
					t.Error("kept profile was not installed")
				}
			},
		},
		{
			name:        "reinstall without changes",
			flags:       func() { install = true },
			wantReloads: 1, // unchanged: no new reload
		},
		{
			name:        "status explicit",
			flags:       func() { status = true },
			wantReloads: 1,
		},
		{
			name:        "uninstall",
			flags:       func() { uninstall = true },
			wantReloads: 2,
			check: func(t *testing.T) {
				if env.targetDir.Join("aa_test_kept").Exist() {
					t.Error("profile still installed after uninstall")
				}
			},
		},
		{
			name:        "uninstall with nothing installed",
			flags:       func() { uninstall = true },
			wantReloads: 2, // nothing removed: no new reload
		},
	}
	for _, step := range steps {
		install, all, complain, enforce = false, false, false, false
		uninstall, status, list = false, false, false
		step.flags()
		if err := run(); err != nil {
			t.Fatalf("run(%s) error = %v", step.name, err)
		}
		if env.reloads != step.wantReloads {
			t.Errorf("run(%s) reloads = %d, want %d", step.name, env.reloads, step.wantReloads)
		}
		if step.check != nil {
			step.check(t)
		}
	}
}
