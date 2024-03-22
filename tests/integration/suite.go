// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2023-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package integration

import (
	"os"

	"github.com/arduino/go-paths-helper"
	"github.com/roddhjav/apparmor.d/pkg/logs"
	"github.com/roddhjav/apparmor.d/pkg/util"
	"gopkg.in/yaml.v2"
)

var (
	// Integration tests standard output
	Stdout *os.File

	// Integration tests standard error output
	Stderr *os.File

	stdoutPath = paths.New("tests/out.log")
	stderrPath = paths.New("tests/err.log")
)

// TestSuite is the apparmod.d integration tests to run
type TestSuite struct {
	Tests     []Test            // List of tests to run
	Ignore    []string          // Do not run some tests
	Arguments map[string]string // Common arguments used across all tests
}

// NewScenarios returns a new list of scenarios
func NewTestSuite() *TestSuite {
	var err error
	Stdout, err = stdoutPath.Create()
	if err != nil {
		panic(err)
	}
	Stderr, err = stderrPath.Create()
	if err != nil {
		panic(err)
	}
	return &TestSuite{
		Tests:     []Test{},
		Ignore:    []string{},
		Arguments: map[string]string{},
	}
}

// Write export the list of scenarios to a file
func (t *TestSuite) Write(path *paths.Path) error {
	jsonString, err := yaml.Marshal(&t.Tests)
	if err != nil {
		return err
	}

	path = path.Clean()
	file, err := path.Create()
	if err != nil {
		return err
	}
	defer file.Close()

	// Cleanup a bit
	res := string(jsonString)
	regClean := util.ToRegexRepl([]string{
		"- name:", "\n- name:",
		`(?m)^.*stdin: \[\].*$`, ``,
		`{{`, `{{ `,
		`}}`, ` }}`,
	})
	res = regClean.Replace(res)
	_, err = file.WriteString("---\n" + res)
	return err
}

// ReadTests import the tests from a file
func (t *TestSuite) ReadTests(path *paths.Path) error {
	content, _ := path.ReadFile()
	return yaml.Unmarshal(content, &t.Tests)
}

// ReadSettings import the common argument and ignore list from a file
func (t *TestSuite) ReadSettings(path *paths.Path) error {
	type temp struct {
		Arguments map[string]string `yaml:"arguments"`
		Ignore    []string          `yaml:"ignore"`
	}
	tmp := temp{}
	content, _ := path.ReadFile()
	if err := yaml.Unmarshal(content, &tmp); err != nil {
		return err
	}
	t.Arguments = tmp.Arguments
	t.Ignore = tmp.Ignore
	return nil
}

// Results returns a sum up of the apparmor logs raised by the scenarios
func (t *TestSuite) Results() string {
	file, _ := logs.GetAuditLogs(logs.LogFiles[0])
	aaLogs := logs.NewApparmorLogs(file, "")
	return aaLogs.String()
}

func (t *TestSuite) GetDependencies() []string {
	res := []string{}
	for _, test := range t.Tests {
		res = append(res, test.Dependencies...)
	}
	return res
}
