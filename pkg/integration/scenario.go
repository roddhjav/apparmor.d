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
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/arduino/go-paths-helper"
	"github.com/roddhjav/apparmor.d/pkg/logging"
)

var (
	Ignore    []string          // Do not run some scenarios
	Arguments map[string]string // Common arguments used across all scenarios
	Profiles  paths.PathList    // List of profiles in apparmor.d
)

// Test represents of a list of tests for a given program
type Test struct {
	Name         string            `yaml:"name"`
	Root         bool              `yaml:"root"`      // Run the test as user or as root
	Dependencies []string          `yaml:"require"`   // Packages required for the tests to run "$(pacman -Qqo Scenario.Name)"
	Arguments    map[string]string `yaml:"arguments"` // Arguments to pass to the program, specific to this scenario
	Commands     []Command         `yaml:"tests"`
}

// Command is a command line to run as part of a test
type Command struct {
	Description string   `yaml:"dsc"`
	Cmd         string   `yaml:"cmd"`
	Stdin       []string `yaml:"stdin"`
}

func NewTest() *Test {
	return &Test{
		Name:         "",
		Root:         false,
		Dependencies: []string{},
		Arguments:    map[string]string{},
		Commands:     []Command{},
	}
}

// HasProfile returns true if the program in the scenario is profiled in apparmor.d
func (t *Test) HasProfile() bool {
	for _, path := range Profiles {
		if t.Name == path.Base() {
			return true
		}
	}
	return false
}

// IsInstalled returns true if the program in the scenario is installed on the system
func (t *Test) IsInstalled() bool {
	if _, err := exec.LookPath(t.Name); err != nil {
		return false
	}
	return true
}

func (t *Test) resolve(in string) string {
	res := in
	for key, value := range t.Arguments {
		res = strings.ReplaceAll(res, "{{ "+key+" }}", value)
	}
	return res
}

// mergeArguments merge the arguments of the scenario with the global arguments
// Test arguments have priority over global arguments
func (t *Test) mergeArguments(args map[string]string) {
	if len(t.Arguments) == 0 {
		t.Arguments = map[string]string{}
	}
	for key, value := range args {
		t.Arguments[key] = value
	}
}

// Run the scenarios tests
func (t *Test) Run(dryRun bool) (ran int, nb int, err error) {
	nb = 0
	if t.HasProfile() && t.IsInstalled() {
		logging.Step("%s", t.Name)
		t.mergeArguments(Arguments)
		for _, test := range t.Commands {
			cmd := t.resolve(test.Cmd)
			if !strings.Contains(cmd, "{{") {
				nb++
				if dryRun {
					logging.Bullet(cmd)
				} else {
					cmdErr := t.run(cmd, strings.Join(test.Stdin, "\n"))
					if cmdErr != nil {
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

func (t *Test) run(cmdline string, in string) error {
	var testErr bytes.Buffer

	// Running the command in a shell ensure it does not run confined under the sudo profile.
	// The shell is run unconfined and therefore the cmdline can be confined without no-new-privs issue.
	sufix := " &" // TODO: we need a goroutine here
	cmd := exec.Command("sh", "-c", cmdline+sufix)
	if t.Root {
		cmd = exec.Command("sudo", "sh", "-c", cmdline+sufix)
	}

	stderr := io.MultiWriter(Stderr, &testErr)
	cmd.Stdin = strings.NewReader(in)
	cmd.Stdout = Stdout
	cmd.Stderr = stderr
	err := cmd.Run()
	if testErr.Len() > 0 {
		return fmt.Errorf("%s", testErr.String())
	}
	return err
}
}
