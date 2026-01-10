// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2023-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

// Package cli provides the command line interface for prebuilding apparmor.d profiles.
// It is separated from the main package as it is also used by downstream projects.

package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/roddhjav/apparmor.d/pkg/builder"
	"github.com/roddhjav/apparmor.d/pkg/configure"
	"github.com/roddhjav/apparmor.d/pkg/directive"
	"github.com/roddhjav/apparmor.d/pkg/logging"
	"github.com/roddhjav/apparmor.d/pkg/paths"
	"github.com/roddhjav/apparmor.d/pkg/runtime"
	"github.com/roddhjav/apparmor.d/pkg/tasks"
)

const (
	nilABI   = 0
	nilVer   = 0.0
	nilBuild = ""
	nilSrc   = ""
	usage    = `aa-prebuild [-h] [--status] [--abi 3|4|5] [--version V] [--fsp] [--src DIR] [--buildir DIR]

    Prebuild apparmor.d profiles for a given distribution and apply
    internal built-in directives.

Options:
    -h, --help        Show this help message and exit.
    -s, --status      Show the status of enabled build tasks.
    -a, --abi ABI     Target apparmor ABI.
    -v, --version V   Target apparmor version.
    -f, --fsp         Configure AppArmor for full system policy and RBAC.
    -S, --src DIR     Profile source directory (default: apparmor.d/).
    -b, --buildir DIR Destination root build directory (default: .build/).
        --test        Enable test mode.
        --debug       Enable debug mode.
`
)

var (
	help    bool
	status  bool
	fsp     bool
	debug   bool
	test    bool
	abi     int
	version float64
	src     string
	buildir string
)

func init() {
	flag.BoolVar(&help, "h", false, "Show this help message and exit.")
	flag.BoolVar(&help, "help", false, "Show this help message and exit.")
	flag.BoolVar(&status, "s", false, "Show the status of enabled build tasks.")
	flag.BoolVar(&status, "status", false, "Show the status of enabled build tasks.")
	flag.BoolVar(&fsp, "f", false, "Configure AppArmor for full system policy and RBAC.")
	flag.BoolVar(&fsp, "fsp", false, "Configure AppArmor for full system policy and RBAC.")
	flag.IntVar(&abi, "a", nilABI, "Target apparmor ABI.")
	flag.IntVar(&abi, "abi", nilABI, "Target apparmor ABI.")
	flag.Float64Var(&version, "v", nilVer, "Target apparmor version.")
	flag.Float64Var(&version, "version", nilVer, "Target apparmor version.")
	flag.StringVar(&src, "S", nilSrc, "Profile source directory.")
	flag.StringVar(&src, "src", nilSrc, "Profile source directory.")
	flag.StringVar(&buildir, "b", nilBuild, "Destination root build directory.")
	flag.StringVar(&buildir, "buildir", nilBuild, "Destination root build directory.")
	flag.BoolVar(&debug, "debug", false, "Enable debug mode.")
	flag.BoolVar(&test, "test", false, "Enable test mode.")
}

func GetPrebuildRoot() *paths.Path {
	if buildir != nilBuild {
		return paths.New(buildir)
	}
	return paths.New(".build")
}

func Configure(r *runtime.Runners) *runtime.Runners {
	flag.Usage = func() { fmt.Print(usage) }
	flag.Parse()
	if help {
		flag.Usage()
		os.Exit(0)
	}

	// Register all directives (always available)
	r.Directives.
		Register(directive.NewDbus()).
		Register(directive.NewExec()).
		Register(directive.NewFilterOnly()).
		Register(directive.NewFilterExclude()).
		Register(directive.NewProfile()).
		Register(directive.NewRestart()).
		Register(directive.NewStack())

	if fsp && paths.New("apparmor.d/groups/_full").Exist() {
		r.Configures.Add(configure.NewFullSystemPolicy())
		r.Builders.Add(builder.NewFSP())
		r.RBAC = true
	}

	if abi != nilABI {
		r.ABI = abi
	}
	switch r.ABI {
	case 3:
		r.Builders.
			Add(builder.NewABI3()).      // Convert all profiles from abi 4.0 to abi 3.0
			Add(builder.NewAPPARMOR40()) // Convert all profiles from apparmor 4.1 to 4.0 or less

	case 4:
		// priority support was added in 4.1
		if r.Version == 4.0 {
			r.Builders.Add(builder.NewAPPARMOR40())
		}

		// Re-attach disconnected path
		if tasks.Distribution == "ubuntu" && r.Version >= 4.1 {
			// Ignored on ubuntu 25.04+ due to a memory leak that fully prevent
			// profiles compilation with re-attached paths.
			// See https://bugs.launchpad.net/ubuntu/+source/linux/+bug/2098730

			// Use stacked-dbus builder to resolve dbus rules
			r.Builders.Add(builder.NewStackedDbus())

		} else {
			if !r.DownStream {
				r.Configures.Add(configure.NewAttach())
			}
			r.Builders.Add(builder.NewAttach())
		}

	case 5:
		r.Builders.Add(builder.NewABI5()) // Convert all profiles from abi 4.0 to abi 5.0

		// Re-attach disconnected path
		if tasks.Distribution == "ubuntu" {
			// Ignored on ubuntu 25.04+ due to a memory leak that fully prevent
			// profiles compilation with re-attached paths.
			// See https://bugs.launchpad.net/ubuntu/+source/linux/+bug/2098730

			// Use stacked-dbus builder to resolve dbus rules
			r.Builders.Add(builder.NewStackedDbus())

		} else {
			if !r.DownStream {
				r.Configures.Add(configure.NewAttach())
			}
			r.Builders.Add(builder.NewAttach())

			// Fix dbus rules for dbus-broker
			r.Builders.Add(builder.NewDbusBroker())
			r.DbusDaemon = false
		}

	default:
		logging.Fatal("Invalid ABI version: %d", r.ABI)
	}

	if version != nilVer {
		r.Version = version
	}

	if status {
		fmt.Printf("%s\n%s\n%s",
			r.Configures.Help("Enabled configure"),
			r.Builders.Help("Enabled build"),
			r.Directives.Usage(),
		)
		os.Exit(0)
	}
	return r
}

func Prebuild(r *runtime.Runners) {
	logging.Step("Building apparmor.d profiles for %s", tasks.Distribution)
	logging.Success("AppArmor ABI targeted: %d", r.ABI)
	logging.Success("AppArmor version targeted: %.1f", r.Version)
	if r.Test {
		logging.Warning("Test mode enabled")
	}
	if fsp {
		logging.Success("Full system policy enabled")
	}
	if err := r.Configure(); err != nil {
		logging.Fatal("%s", err.Error())
	}
	if err := r.Build(); err != nil {
		logging.Fatal("%s", err.Error())
	}
}
