// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

// Package prebuild defines the directory layout, manifest readers, and
// shared state used by the prebuild pipeline that turns the source profiles
// into a distribution-ready tree.
package prebuild

import "github.com/roddhjav/apparmor.d/pkg/paths"

var (
	// DistDir is the directory where the distribution specific files are stored
	DistDir *paths.Path = paths.New("dists")

	// FlagDir is the directory where the distribution flags are stored. It
	// mirrors the user flags.d structure: <name>.conf drop-in files.
	FlagDir *paths.Path = DistDir.Join("flags.d")

	// IgnoreDir is the directory where the distribution ignore files are
	// stored. It mirrors the user ignore.d structure: <name>.conf drop-ins.
	IgnoreDir *paths.Path = DistDir.Join("ignore.d")

	// SystemdDir is the directory where the systemd drop-in files are stored
	SystemdDir *paths.Path = paths.New("systemd")

	// DebianDir is the directory where the debian specific files are stored
	DebianDir *paths.Path = paths.New("debian")

	// DebianHide is the path to the debian/common.hide file
	DebianHide = DebianHider{path: DebianDir.Join("apparmor.d.hide")}

	Ignore = Ignorer{}
	Flags  = Flagger{}
)
