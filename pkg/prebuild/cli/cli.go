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
)

const (
	nilABI uint = 0
	usage       = `aa-prebuild [-h] [-s] [--complain|--enforce] [--packages] [--full] [--abi 3|4]

    Prebuild apparmor.d profiles for a given distribution, apply
    internal built-in directives and build sub-packages structure.

Options:
    -h, --help      Show this help message and exit.
    -c, --complain  Set complain flag on all profiles.
    -e, --enforce   Set enforce flag on ALL profiles.
    -a, --abi ABI   Target apparmor ABI.
    -f, --full      Set AppArmor for full system policy.
    -p, --packages  Build all split packages.
    -s, --status    Show build configuration.
    -F, --file      Only prebuild a given file.
`
)

var (
	help     bool
	complain bool
	enforce  bool
	full     bool
	packages bool
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
	flag.BoolVar(&packages, "p", false, "Build all split packages.")
	flag.BoolVar(&packages, "packages", false, "Build all split packages.")
	flag.StringVar(&file, "F", "", "Only prebuild a given file.")
	flag.StringVar(&file, "file", "", "Only prebuild a given file.")
}

func Prebuild() {
	flag.Usage = func() {
		fmt.Printf("%s\n%s\n%s\n%s", usage,
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

	if full && paths.New("apparmor.d/groups/_full").Exist() {
		prepare.Register("fsp")
		builder.Register("fsp")
	} else if prebuild.SystemdDir.Exist() {
		prepare.Register("systemd-early")
	}

	if abi != nilABI {
		prebuild.ABI = abi
	}
	switch prebuild.ABI {
	case 3:
		builder.Register("abi3") // Convert all profiles from abi 4.0 to abi 3.0
	case 4:
		// builder.Register("attach") // Re-attach disconnect path
	default:
		logging.Fatal("Invalid ABI version: %d", prebuild.ABI)
	}

	if file != "" {
		sync, _ := prepare.Tasks["synchronise"].(*prepare.Synchronise)
		sync.Path = file
		overwrite, _ := prepare.Tasks["overwrite"].(*prepare.Overwrite)
		overwrite.OneFile = true
	}

	// Prepare the build directories
	logging.Step("Building apparmor.d profiles for %s (abi%d).", prebuild.Distribution, prebuild.ABI)
	prebuild.RootApparmord = prebuild.Root.Join(prebuild.Src)
	if err := Prepare(); err != nil {
		logging.Fatal("%s", err.Error())
	}

	// Generate the packages
	if packages {
		if err := Packages(); err != nil {
			logging.Fatal("%s", err.Error())
		}
	}

	// Build the apparmor.d profiles
	if err := Build(); err != nil {
		logging.Fatal("%s", err.Error())
	}

	if packages {
		// Move all other profiles to apparmor.d.other
		prebuild.RootApparmord = prebuild.Root.Join(prebuild.Src)
		if err := prebuild.RootApparmord.Rename(prebuild.Root.Join("other")); err != nil {
			logging.Fatal("%s", err.Error())
		}
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

func Packages() error {
	logging.Success("Building apparmor.d.* packages structure:")

	for _, name := range prebuild.Packages {
		pkg := prebuild.NewPackage(name)
		msg, err := pkg.Generate()
		if err != nil {
			return err
		}
		if err = pkg.Validate(); err != nil {
			return err
		}
		logging.Indent = "   "
		logging.Bullet("apparmor.d.%s", name)
		logging.Indent += "   "
		for _, line := range util.RemoveDuplicate(msg) {
			logging.Warning("%s", line)
		}
		logging.Indent = ""
	}
	return nil
}

func Build() error {
	sources := []string{prebuild.Src}
	if packages {
		sources = append(sources, prebuild.Packages...)
	}

	for _, src := range sources {
		prebuild.RootApparmord = prebuild.Root.Join(src)
		if src == prebuild.Src {
			setMode("")
		} else {
			pkg := prebuild.NewPackage(src)
			setMode(pkg.Mode)
		}

		files, _ := prebuild.RootApparmord.ReadDirRecursiveFiltered(nil, paths.FilterOutDirectories())
		for _, file := range files {
			if !file.Exist() {
				continue
			}
			profile, err := file.ReadFileAsString()
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

func setMode(mode string) {
	if mode == "" {
		if complain {
			mode = "complain"
		} else if enforce {
			mode = "enforce"
		}
	}
	switch mode {
	case "complain":
		builder.Register("complain")
		builder.Unregister("enforce")
	case "enforce":
		builder.Unregister("complain")
	}
}
