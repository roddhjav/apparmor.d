// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package detector

import "github.com/roddhjav/apparmor.d/pkg/tasks"

func init() {
	Register(
		osRelease{name: "OS_FAMILY"},
		osRelease{name: "OS_ID", key: "ID"},
		osRelease{name: "OS_VERSION_ID", key: "VERSION_ID"},
	)
}

// osRelease detects an @{OS_*} variable from os-release. With no key it
// reports the distribution family.
type osRelease struct {
	name string
	key  string
}

func (d osRelease) Name() string { return d.name }

func (d osRelease) Detect(*System) []string {
	value := tasks.Family
	if d.key != "" {
		value = tasks.Release[d.key]
	}
	if value == "" {
		return nil
	}
	return []string{value}
}
