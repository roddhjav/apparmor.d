// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2023 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package integration

import (
	"strings"

	"github.com/arduino/go-paths-helper"
	"github.com/roddhjav/apparmor.d/pkg/logs"
	"gopkg.in/yaml.v2"
)

// TestSuite is the apparmod.d integration tests to run
type TestSuite struct {
	Scenarios []Scenario        // List of scenarios to run
	Ignore    []string          // Do not run some scenarios
	Arguments map[string]string // Common arguments used across all scenarios
}

// NewScenarios returns a new list of scenarios
func NewTestSuite() *TestSuite {
	return &TestSuite{
		Scenarios: []Scenario{},
		Ignore:    []string{},
		Arguments: map[string]string{},
	}
}

// Write export the list of scenarios to a file
func (t *TestSuite) Write(path *paths.Path) error {
	jsonString, err := yaml.Marshal(&t.Scenarios)
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
	res = strings.Replace(res, "- name:", "\n- name:", -1)
	_, err = file.WriteString("---\n" + res)
	return err
}

// ReadScenarios import the scenarios from a file
func (t *TestSuite) ReadScenarios(path *paths.Path) error {
	content, _ := path.ReadFile()
	return yaml.Unmarshal(content, &t.Scenarios)
}

// ReadSettings import the common argument and ignore list from a file
func (t *TestSuite) ReadSettings(path *paths.Path) error {
	type temp struct {
		Arguments map[string]string `yaml:"args"`
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
