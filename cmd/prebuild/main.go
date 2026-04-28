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
	case "arch":

	case "ubuntu":
		switch tasks.Release["VERSION_CODENAME"] {
		case "jammy":
			c.ABI = 3
			c.Version = 3.0
		case "noble":
			c.ABI = 4
			c.Version = 4.0
		case "questing":
			c.ABI = 4
			c.Version = 4.0
		case "resolute":
			c.ABI = 5
			c.Version = 5.0
		}

	case "debian":
		switch tasks.Release["VERSION_CODENAME"] {
		case "bullseye", "bookworm":
			c.ABI = 3
			c.Version = 3.0
		}

	case "whonix":
		c.ABI = 3
		c.Version = 3.0

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

		// Merge profiles (from group/, profiles-*-*/) to a unified apparmor.d directory
		Add(configure.NewMerge()).

		// Set distribution specificities
		Add(configure.NewConfigure()).

		// Overwrite dummy upstream profile
		Add(configure.NewOverwrite(false)).

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
