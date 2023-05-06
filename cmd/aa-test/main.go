// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2023 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/arduino/go-paths-helper"
	intg "github.com/roddhjav/apparmor.d/pkg/integration"
	"github.com/roddhjav/apparmor.d/pkg/logging"
)

const usage = `aa-test [-h] [--bootstrap | --run | --list]

    Integration tests manager tool for apparmor.d

Options:
    -h, --help         Show this help message and exit.
    -b, --bootstrap    Bootstrap test scenarios using tldr pages.
    -r, --run          Run a predefined list of tests.
    -l, --list         List the configured tests.
    -f, --file FILE    Set a scenario test file. Default: tests/scenarios.yml

`

const defaultScenariosFile = "scenarios.yml"

var (
	help      bool
	bootstrap bool
	run       bool
	list      bool
	path      string
	Cfg       Config
)

type Config struct {
	TldrDir       *paths.Path // Default: tests/tldr
	ScenariosDir  *paths.Path // Default: tests
	ScenariosFile *paths.Path // Default: tests/tldr.yml
	ProfilesDir   *paths.Path // Default: /etc/apparmor.d
	Profiles      paths.PathList
}

func NewConfig() Config {
	cfg := Config{
		TldrDir:      paths.New("tests/tldr"),
		ScenariosDir: paths.New("tests/"),
		ProfilesDir:  paths.New("/etc/apparmor.d"),
		Profiles:     paths.PathList{},
	}
	cfg.ScenariosFile = cfg.ScenariosDir.Join("tldr.yml")
	return cfg
}

func init() {
	Cfg = NewConfig()
	files, _ := Cfg.ProfilesDir.ReadDir(paths.FilterOutDirectories())
	for _, path := range files {
		Cfg.Profiles.Add(path)
	}

	flag.BoolVar(&help, "h", false, "Show this help message and exit.")
	flag.BoolVar(&help, "help", false, "Show this help message and exit.")
	flag.BoolVar(&bootstrap, "b", false, "Bootstrap test scenarios using tldr pages.")
	flag.BoolVar(&bootstrap, "bootstrap", false, "Bootstrap test scenarios using tldr pages.")
	flag.BoolVar(&run, "r", false, "Run a predefined list of tests.")
	flag.BoolVar(&run, "run", false, "Run a predefined list of tests.")
	flag.BoolVar(&list, "l", false, "List the test to run.")
	flag.BoolVar(&list, "list", false, "List the test to run.")
	flag.StringVar(&path, "f", defaultScenariosFile, "Set a scenario test file.")
	flag.StringVar(&path, "file", defaultScenariosFile, "Set a scenario test file.")
}

func apparmorTestBootstrap() error {
	tldr := intg.NewTldr(Cfg.TldrDir)
	if err := tldr.Download(); err != nil {
		return err
	}

	tSuite, err := tldr.Parse(Cfg.Profiles)
	if err != nil {
		return err
	}

	// Default bootstraped scenarios file
	if err := tSuite.Write(Cfg.ScenariosFile); err != nil {
		return err
	}
	logging.Bullet("Default scenarios saved: %s", Cfg.ScenariosFile)

	// Scenarios file with only profiled programs
	tSuiteWithProfile := intg.NewTestSuite()
	for _, s := range tSuite.Scenarios {
		if s.Profiled {
			tSuiteWithProfile.Scenarios = append(tSuiteWithProfile.Scenarios, s)
		}
	}
	scnPath := Cfg.ScenariosDir.Join(strings.TrimSuffix(path, filepath.Ext(path)) + ".new.yml")
	if err := tSuiteWithProfile.Write(scnPath); err != nil {
		return err
	}
	logging.Bullet("Scenarios with profile saved: %s", scnPath)

	logging.Bullet("Number of scenarios found %d", len(tSuite.Scenarios))
	logging.Bullet("Number of scenarios with profiles in apparmor.d %d", len(tSuiteWithProfile.Scenarios))
	return nil
}

func apparmorTestRun(dryRun bool) error {
	// FIXME: Safety settings. For default scenario set with sudo enabled the apparmor.d folder gets removed.
	dryRun = true
	if dryRun {
		logging.Step("List tests")
	} else {
		logging.Step("Run tests")
	}

	tSuite := intg.NewTestSuite()
	scnPath := Cfg.ScenariosDir.Join(path)
	if err := tSuite.ReadScenarios(scnPath); err != nil {
		return err
	}
	cfgPath := Cfg.ScenariosDir.Join("integration.yml")
	if err := tSuite.ReadSettings(cfgPath); err != nil {
		return err
	}
	intg.Arguments = tSuite.Arguments
	intg.Ignore = tSuite.Ignore

	nbTest := 0
	nbScn := 0
	for _, scn := range tSuite.Scenarios {
		ran, nb, err := scn.Run(dryRun)
		nbScn += ran
		nbTest += nb
		if err != nil {
			return err
		}
	}

	if dryRun {
		logging.Bullet("Number of scenarios to run %d", nbScn)
		logging.Bullet("Number of tests to run %d", nbTest)
	} else {
		logging.Success("Number of scenarios ran %d", nbScn)
		logging.Success("Number of tests to ran %d", nbTest)
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

	var err error
	if bootstrap {
		logging.Step("Bootstraping tests")
		err = apparmorTestBootstrap()
	} else if run || list {
		err = apparmorTestRun(list)
	} else {
		flag.Usage()
		os.Exit(1)
	}
	if err != nil {
		logging.Fatal(err.Error())
	}
}
