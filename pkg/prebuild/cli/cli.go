// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2023-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package cli

import (
	"flag"
	"fmt"
	"strings"

	"github.com/roddhjav/apparmor.d/pkg/logging"
	"github.com/roddhjav/apparmor.d/pkg/paths"
	"github.com/roddhjav/apparmor.d/pkg/prebuild"
	"github.com/roddhjav/apparmor.d/pkg/prebuild/builder"
	"github.com/roddhjav/apparmor.d/pkg/prebuild/directive"
	"github.com/roddhjav/apparmor.d/pkg/prebuild/prepare"
	"github.com/roddhjav/apparmor.d/pkg/util"
)

const (
	nilABI uint = 0
	usage       = `aa-prebuild [-h] [--complain | --enforce] [--full] [--abi 3|4]

    Prebuild apparmor.d profiles for a given distribution and apply
    internal built-in directives.

Options:
    -h, --help      Show this help message and exit.
    -c, --complain  Set complain flag on all profiles.
    -e, --enforce   Set enforce flag on all profiles.
    -a, --abi ABI   Target apparmor ABI.
    -f, --full      Set AppArmor for full system policy.
    -F, --file      Only prebuild a given file.
`
)

var (
	help     bool
	complain bool
	enforce  bool
	full     bool
	abi      uint
	file     string
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
	flag.UintVar(&abi, "a", nilABI, "Target apparmor ABI.")
	flag.UintVar(&abi, "abi", nilABI, "Target apparmor ABI.")
	flag.StringVar(&file, "F", "", "Only prebuild a given file.")
	flag.StringVar(&file, "file", "", "Only prebuild a given file.")
}

func Prebuild() {
	flag.Usage = func() {
		fmt.Printf("%s%s\n%s\n%s", usage,
			prebuild.Help("Prepare", prepare.Tasks),
			prebuild.Help("Build", builder.Builders),
			directive.Usage(),
		)
	}
	flag.Parse()
	if help {
		flag.Usage()
		return
	}

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

	switch abi {
	case 3:
		prebuild.ABI = 3
		builder.Register("abi3")
	case 4:
		prebuild.ABI = 4
		for i, b := range builder.Builds {
			if b.Name() == "abi3" {
				builder.Builds = append(builder.Builds[:i], builder.Builds[i+1:]...)
				break
			}
		}
	case nilABI:
	default:
		logging.Fatal("ABI %d not supported", abi)
	}

	if file != "" {
		sync, _ := prepare.Tasks["synchronise"].(*prepare.Synchronise)
		sync.Path = file
		overwrite, _ := prepare.Tasks["overwrite"].(*prepare.Overwrite)
		overwrite.OneFile = true
	}

	logging.Step("Building apparmor.d profiles for %s on ABI%d.", prebuild.Distribution, prebuild.ABI)
	if err := Prepare(); err != nil {
		logging.Fatal("%s", err.Error())
	}
	if err := Build(); err != nil {
		logging.Fatal("%s", err.Error())
	}
}

func Prepare() error {
	for _, task := range prepare.Prepares {
		msg, err := task.Apply()
		if err != nil {
			return err
		}
		if file != "" && task.Name() == "setflags" {
			continue
		}
		logging.Success("%s", task.Message())
		logging.Indent = "   "
		for _, line := range msg {
			if strings.Contains(line, "not found") {
				logging.Warning("%s", line)
			} else {
				logging.Bullet("%s", line)
			}
		}
		logging.Indent = ""
	}
	return nil
}

func Build() error {
	files, _ := prebuild.RootApparmord.ReadDirRecursiveFiltered(nil, paths.FilterOutDirectories())
	for _, file := range files {
		if !file.Exist() {
			continue
		}
		profile, err := util.ReadFile(file)
		if err != nil {
			return err
		}
		profile, err = builder.Run(file, profile)
		if err != nil {
			return err
		}
		profile, err = directive.Run(file, profile)
		if err != nil {
			return err
		}
		if err := file.WriteFile([]byte(profile)); err != nil {
			return err
		}
	}

	logging.Success("Build tasks:")
	logging.Indent = "   "
	for _, task := range builder.Builds {
		logging.Bullet("%s", task.Message())
	}
	logging.Indent = ""
	logging.Success("Directives processed:")
	logging.Indent = "   "
	for _, dir := range directive.Directives {
		logging.Bullet("%s%s", directive.Keyword, dir.Name())
	}
	logging.Indent = ""
	return nil
}
