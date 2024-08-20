// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2023-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/roddhjav/apparmor.d/pkg/logging"
	"github.com/roddhjav/apparmor.d/pkg/prebuild"
	"github.com/roddhjav/apparmor.d/pkg/prebuild/builder"
	"github.com/roddhjav/apparmor.d/pkg/prebuild/cfg"
	"github.com/roddhjav/apparmor.d/pkg/prebuild/directive"
	"github.com/roddhjav/apparmor.d/pkg/prebuild/prepare"
)

const usage = `prebuild [-h] [--full] [--complain | --enforce]

    Prebuild apparmor.d profiles for a given distribution and apply
    internal built-in directives.

Options:
    -h, --help      Show this help message and exit.
    -f, --full      Set AppArmor for full system policy.
    -c, --complain  Set complain flag on all profiles.
    -e, --enforce   Set enforce flag on all profiles.
        --abi4      Convert the profiles to Apparmor abi/4.0.

`

var (
	help     bool
	full     bool
	complain bool
	enforce  bool
	abi4     bool
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
	flag.BoolVar(&abi4, "abi4", false, "Convert the profiles to Apparmor abi/4.0.")
}

func aaPrebuild() error {
	logging.Step("Building apparmor.d profiles for %s.", cfg.Distribution)

	if full {
		prepare.Register("fsp")
		builder.Register("fsp")
	} else {
		prepare.Register("systemd-early")
	}

	if complain {
		builder.Register("complain")
	} else if enforce {
		builder.Register("enforce")
	}

	if abi4 {
		builder.Register("abi3")
	}

	if err := prebuild.Prepare(); err != nil {
		return err
	}
	return prebuild.Build()
}

func main() {
	flag.Usage = func() {
		fmt.Printf("%s%s\n%s\n%s", usage,
			cfg.Help("Prepare", prepare.Tasks),
			cfg.Help("Build", builder.Builders),
			cfg.Usage("Directives", directive.Directives),
		)
	}
	flag.Parse()
	if help {
		flag.Usage()
		os.Exit(0)
	}
	if err := aaPrebuild(); err != nil {
		logging.Fatal("%s", err.Error())
	}
}
