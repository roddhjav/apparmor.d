// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2023 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/arduino/go-paths-helper"
	"github.com/roddhjav/apparmor.d/pkg/logging"
	"github.com/roddhjav/apparmor.d/pkg/util"
)

const usage = `prebuild [-h] [--full] [--complain]

    Internal tool to prebuild apparmor.d profiles for a given distribution.

Options:
    -h, --help      Show this help message and exit.
    -d, --dist      The target Linux distribution.
    -f, --full      Set AppArmor for full system policy.
    -c, --complain  Set complain flag on all profiles.
`

var (
	help          bool
	Full          bool
	Complain      bool
	Distribution  string
	DistDir       *paths.Path
	Root          *paths.Path
	RootApparmord *paths.Path

	// Prepare the build directory with the following tasks
	prepare = []prepareFunc{Synchronise, Ignore, Merge, Configure, SetFlags, SetFullSystemPolicy}

	// Build the profiles with the following build tasks
	build = []buildFunc{BuildUserspace, BuildComplain, BuildABI}
)

type prepareFunc func() error
type buildFunc func(string) string

func init() {
	DistDir = paths.New("dists")
	Root = paths.New(".build")
	RootApparmord = Root.Join("apparmor.d")
	Distribution, _ = util.GetSupportedDistribution()
	flag.BoolVar(&help, "h", false, "Show this help message and exit.")
	flag.BoolVar(&help, "help", false, "Show this help message and exit.")
	flag.BoolVar(&Full, "f", false, "Set AppArmor for full system policy.")
	flag.BoolVar(&Full, "full", false, "Set AppArmor for full system policy.")
	flag.BoolVar(&Complain, "c", false, "Set complain flag on all profiles.")
	flag.BoolVar(&Complain, "complain", false, "Set complain flag on all profiles.")
}

// Build the profiles.
func buildProfiles() error {
	files, _ := RootApparmord.ReadDir(paths.FilterOutDirectories())
	for _, file := range files {
		if !file.Exist() {
			continue
		}
		content, _ := file.ReadFile()
		profile := string(content)
		for _, fct := range build {
			profile = fct(profile)
		}
		if err := file.WriteFile([]byte(profile)); err != nil {
			panic(err)
		}
	}
	return nil
}

func aaPrebuild() error {
	logging.Step("Building apparmor.d profiles for %s.", Distribution)

	for _, fct := range prepare {
		if err := fct(); err != nil {
			return err
		}
	}

	if err := buildProfiles(); err != nil {
		return err
	}
	logging.Success("Builded profiles with: ")
	logging.Bullet("Bypass userspace tools restriction")
	if Complain {
		logging.Bullet("Set complain flag on all profiles")
	}
	switch Distribution {
	case "debian", "whonix":
		logging.Bullet("%s does not support abi 3.0 yet", Distribution)
	}
	return nil
}

func main() {
	flag.Usage = func() { fmt.Print(usage) }
	flag.Parse()
	if help {
		flag.Usage()
		os.Exit(0)
	}
	err := aaPrebuild()
	if err != nil {
		logging.Fatal(err.Error())
	}
}
