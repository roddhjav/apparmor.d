// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package prebuild

import (
	"os"
	"strings"

	"github.com/roddhjav/apparmor.d/pkg/paths"
)

var (
	// AppArmor ABI version
	ABI uint = 0

	// Root is the root directory for the build (default: ./.build)
	Root *paths.Path = getRootBuild()

	// RootApparmord is the final built apparmor.d directory (default: .build/apparmor.d)
	RootApparmord *paths.Path = Root.Join(Src)

	// src is the basename of the source directory (default: apparmor.d)
	Src = "apparmor.d"

	// SrcApparmord is the source apparmor.d directory (default: ./apparmor.d)
	SrcApparmord *paths.Path = paths.New(Src)

	// DistDir is the directory where the distribution specific files are stored
	DistDir *paths.Path = paths.New("dists")

	// FlagDir is the directory where the flags are stored
	FlagDir *paths.Path = DistDir.Join("flags")

	// IgnoreDir is the directory where the ignore files are stored
	IgnoreDir *paths.Path = DistDir.Join("ignore")

	// PkgDir is the directory where the packages files are stored
	PkgDir *paths.Path = DistDir.Join("packages")

	// SystemdDir is the directory where the systemd drop-in files are stored
	SystemdDir *paths.Path = paths.New("systemd")

	// DebianDir is the directory where the debian specific files are stored
	DebianDir *paths.Path = paths.New("debian")

	// DebianHide is the path to the debian/apparmor.d.hide file
	DebianHide = DebianHider{path: DebianDir.Join("apparmor.d.hide")}

	// Packages are the packages to build
	Packages = getPackages()

	Ignore = Ignorer{}
	Flags  = Flagger{}
)

func getRootBuild() *paths.Path {
	root, present := os.LookupEnv("BUILD")
	if !present {
		root = ".build"
	}
	return paths.New(root)
}

func getPackages() []string {
	files, err := PkgDir.ReadDirRecursiveFiltered(nil, paths.FilterOutDirectories())
	if err != nil {
		return []string{}
	}
	packages := make([]string, 0, len(files))
	for _, file := range files {
		packages = append(packages, strings.TrimSuffix(file.Base(), ".conf"))
	}
	return packages
}
