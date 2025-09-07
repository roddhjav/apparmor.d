// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package prebuild

import "github.com/roddhjav/apparmor.d/pkg/paths"

var (
	// AppArmor ABI version
	ABI = 0

	// AppArmor version
	Version = 4.0

	// Tells the build we are a downstream project using apparmor.d as dependency
	DownStream = false

	// Either or not RBAC is enabled
	RBAC = false

	// Pkgname is the name of the package
	Pkgname = "apparmor.d"

	// Root is the root directory for the build (default: .build)
	Root *paths.Path = paths.New(".build")

	// RootApparmord is the final built apparmor.d directory (default: .build/apparmor.d)
	RootApparmord *paths.Path = Root.Join("apparmor.d")

	// DistDir is the directory where the distribution specific files are stored
	DistDir *paths.Path = paths.New("dists")

	// FlagDir is the directory where the flags are stored
	FlagDir *paths.Path = DistDir.Join("flags")

	// IgnoreDir is the directory where the ignore files are stored
	IgnoreDir *paths.Path = DistDir.Join("ignore")

	// SystemdDir is the directory where the systemd drop-in files are stored
	SystemdDir *paths.Path = paths.New("systemd")

	// DebianDir is the directory where the debian specific files are stored
	DebianDir *paths.Path = paths.New("debian")

	// DebianHide is the path to the debian/apparmor.d.hide file
	DebianHide = DebianHider{path: DebianDir.Join("apparmor.d.hide")}

	Ignore = Ignorer{}
	Flags  = Flagger{}
)
