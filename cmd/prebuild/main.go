// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2023 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/roddhjav/apparmor.d/pkg/logging"
	"github.com/roddhjav/apparmor.d/pkg/prebuild"
)

const usage = `prebuild [-h] [--full] [--complain | --enforce]

    Prebuild apparmor.d profiles for a given distribution.

Options:
    -h, --help      Show this help message and exit.
    -f, --full      Set AppArmor for full system policy.
    -c, --complain  Set complain flag on all profiles.
    -e, --enforce   Set enforce flag on all profiles.
`

var (
	help     bool
	full     bool
	complain bool
	enforce  bool
)

func init() {
	flag.BoolVar(&help, "h", false, "Show this help message and exit.")
	flag.BoolVar(&help, "help", false, "Show this help message and exit.")
	flag.BoolVar(&full, "f", false, "Set AppArmor for full system policy.")
	flag.BoolVar(&full, "full", false, "Set AppArmor for full system policy.")
	flag.BoolVar(&complain, "c", false, "Set complain flag on all profiles.")
	flag.BoolVar(&complain, "complain", false, "Set complain flag on all profiles.")
	flag.BoolVar(&enforce, "e", false, "Set enforce flag on all profiles.")
	flag.BoolVar(&enforce, "enforce", false, "Set enforce flag on all profiles.")
}

func aaPrebuild() error {
	logging.Step("Building apparmor.d profiles for %s.", prebuild.Distribution)

	if full {
		prebuild.Prepares = append(prebuild.Prepares, prebuild.SetFullSystemPolicy)
	}
	if complain {
		prebuild.Builds = append(prebuild.Builds, prebuild.BuildComplain)
	} else if enforce {
		prebuild.Builds = append(prebuild.Builds, prebuild.BuildEnforce)
	}

	if err := prebuild.Prepare(); err != nil {
		return err
	}

	if err := prebuild.Build(); err != nil {
		return err
	}

	logging.Success("Builded profiles with: ")
	logging.Bullet("Bypass userspace tools restriction")
	if complain {
		logging.Bullet("Set complain flag on all profiles")
	} else if enforce {
		logging.Bullet("All profiles have been enforced")
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
	if err := aaPrebuild(); err != nil {
		logging.Fatal(err.Error())
	}
}
