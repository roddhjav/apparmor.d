// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2023-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/roddhjav/apparmor.d/pkg/logging"
	"github.com/roddhjav/apparmor.d/pkg/paths"
)

const usage = `aa-test [-h] --bootstrap

    Integration tests manager tool for apparmor.d

Options:
    -h, --help         Show this help message and exit.
    -b, --bootstrap    Download tests using tldr pages and generate Bats tests.

`

var (
	help      bool
	bootstrap bool
)

func init() {
	flag.BoolVar(&help, "h", false, "Show this help message and exit.")
	flag.BoolVar(&help, "help", false, "Show this help message and exit.")
	flag.BoolVar(&bootstrap, "b", false, "Download tests using tldr pages and generate Bats tests.")
	flag.BoolVar(&bootstrap, "bootstrap", false, "Download tests using tldr pages and generate Bats tests.")
}

type Config struct {
	TestsDir  *paths.Path // Default: tests
	TldrDir   *paths.Path // Default: tests/tldr
	TldrFile  *paths.Path // Default: tests/tldr.yml
	TestsFile *paths.Path // Default: tests/tests.yml
	BatsDir   *paths.Path // Default: tests/bats_dirty
}

func NewConfig() *Config {
	testsDir := paths.New("tests")
	cfg := Config{
		TestsDir:  testsDir,
		TldrDir:   testsDir.Join("tldr"),
		TldrFile:  testsDir.Join("tldr.yml"),
		TestsFile: testsDir.Join("tldr.yml"),
		BatsDir:   testsDir.Join("bats_dirty"),
	}
	return &cfg
}

func run() error {
	logging.Step("Bootstraping tests")
	cfg := NewConfig()

	tldr := NewTldr(cfg.TldrDir)
	if err := tldr.Download(); err != nil {
		return err
	}

	tests, err := tldr.Parse()
	if err != nil {
		return err
	}

	if err := cfg.BatsDir.RemoveAll(); err != nil {
		return err
	}
	if err := cfg.BatsDir.MkdirAll(); err != nil {
		return err
	}
	if err := cfg.BatsDir.Join("profiled").MkdirAll(); err != nil {
		return err
	}
	if err := cfg.BatsDir.Join("unprofiled").MkdirAll(); err != nil {
		return err
	}
	for _, test := range tests {
		if err := test.Write(cfg.BatsDir); err != nil {
			return err
		}
	}

	logging.Bullet("Bats tests directory: %s", cfg.BatsDir)
	logging.Bullet("Number of profiles with tests found %d", len(tests))
	logging.Bullet("Number of programs without profile found %d", len(tests))
	return nil
}

func main() {
	flag.Usage = func() { fmt.Print(usage) }
	flag.Parse()
	if help {
		flag.Usage()
		os.Exit(0)
	}

	if !bootstrap {
		flag.Usage()
		os.Exit(1)
	}

	err := run()
	if err != nil {
		logging.Fatal("%s", err.Error())
	}
}
