// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package builder

import (
	"fmt"

	"github.com/roddhjav/apparmor.d/pkg/prebuild/cfg"
)

var (
	// Build the profiles with the following directive applied
	Builds = []Builder{}

	// Available builders
	Builders = map[string]Builder{}
)

// Main directive interface
type Builder interface {
	cfg.BaseInterface
	Apply(profile string) (string, error)
}

func Register(names ...string) {
	for _, name := range names {
		if b, present := Builders[name]; present {
			Builds = append(Builds, b)
		} else {
			panic(fmt.Sprintf("Unknown builder: %s", name))
		}
	}
}

func RegisterBuilder(d Builder) {
	Builders[d.Name()] = d
}
