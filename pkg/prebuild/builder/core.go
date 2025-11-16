// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package builder

import (
	"fmt"
	"strings"

	"github.com/roddhjav/apparmor.d/pkg/aa"
	"github.com/roddhjav/apparmor.d/pkg/paths"
	"github.com/roddhjav/apparmor.d/pkg/prebuild"
)

var (
	// Build the profiles with the following directive applied
	Builds = []Builder{}

	// Available builders
	Builders = map[string]Builder{}
)

// Builder main directive interface
type Builder interface {
	prebuild.BaseInterface
	Apply(opt *Option, profile string) (string, error)
}

// Option for a builder
type Option struct {
	Name string
	File *paths.Path
	Kind aa.FileKind
}

func NewOption(file *paths.Path) *Option {
	return &Option{
		Name: strings.TrimSuffix(file.Base(), ".apparmor.d"),
		File: file,
		Kind: aa.KindFromPath(file),
	}
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

func Run(file *paths.Path, profile string) (string, error) {
	var err error
	opt := NewOption(file)
	for _, b := range Builds {
		profile, err = b.Apply(opt, profile)
		if err != nil {
			return "", fmt.Errorf("%s %s: %w", b.Name(), opt.File, err)
		}
	}
	return profile, nil
}
