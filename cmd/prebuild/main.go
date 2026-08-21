// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2023-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package main

import (
	"github.com/roddhjav/apparmor.d/pkg/builder"
	"github.com/roddhjav/apparmor.d/pkg/configure"
	"github.com/roddhjav/apparmor.d/pkg/paths"
	"github.com/roddhjav/apparmor.d/pkg/prebuild/cli"
	"github.com/roddhjav/apparmor.d/pkg/runtime"
	"github.com/roddhjav/apparmor.d/pkg/tasks"
)

// Cli arguments have priority over the settings entered here
func configInit() *tasks.TaskConfig {
	c := tasks.NewTaskConfig(cli.GetPrebuildRoot())

	// Define the default ABI
	c.ABI = 4

	// Define the default version
	c.Version = 4.1

	// Matrix of ABI/Apparmor version to integrate with
	switch tasks.Distribution {
	case "opensuse":
		c.ABI = 4
		c.Version = 5.0

	case "arch":

	case "ubuntu":
		switch tasks.Release["VERSION_CODENAME"] {
		case "resolute":
			c.ABI = 4
			c.Version = 5.0
		default:
			panic("Unsupported Ubuntu version: " + tasks.Release["VERSION_CODENAME"])
		}

	case "debian":
		switch tasks.Release["VERSION_CODENAME"] {
		case "trixie":
			c.ABI = 4
			c.Version = 4.1
		case "forky":
			c.ABI = 5
			c.Version = 5.0
		default:
			panic("Unsupported Debian version: " + tasks.Release["VERSION_CODENAME"])
		}

	}
	return c
}

func main() {
	cli.ParseFlags()
	c := configInit()
	r := runtime.NewRunners(c)

	// Add default configure tasks
	r.Configures.
		// Initialize a new clean apparmor.d build directory
		Add(configure.NewSynchronise(
			[]*paths.Path{paths.New("apparmor.d"), paths.New("share")},
		)).

		// Ignore profiles and files from dist/ignore
		Add(configure.NewIgnore()).

		// Set distribution specificities
		Add(NewConfigure()).

		// Set systemd unit drop in files for dbus profiles
		Add(configure.NewSystemdDefault())

	// Default build tasks
	r.Builders.
		// Resolve variable in profile attachments
		Add(builder.NewUserspace()).

		// Temporary fix for #74, #80 & #235
		Add(builder.NewHotFix()).

		// Use base-strict as base abstraction
		Add(builder.NewBaseStrict())

	r = cli.Configure(r)
	cli.Prebuild(r)
}
