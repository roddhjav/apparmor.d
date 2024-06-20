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
	"github.com/roddhjav/apparmor.d/pkg/util"
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

func parse(profile string) (*aa.AppArmorProfileFile, []aa.Rules, []string, error) {
	f := &aa.AppArmorProfileFile{}
	nb, err := f.Parse(profile)
	if err != nil {
		return nil, nil, nil, err
	}
	paragraphRules, paragraphs, err := aa.ParseRules(strings.Join(strings.Split(profile, "\n")[nb:], "\n"))
	if err != nil {
		return nil, nil, nil, err
	}
	return f, paragraphRules, paragraphs, nil
}

func formatFile(profile string) (string, error) {
	_, paragraphRules, paragraphs, err := parse(profile)
	if err != nil {
		return "", err
	}
	for idx, rules := range paragraphRules {
		if err := rules.Validate(); err != nil {
			return "", err
		}
		aa.IndentationLevel = getIndentationLevel(paragraphs[idx])
		rules = rules.Merge().Sort().Format()
		profile = strings.Replace(profile, paragraphs[idx], rules.String()+"\n", -1)
	}
	return profile, nil
}

func aaFormat(files paths.PathList) error {
	for _, file := range files {
		if !file.Exist() {
			return nil
		}
		profile, err := util.ReadFile(file)
		if err != nil {
			return err
		}
		profile, err = formatFile(profile)
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
			logging.Fatal(err.Error())
		}
		err = aaFormat(files)
	case tree:
		err = aaTree()
	}

	if err != nil {
		logging.Fatal(err.Error())
	}
}
