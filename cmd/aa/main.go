// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/roddhjav/apparmor.d/pkg/aa"
	"github.com/roddhjav/apparmor.d/pkg/logging"
	"github.com/roddhjav/apparmor.d/pkg/paths"
)

const usage = `aa [-h] <command> [profiles...]

    Various unprivileged AppArmor profiles tools.

Command:
    lint        Lint the AppArmor profiles.
    format      Format the AppArmor profiles.
    tree        Generate a tree of visited profiles.

Options:
    -h, --help  Show this help message and exit.

`

// Command line options
var (
	help    bool
	command string
)

func init() {
	flag.BoolVar(&help, "h", false, "Show this help message and exit.")
	flag.BoolVar(&help, "help", false, "Show this help message and exit.")
}

func getIndentationLevel(input string) int {
	level := 0
	after, found := strings.CutPrefix(input, aa.Indentation)
	for found {
		after, found = strings.CutPrefix(after, aa.Indentation)
		level++
	}
	if strings.Contains(input, "owner") && level > 3 {
		level = level - 3
	}
	return level
}

func parse(kind aa.FileKind, profile string) (aa.ParaRules, []string, error) {
	var raw string
	paragraphs := []string{}
	rulesByParagraph := aa.ParaRules{}

	switch kind {
	case aa.TunableKind, aa.ProfileKind:
		f := &aa.AppArmorProfileFile{}
		nb, err := f.Parse(profile)
		if err != nil {
			return nil, nil, err
		}
		lines := strings.Split(profile, "\n")
		raw = strings.Join(lines[nb:], "\n")

	case aa.AbstractionKind:
		raw = profile
	}

	r, par, err := aa.ParseRules(raw)
	if err != nil {
		return nil, nil, err
	}
	rulesByParagraph = append(rulesByParagraph, r...)
	paragraphs = append(paragraphs, par...)
	return rulesByParagraph, paragraphs, nil
}

func formatFile(kind aa.FileKind, profile string) (string, error) {
	rulesByParagraph, paragraphs, err := parse(kind, profile)
	if err != nil {
		return "", err
	}
	for idx, rules := range rulesByParagraph {
		aa.IndentationLevel = getIndentationLevel(paragraphs[idx])
		rules = rules.Merge().Sort().Format()
		fmt.Print(rules.String() + "\n")
	}
	return profile, nil
}

func aaFormat(files paths.PathList) error {
	for _, file := range files {
		if !file.Exist() {
			return nil
		}

		raw, err := file.ReadFileAsString()
		if err != nil {
			return err
		}

		raw, err = formatFile(aa.KindFromPath(file), raw)
		if err != nil {
			return err
		}
		if err := file.WriteFile([]byte(raw)); err != nil {
			return err
		}
		logging.Success("Formatted: %s", file)
	}
	return nil
}

func aaLint(files paths.PathList) error {
	for _, file := range files {
		fmt.Printf("wip: %v\n", file)
	}
	return nil
}

func aaTree() error {
	return nil
}

func main() {
	flag.Usage = func() { fmt.Print(usage) }
	flag.Parse()
	if help {
		flag.Usage()
		os.Exit(0)
	}
	if len(flag.Args()) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	command = flag.Args()[0]
	if err := flag.CommandLine.Parse(flag.Args()[1:]); err != nil {
		logging.Fatal("%s", err.Error())
	}

	var err error
	var files paths.PathList
	switch command {
	case "lint":
		files, err = paths.PathListFromArgs(flag.Args(), aa.MagicRoot)
		if err != nil {
			logging.Fatal("%s", err.Error())
		}
		err = aaLint(files)

	case "format":
		files, err = paths.PathListFromArgs(flag.Args(), aa.MagicRoot)
		if err != nil {
			logging.Fatal("%s", err.Error())
		}
		err = aaFormat(files)

	case "tree":
		err = aaTree()

	default:
		flag.Usage()
		os.Exit(1)
	}

	if err != nil {
		logging.Fatal("%s", err.Error())
	}
}
