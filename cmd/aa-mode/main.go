// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strings"

	"github.com/roddhjav/apparmor.d/pkg/aa"
	"github.com/roddhjav/apparmor.d/pkg/logging"
	"github.com/roddhjav/apparmor.d/pkg/paths"
	"github.com/roddhjav/apparmor.d/pkg/util"
)

const usage = `aa-mode [-h] (-e|-c|-k|-a|-u|-p|-A|-R) [profiles...]

    Switch the given program to an AppArmor mode.

    If a profile name is given without a path, it is looked up in '/etc/apparmor.d/'.
    If a directory is given, all profiles in it are processed recursively.
    A subprofile is given as 'profile//subprofile'; otherwise the main profile
    and all its subprofiles are processed.

Options:
    -h, --help             Show this help message and exit.
    -e, --enforce          Set the profile in enforce mode.
    -c, --complain         Set the profile in complain mode.
    -k, --kill             Set the profile in kill mode.
    -a, --default-allow    Set the profile in default_allow mode.
    -u, --unconfined       Set the profile in unconfined mode.
    -p, --prompt           Set the profile in prompt mode.
    -A, --all              Insert the 'all' rule in the profile (allow everything).
    -R, --restrict         Remove the 'all' rule from the profile (reverse -A).
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
	all          bool
	restrict     bool
	noReload     bool

	regHeader = regexp.MustCompile(`^([ \t]*)profile[ \t]+(\S+).*\{$`)
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
	flag.BoolVar(&all, "A", false, "Insert the 'all' rule in the profile.")
	flag.BoolVar(&all, "all", false, "Insert the 'all' rule in the profile.")
	flag.BoolVar(&restrict, "R", false, "Remove the 'all' rule from the profile.")
	flag.BoolVar(&restrict, "restrict", false, "Remove the 'all' rule from the profile.")
	flag.BoolVar(&noReload, "no-reload", false, "Do not reload the profile after modifying it.")
}

// target is one profile to modify: the file it lives in and, for a
// subprofile, its name. An empty sub means the main profile and all its
// subprofiles.
type target struct {
	file *paths.Path
	sub  string
}

func selectedMode() (string, error) {
	flagsByMode := map[string]bool{
		"enforce":       enforce,
		"complain":      complain,
		"kill":          kill,
		"default_allow": defaultAllow,
		"unconfined":    unconfined,
		"prompt":        prompt,
		"all":           all,
		"restrict":      restrict,
	}
	var selected string
	for _, mode := range append(slices.Clone(util.ProfileModes), "all", "restrict") {
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

func targetsFromArgs(args []string) ([]target, error) {
	res := []target{}
	for _, arg := range args {
		parts := strings.Split(arg, "//")
		sub := ""
		if len(parts) > 1 {
			sub = parts[len(parts)-1]
		}
		files, err := paths.PathListFromArgs([]string{parts[0]}, aa.MagicRoot)
		if err != nil {
			return nil, err
		}
		for _, file := range files {
			res = append(res, target{file: file, sub: sub})
		}
	}
	return res, nil
}

// setMode applies mode to the profile headers of a profile file: all of them
// when sub is empty, only the named subprofile otherwise. The "all" mode
// inserts the `all,` rule right after the header, "restrict" removes it.
func setMode(profile string, sub string, mode string) (string, error) {
	found := false
	lines := strings.Split(profile, "\n")
	out := make([]string, 0, len(lines)+1)
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		m := regHeader.FindStringSubmatch(line)
		if m == nil || (sub != "" && m[2] != sub) {
			out = append(out, line)
			continue
		}
		found = true
		hasAll := i+1 < len(lines) && strings.TrimSpace(lines[i+1]) == "all,"
		switch mode {
		case "all":
			out = append(out, line)
			if !hasAll {
				out = append(out, m[1]+"  all,")
			}
		case "restrict":
			out = append(out, line)
			if hasAll {
				i++ // drop the all rule
			}
		default:
			header, err := util.SetMode(line+"\n", mode)
			if err != nil {
				return profile, err
			}
			out = append(out, strings.TrimSuffix(header, "\n"))
		}
	}
	if sub != "" && !found {
		return profile, fmt.Errorf("subprofile %s not found", sub)
	}
	return strings.Join(out, "\n"), nil
}

func aaSetMode(targets []target, mode string) error {
	modified := paths.PathList{}
	for _, t := range targets {
		profile, err := t.file.ReadFileAsString()
		if err != nil {
			return err
		}
		if util.IsUnconfined(profile) {
			logging.Warning("skipping %s: profile is in unconfined mode", t.file)
			continue
		}
		profile, err = setMode(profile, t.sub, mode)
		if err != nil {
			return fmt.Errorf("%s: %w", t.file, err)
		}
		if err = t.file.WriteFile([]byte(profile)); err != nil {
			return err
		}
		modified = append(modified, t.file)
	}
	if noReload {
		return nil
	}
	return util.ReloadProfiles(modified)
}

func run() error {
	mode, err := selectedMode()
	if err != nil {
		return err
	}
	targets, err := targetsFromArgs(flag.Args())
	if err != nil {
		return err
	}
	return aaSetMode(targets, mode)
}

func main() {
	flag.Usage = func() { fmt.Print(usage) }
	flag.Parse()
	if help {
		flag.Usage()
		os.Exit(0)
	}
	if err := run(); err != nil {
		logging.Fatal("%s", err.Error())
	}
}
