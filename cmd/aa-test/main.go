// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2023 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/arduino/go-paths-helper"
	"github.com/roddhjav/apparmor.d/pkg/aa"
	"github.com/roddhjav/apparmor.d/pkg/integration"
	"github.com/roddhjav/apparmor.d/pkg/logging"
)

const usage = `aa-test [-h] [--bootstrap | --run | --list]

    Integration tests manager tool for apparmor.d

Options:
    -h, --help         Show this help message and exit.
    -b, --bootstrap    Bootstrap tests using tldr pages.
    -r, --run          Run a predefined list of tests.
    -l, --list         List the configured tests.
    -f, --file FILE    Set a tests file. Default: tests/tests.yml

`

var (
	help      bool
	bootstrap bool
	run       bool
	list      bool
	cfg       Config
)

type Config struct {
	TldrDir      *paths.Path // Default: tests/tldr
	ScenariosDir *paths.Path // Default: tests
	TldrFile     *paths.Path // Default: tests/tldr.yml
	TestsFile    *paths.Path // Default: tests/tests.yml
	Profiles     paths.PathList
}

func NewConfig() Config {
	cfg := Config{
		TldrDir:      paths.New("tests/tldr"),
		ScenariosDir: paths.New("tests/"),
		Profiles:     paths.PathList{},
	}
	cfg.TldrFile = cfg.ScenariosDir.Join("tldr.yml")
	cfg.TestsFile = cfg.ScenariosDir.Join("tests.yml")
	return cfg
}

func init() {
	cfg = NewConfig()
	files, _ := aa.MagicRoot.ReadDir(paths.FilterOutDirectories())
	for _, path := range files {
		cfg.Profiles.Add(path)
	}

	flag.BoolVar(&help, "h", false, "Show this help message and exit.")
	flag.BoolVar(&help, "help", false, "Show this help message and exit.")
	flag.BoolVar(&bootstrap, "b", false, "Bootstrap tests using tldr pages.")
	flag.BoolVar(&bootstrap, "bootstrap", false, "Bootstrap tests using tldr pages.")
	flag.BoolVar(&run, "r", false, "Run a predefined list of tests.")
	flag.BoolVar(&run, "run", false, "Run a predefined list of tests.")
	flag.BoolVar(&list, "l", false, "List the tests to run.")
	flag.BoolVar(&list, "list", false, "List the tests to run.")
}

func testDownload() error {
	tldr := integration.NewTldr(cfg.TldrDir)
	if err := tldr.Download(); err != nil {
		return err
	}

	tSuite, err := tldr.Parse(cfg.Profiles)
	if err != nil {
		return err
	}

	// Default bootstraped scenarios file
	if err := tSuite.Write(cfg.TldrFile); err != nil {
		return err
	}
	logging.Bullet("Default scenarios saved: %s", cfg.TldrFile)

	// Scenarios file with only profiled programs
	tSuiteWithProfile := integration.NewTestSuite()
	for _, t := range tSuite.Tests {
		if t.Profiled {
			tSuiteWithProfile.Tests = append(tSuiteWithProfile.Tests, t)
		}
	}

	testWithProfilePath := cfg.TldrFile.Parent().Join(
		strings.TrimSuffix(cfg.TldrFile.Base(), cfg.TldrFile.Ext()) + ".new.yml")
	if err := tSuiteWithProfile.Write(testWithProfilePath); err != nil {
		return err
	}
	logging.Bullet("Tests with profile saved: %s", testWithProfilePath)

	logging.Bullet("Number of tests found %d", len(tSuite.Tests))
	logging.Bullet("Number of tests with profiles in apparmor.d %d", len(tSuiteWithProfile.Tests))
	return nil
}

func testRun(dryRun bool) error {
	// Warning: There is no guarantee that the tests are not destructive
	if dryRun {
		logging.Step("List tests")
	} else {
		logging.Step("Run tests")
	}

	tSuite := integration.NewTestSuite()
	if err := tSuite.ReadScenarios(cfg.TestsFile); err != nil {
		return err
	}
	cfgPath := cfg.ScenariosDir.Join("integration.yml")
	if err := tSuite.ReadSettings(cfgPath); err != nil {
		return err
	}
	integration.Arguments = tSuite.Arguments
	integration.Ignore = tSuite.Ignore
	nbCmd := 0
	nbTest := 0
	for _, test := range tSuite.Tests {
		ran, nb, err := test.Run(dryRun)
		nbTest += ran
		nbCmd += nb
		if err != nil {
			return err
		}
	}

	if dryRun {
		logging.Bullet("Number of tests to run %d", nbTest)
		logging.Bullet("Number of test commands to run %d", nbCmd)
	} else {
		logging.Success("Number of tests ran %d", nbTest)
		logging.Success("Number of test command to ran %d", nbCmd)
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
		err = testDownload()
	} else if run || list {
		err = testRun(list)
	} else {
		flag.Usage()
		os.Exit(1)
	}
	if err != nil {
		logging.Fatal(err.Error())
	}
}
