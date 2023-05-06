// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2023 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package prebuild

import (
	"github.com/arduino/go-paths-helper"
)

var (
	Distribution  string
	DistDir       *paths.Path
	Root          *paths.Path
	RootApparmord *paths.Path
)

func init() {
	DistDir = paths.New("dists")
	Root = paths.New(".build")
	RootApparmord = Root.Join("apparmor.d")
	Distribution = getSupportedDistribution()
}

func Prepare() error {
	for _, fct := range Prepares {
		if err := fct(); err != nil {
			return err
		}
	}
	return nil
}

func Build() error {
	files, _ := RootApparmord.ReadDirRecursiveFiltered(nil, paths.FilterOutDirectories())
	for _, file := range files {
		if !file.Exist() {
			continue
		}
		content, _ := file.ReadFile()
		profile := string(content)
		for _, fct := range Builds {
			profile = fct(profile)
		}
		if err := file.WriteFile([]byte(profile)); err != nil {
			panic(err)
		}
	}
	return nil
}
