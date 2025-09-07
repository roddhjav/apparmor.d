// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2023-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package main

import (
	"github.com/roddhjav/apparmor.d/pkg/prebuild"
	"github.com/roddhjav/apparmor.d/pkg/prebuild/builder"
	"github.com/roddhjav/apparmor.d/pkg/prebuild/cli"
	"github.com/roddhjav/apparmor.d/pkg/prebuild/prepare"
)

// Cli arguments have priority over the settings entered here
func init() {
	// Define the default ABI
	prebuild.ABI = 4

	// Define the default version
	prebuild.Version = 4.1

	// Define the tasks applied by default
	prepare.Register(
		"synchronise",     // Initialize a new clean apparmor.d build directory
		"ignore",          // Ignore profiles and files from dist/ignore
		"merge",           // Merge profiles (from group/, profiles-*-*/) to a unified apparmor.d directory
		"configure",       // Set distribution specificities
		"setflags",        // Set flags as definied in dist/flags
		"overwrite",       // Overwrite dummy upstream profiles
		"systemd-default", // Set systemd unit drop in files for dbus profiles
	)

	// Build tasks applied by default
	builder.Register(
		"userspace",   // Resolve variable in profile attachments
		"hotfix",      // Temporary fix for #74, #80 & #235
		"base-strict", // Use base-strict as base abstraction
	)

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
	cli.Configure()
	cli.Prebuild()
}
