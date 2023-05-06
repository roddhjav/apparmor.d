// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2023 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

// TODO:
// - Finish templating
// - Provide a large selection of resources: files, disks, http server... for automatic test on them
// - Expand support for interactive program (stdin and Control-D)
// - Properlly log the test result
// - Dbus integration

package integration

import (
	"os"
	"os/exec"
	"strings"

	"github.com/arduino/go-paths-helper"
	"github.com/roddhjav/apparmor.d/pkg/logging"
	"golang.org/x/exp/slices"
)

// Scenario represents of a list of tests for a given program
type (
	Scenario struct {
		Name         string            `yaml:"name"`
		Profiled     bool              `yaml:"profiled"`  // The program is profiled in apparmor.d
		Root         bool              `yaml:"root"`      // Run the test as user or as root
		Dependencies []string          `yaml:"require"`   // Packages required for the tests to run "$(pacman -Qqo Scenario.Name)"
		Arguments    map[string]string `yaml:"arguments"` // Arguments to pass to the program. Sepicific to this scenario
		Tests        []Test            `yaml:"tests"`
	}
	Test struct {
		Description string   `yaml:"dsc"`
		Command     string   `yaml:"cmd"`
		Stdin       []string `yaml:"stdin"`
	}
)

func NewScenario() *Scenario {
	return &Scenario{
		Name:         "",
		Profiled:     false,
		Root:         false,
		Dependencies: []string{},
		Arguments:    map[string]string{},
		Tests:        []Test{},
	}
}

// HasProfile returns true if the program in the scenario is profiled in apparmor.d
func (s *Scenario) hasProfile(profiles paths.PathList) bool {
	for _, path := range profiles {
		if s.Name == path.Base() {
			return true
		}
	}
	return false
}

func (s *Scenario) installed() bool {
	if _, err := exec.LookPath(s.Name); err != nil {
		return false
	}
	return true
}

func (s *Scenario) resolve(in string) string {
	res := in
	for key, value := range s.Arguments {
		res = strings.ReplaceAll(res, "{{"+key+"}}", value)
	}
	return res
}

// mergeArguments merge the arguments of the scenario with the global arguments
// Scenarios arguments have priority over global arguments
func (s *Scenario) mergeArguments(args map[string]string) {
	for key, value := range args {
		s.Arguments[key] = value
	}
}

// Run the scenarios tests
func (s *Scenario) Run(dryRun bool) (ran int, nb int, err error) {
	nb = 0
	if s.Profiled && s.installed() {
		if slices.Contains(Ignore, s.Name) {
			return 0, nb, err
		}
		logging.Step("%s", s.Name)
		s.mergeArguments(Arguments)
		for _, test := range s.Tests {
			cmd := s.resolve(test.Command)
			if !strings.Contains(cmd, "{{") {
				nb++
				if dryRun {
					logging.Bullet(cmd)
				} else {
					cmdErr := s.run(cmd, strings.Join(test.Stdin, "\n"))
					if cmdErr != nil {
						// TODO: log the error
						logging.Error("%v", cmdErr)
					} else {
						logging.Success(cmd)
					}
				}
			}
		}
		return 1, nb, err
	}
	return 0, nb, err
}

func (s *Scenario) run(cmdline string, in string) error {
	// Running the command in a shell ensure it does not run confined under the sudo profile.
	// The shell is run unconfined and therefore the cmdline can be confined without no-new-privs issue.
	sufix := " &" // TODO: we need a goroutine here
	cmd := exec.Command("sh", "-c", cmdline+sufix)
	if s.Root {
		cmd = exec.Command("sudo", "sh", "-c", cmdline+sufix)
	}
	cmd.Stdin = strings.NewReader(in)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
