// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
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

const usage = `aa [-h] [--lint | --format | --tree] [-s] [-F file] [profiles...]

    Various AppArmor profiles development tools

Options:
    -h, --help         Show this help message and exit.
    -f, --format       Format the AppArmor profiles.
    -l, --lint         Lint the AppArmor profiles.
    -t, --tree         Generate a tree of visited profiles.
    -F, --file FILE    Set a logfile or a suffix to the default log file.
    -s, --systemd      Parse systemd logs from journalctl.

`

// Command line options
var (
	help    bool
	path    string
	systemd bool
	lint    bool
	format  bool
	tree    bool
)

type kind uint8

const (
	isProfile kind = iota
	isAbstraction
	isTunable
)

func init() {
	flag.BoolVar(&help, "h", false, "Show this help message and exit.")
	flag.BoolVar(&help, "help", false, "Show this help message and exit.")
	flag.BoolVar(&lint, "l", false, "Lint the AppArmor profiles.")
	flag.BoolVar(&lint, "lint", false, "Lint the AppArmor profiles.")
	flag.BoolVar(&format, "f", false, "Format the AppArmor profiles.")
	flag.BoolVar(&format, "format", false, "Format the AppArmor profiles.")
	flag.BoolVar(&tree, "t", false, "Generate a tree of visited profiles.")
	flag.BoolVar(&tree, "tree", false, "Generate a tree of visited profiles.")
	flag.StringVar(&path, "F", "", "Set a logfile or a suffix to the default log file.")
	flag.StringVar(&path, "file", "", "Set a logfile or a suffix to the default log file.")
	flag.BoolVar(&systemd, "s", false, "Parse systemd logs from journalctl.")
	flag.BoolVar(&systemd, "systemd", false, "Parse systemd logs from journalctl.")
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

func parse(kind kind, profile string) (aa.ParaRules, []string, error) {
	var raw string
	paragraphs := []string{}
	rulesByParagraph := aa.ParaRules{}

	switch kind {
	case isTunable, isProfile:
		f := &aa.AppArmorProfileFile{}
		nb, err := f.Parse(profile)
		if err != nil {
			return nil, nil, err
		}
		lines := strings.Split(profile, "\n")
		raw = strings.Join(lines[nb:], "\n")

	case isAbstraction:
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

func formatFile(kind kind, profile string) (string, error) {
	rulesByParagraph, paragraphs, err := parse(kind, profile)
	if err != nil {
		return "", err
	}
	for idx, rules := range rulesByParagraph {
		aa.IndentationLevel = getIndentationLevel(paragraphs[idx])
		rules = rules.Merge().Sort().Format()
		profile = strings.ReplaceAll(profile, paragraphs[idx], rules.String()+"\n")
	}
	return profile, nil
}

// getKind checks if the file is a full apparmor profile file or an
// included (abstraction or tunable) file.
func getKind(file *paths.Path) kind {
	dirname := file.Parent().String()
	switch {
	case strings.Contains(dirname, "abstractions"):
		return isAbstraction
	case strings.Contains(dirname, "tunables"):
		return isTunable
	default:
		return isProfile
	}
}

func aaFormat(files paths.PathList) error {
	for _, file := range files {
		if !file.Exist() {
			return nil
		}
		profile, err := file.ReadFileAsString()
		if err != nil {
			return err
		}

		profile, err = formatFile(getKind(file), profile)
		if err != nil {
			return err
		}
		if err := file.WriteFile([]byte(profile)); err != nil {
			return err
		}
		logging.Success("Formatted: %s", file)
	}
	return nil
}

func aaTree() error {
	return nil
}

func pathsFromArgs() (paths.PathList, error) {
	res := paths.PathList{}
	for _, arg := range flag.Args() {
		path := paths.New(arg)
		switch {
		case !path.Exist():
			return nil, fmt.Errorf("file %s not found", path)
		case path.IsDir():
			files, err := path.ReadDirRecursiveFiltered(nil,
				paths.FilterOutDirectories(),
				paths.FilterOutNames("README.md"),
			)
			if err != nil {
				return nil, err
			}
			res = append(res, files...)
		case path.Exist():
			res = append(res, path)
		}
	}
	return res, nil
}

func main() {
	flag.Usage = func() { fmt.Print(usage) }
	flag.Parse()
	if help {
		flag.Usage()
		os.Exit(0)
	}

	var err error
	var files paths.PathList
	switch {
	case lint:

	case format:
		files, err = pathsFromArgs()
		if err != nil {
			logging.Fatal("%s", err.Error())
		}
		err = aaFormat(files)

	case tree:
		err = aaTree()

	default:
		flag.Usage()
	}

	if err != nil {
		logging.Fatal("%s", err.Error())
	}
}
