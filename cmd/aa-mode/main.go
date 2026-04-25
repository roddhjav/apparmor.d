// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/roddhjav/apparmor.d/pkg/aa"
	"github.com/roddhjav/apparmor.d/pkg/logging"
	"github.com/roddhjav/apparmor.d/pkg/paths"
	"github.com/roddhjav/apparmor.d/pkg/util"
)

const usage = `aa-mode [-h] (-e|-c|-k|-a|-u|-p) [profiles...]

	Switch the given program to an AppArmor mode.

	If a profile name is given without a path, it is looked up in '/etc/apparmor.d/'.
	If a directory is given, all profiles in it are processed recursively.

Options:
    -h, --help             Show this help message and exit.
    -e, --enforce          Set the profile in enforce mode.
    -c, --complain         Set the profile in complain mode.
    -k, --kill             Set the profile in kill mode.
    -a, --default-allow    Set the profile in default_allow mode.
    -u, --unconfined       Set the profile in unconfined mode.
    -p, --prompt           Set the profile in prompt mode.
        --no-reload        Do not reload the profile after modifying it.

`

var (
	help         bool
	enforce      bool
	complain     bool
	kill         bool
	defaultAllow bool
	unconfined   bool
	prompt       bool
	noReload     bool
)

func init() {
	flag.BoolVar(&help, "h", false, "Show this help message and exit.")
	flag.BoolVar(&help, "help", false, "Show this help message and exit.")
	flag.BoolVar(&enforce, "e", false, "Set the profile in enforce mode.")
	flag.BoolVar(&enforce, "enforce", false, "Set the profile in enforce mode.")
	flag.BoolVar(&complain, "c", false, "Set the profile in complain mode.")
	flag.BoolVar(&complain, "complain", false, "Set the profile in complain mode.")
	flag.BoolVar(&kill, "k", false, "Set the profile in kill mode.")
	flag.BoolVar(&kill, "kill", false, "Set the profile in kill mode.")
	flag.BoolVar(&defaultAllow, "a", false, "Set the profile in default_allow mode.")
	flag.BoolVar(&defaultAllow, "default-allow", false, "Set the profile in default_allow mode.")
	flag.BoolVar(&unconfined, "u", false, "Set the profile in unconfined mode.")
	flag.BoolVar(&unconfined, "unconfined", false, "Set the profile in unconfined mode.")
	flag.BoolVar(&prompt, "p", false, "Set the profile in prompt mode.")
	flag.BoolVar(&prompt, "prompt", false, "Set the profile in prompt mode.")
	flag.BoolVar(&noReload, "no-reload", false, "Do not reload the profile after modifying it.")
}

func selectedMode() (string, error) {
	flagsByMode := map[string]bool{
		"enforce":       enforce,
		"complain":      complain,
		"kill":          kill,
		"default_allow": defaultAllow,
		"unconfined":    unconfined,
		"prompt":        prompt,
	}
	var selected string
	for _, mode := range util.ProfileModes {
		if !flagsByMode[mode] {
			continue
		}
		if selected != "" {
			return "", fmt.Errorf("only one mode can be set, got %s and %s", selected, mode)
		}
		selected = mode
	}
	if selected == "" {
		return "", fmt.Errorf("a mode must be set")
	}
	return selected, nil
}

func aaSetMode(files paths.PathList, mode string) error {
	modified := paths.PathList{}
	for _, file := range files {
		profile, err := file.ReadFileAsString()
		if err != nil {
			return err
		}
		if util.IsUnconfined(profile) {
			logging.Warning("skipping %s: profile is in unconfined mode", file)
			continue
		}
		profile, err = util.SetMode(profile, mode)
		if err != nil {
			return err
		}
		if err = file.WriteFile([]byte(profile)); err != nil {
			return err
		}
		modified = append(modified, file)
	}
	if noReload {
		return nil
	}
	return util.ReloadProfiles(modified)
}

func main() {
	flag.Usage = func() { fmt.Print(usage) }
	flag.Parse()
	if help || flag.NArg() < 1 {
		flag.Usage()
		os.Exit(0)
	}

	mode, err := selectedMode()
	if err != nil {
		logging.Fatal("%s", err.Error())
	}
	files, err := paths.PathListFromArgs(flag.Args(), aa.MagicRoot)
	if err != nil {
		logging.Fatal("%s", err.Error())
	}
	if err = aaSetMode(files, mode); err != nil {
		logging.Fatal("%s", err.Error())
	}
}
