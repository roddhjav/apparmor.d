// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2023-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/arduino/go-paths-helper"
	"github.com/roddhjav/apparmor.d/pkg/aa"
	"github.com/roddhjav/apparmor.d/pkg/logging"
	bcfg "github.com/roddhjav/apparmor.d/pkg/prebuild/cfg"
	"github.com/roddhjav/apparmor.d/tests/integration"
)

const usage = `aa-test [-h] [--bootstrap | --run | --list]

    Integration tests manager tool for apparmor.d

Options:
    -h, --help         Show this help message and exit.
    -b, --bootstrap    Bootstrap tests using tldr pages.
    -r, --run          Run a predefined list of tests.
    -l, --list         List the configured tests.
    -f, --file FILE    Set a tests file. Default: tests/tests.yml
	-d, --deps         Install tests dependencies.
	-D, --dryrun       Do not do the action, list it.

`

var (
	help      bool
	bootstrap bool
	run       bool
	list      bool
	deps      bool
	dryRun    bool
	cfg       Config
)

type Config struct {
	TldrDir      *paths.Path    // Default: tests/tldr
	ScenariosDir *paths.Path    // Default: tests
	TldrFile     *paths.Path    // Default: tests/tldr.yml
	TestsFile    *paths.Path    // Default: tests/tests.yml
	SettingsFile *paths.Path    // Default: tests/settings.yml
	Profiles     paths.PathList // List of profiles
}

func NewConfig() Config {
	cfg := Config{
		TldrDir:      paths.New("tests/tldr"),
		ScenariosDir: paths.New("tests/"),
		Profiles:     paths.PathList{},
	}
	cfg.TldrFile = cfg.ScenariosDir.Join("tldr.yml")
	cfg.TestsFile = cfg.ScenariosDir.Join("tests.yml")
	cfg.SettingsFile = cfg.ScenariosDir.Join("settings.yml")
	return cfg
}

func LoadTestSuite() (*integration.TestSuite, error) {
	tSuite := integration.NewTestSuite()
	if err := tSuite.ReadTests(cfg.TestsFile); err != nil {
		return tSuite, err
	}
	if err := tSuite.ReadSettings(cfg.SettingsFile); err != nil {
		return tSuite, err
	}
	return tSuite, nil
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
	flag.BoolVar(&deps, "d", false, "Install tests dependencies.")
	flag.BoolVar(&deps, "deps", false, "Install tests dependencies.")
	flag.BoolVar(&dryRun, "D", false, "Do not do the action, list it.")
	flag.BoolVar(&dryRun, "dryrun", false, "Do not do the action, list it.")
}

func testDownload() error {
	tldr := integration.NewTldr(cfg.TldrDir)
	if err := tldr.Download(); err != nil {
		return err
	}

	tSuite, err := tldr.Parse()
	if err != nil {
		return err
	}

	// Default bootstraped scenarios file
	if err := tSuite.Write(cfg.TldrFile); err != nil {
		return err
	}
	logging.Bullet("Default scenarios saved: %s", cfg.TldrFile)
	logging.Bullet("Number of tests found %d", len(tSuite.Tests))
	return nil
}

func testDeps(dryRun bool) error {
	tSuite, err := LoadTestSuite()
	if err != nil {
		return nil
	}

	deps := tSuite.GetDependencies()
	switch bcfg.Distribution {
	case "arch":
		arg := []string{"pacman", "-Sy", "--noconfirm"}
		arg = append(arg, deps...)
		cmd := exec.Command("sudo", arg...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if dryRun {
			fmt.Println(strings.Join(cmd.Args, " "))
		} else {
			return cmd.Run()
		}
	default:
	}
	return nil
}

func testRun(dryRun bool) error {
	// Warning: There is no guarantee that the tests are not destructive
	if dryRun {
		logging.Step("List tests")
	} else {
		logging.Step("Run tests")
	}

	tSuite, err := LoadTestSuite()
	if err != nil {
		return nil
	}
	integration.Arguments = tSuite.Arguments
	integration.Ignore = tSuite.Ignore
	integration.Profiles = cfg.Profiles
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
	} else if deps {
		err = testDeps(dryRun)
	} else {
		flag.Usage()
		os.Exit(1)
	}
	if err != nil {
		logging.Fatal(err.Error())
	}
}
