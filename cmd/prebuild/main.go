// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2023-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package main

import (
	"github.com/roddhjav/apparmor.d/pkg/builder"
	"github.com/roddhjav/apparmor.d/pkg/configure"
	"github.com/roddhjav/apparmor.d/pkg/paths"
	"github.com/roddhjav/apparmor.d/pkg/prebuild"
	"github.com/roddhjav/apparmor.d/pkg/prebuild/cli"
	"github.com/roddhjav/apparmor.d/pkg/run"
	"github.com/roddhjav/apparmor.d/pkg/tasks"
)

// Cli arguments have priority over the settings entered here
func init() {
	// Define the default ABI
	prebuild.ABI = 4

	// Define the default version
	prebuild.Version = 4.1

	// Matrix of ABI/Apparmor version to integrate with
	switch prebuild.Distribution {
	case "arch":

	case "ubuntu":
		switch prebuild.Release["VERSION_CODENAME"] {
		case "jammy":
			prebuild.ABI = 3
			prebuild.Version = 3.0
		case "noble":
			prebuild.ABI = 4
			prebuild.Version = 4.0
		case "questing":
			prebuild.ABI = 4
			prebuild.Version = 5.0
		case "resolute":
			prebuild.ABI = 4
			prebuild.Version = 5.0
		}

	case "debian":
		switch prebuild.Release["VERSION_CODENAME"] {
		case "bullseye", "bookworm":
			prebuild.ABI = 3
			prebuild.Version = 3.0
		}

	case "whonix":
		prebuild.ABI = 3
		prebuild.Version = 3.0

		// Hide rewritten Whonix profiles
		prebuild.Hide += `/etc/apparmor.d/abstractions/base.d/kicksecure
		/etc/apparmor.d/home.tor-browser.firefox
		/etc/apparmor.d/tunables/homsanitycheck
		/etc/apparmor.d/usr.bin.url_e.d/anondist
		/etc/apparmor.d/tunables/home.d/live-mode
		/etc/apparmor.d/tunables/home.d/qubes-whonix-anondist
		/etc/apparmor.d/usr.bin.hexchat
		/etc/apparmor.d/usr.bin.sdwdate
		/etc/apparmor.d/usr.bin.systemcheck
		/etc/apparmor.d/usr.bin.timeto_unixtime
		/etc/apparmor.d/whonix-firewall
		`
	}
}

func main() {
	c := tasks.NewTaskConfig(cli.GetPrebuildRoot())
	r := run.NewRunners(c)

	// Add default configure tasks
	r.Configures.
		// Initialize a new clean apparmor.d build directory
		Add(configure.NewSynchronise(
			[]*paths.Path{paths.New("apparmor.d"), paths.New("share")},
		)).

		// Ignore profiles and files from dist/ignore
		Add(configure.NewIgnore()). // TODO: Keep it here, have one in aa-install, as well as a Include

		// Set distribution specificities
		Add(configure.NewConfigure()).
		// Add(configure.NewSetFlags()). // Set flags as definied in dist/flags

		// Overwrite dummy upstream profile
		Add(configure.NewOverwrite(false)). // TODO: Move in aa-install

		// Set systemd unit drop in files for dbus profiles
		Add(configure.NewSystemd())

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
