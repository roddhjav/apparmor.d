// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package prebuild

import (
	"github.com/roddhjav/apparmor.d/pkg/util"
)

type Flagger struct{}

func (f Flagger) Read(name string) map[string][]string {
	return util.ReadFlagsFile(FlagDir.Join(name + ".conf"))
}

type Ignorer struct{}

func (i Ignorer) Read(name string) []string {
	path := IgnoreDir.Join(name + ".conf")
	if !path.Exist() {
		return []string{}
	}
	return path.MustReadFilteredFileAsLines()
}
