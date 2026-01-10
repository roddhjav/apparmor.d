// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package tasks

import "github.com/roddhjav/apparmor.d/pkg/paths"

type TaskConfig struct {
	// AppArmor ABI version
	ABI int

	// AppArmor version
	Version float64

	// Tells the build we are a downstream project using apparmor.d as dependency
	DownStream bool

	// Either or not RBAC is enabled
	RBAC bool

	// Either or not we are in test mode
	Test bool

	// The dbus implementation used (true for dbus-daemon, false for dbus-broker)
	DbusDaemon bool

	// Pkgname is the name of the package
	Pkgname string

	// Root is the root directory for the runner (e.g. .build)
	Root *paths.Path

	// RootApparmor is the source apparmor.d directory (e.g. .build/apparmor.d)
	RootApparmor *paths.Path
}

func NewTaskConfig(root *paths.Path) *TaskConfig {
	return &TaskConfig{
		ABI:          0,
		Version:      4.0,
		DownStream:   false,
		RBAC:         false,
		Test:         false,
		DbusDaemon:   true,
		Pkgname:      "apparmor.d",
		Root:         root,
		RootApparmor: root.Join("apparmor.d"),
	}
}
