// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2023-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package main

import (
	"html/template"
	"os/exec"
	"slices"
	"strings"

	"github.com/roddhjav/apparmor.d/pkg/aa"
	"github.com/roddhjav/apparmor.d/pkg/paths"
)

const tmplTest = `#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load common
{{ $name := .Name -}}
{{ range .Commands }}
@test "{{ $name }}: {{ .Description }}" {
    {{ .Cmd }}
}
{{ end }}
`

var (
	Profiles = getProfiles() // List of profiles in apparmor.d
	tmpl     = template.Must(template.New("bats").Parse(tmplTest))
)

type Tests []Test

// Filter returns a new list of tests with only the ones that have a profile
func (t Tests) Filter() Tests {
	for i := len(t) - 1; i >= 0; i-- {
		if !t[i].HasProfile() {
			t = slices.Delete(t, i, i+1)
		}
	}
	return t
}

// Test represents of a list of tests for a given program
type Test struct {
	Name     string
	Commands []Command
}

// Command is a command line to run as part of a test
type Command struct {
	Description string
	Cmd         string
}

// HasProfile returns true if the program in the scenario is profiled in apparmor.d
func (t Test) HasProfile() bool {
	return slices.Contains(Profiles, t.Name)
}

// IsInstalled returns true if the program in the scenario is installed on the system
func (t Test) IsInstalled() bool {
	if _, err := exec.LookPath(t.Name); err != nil {
		return false
	}
	return true
}

func (t Test) Write(dir *paths.Path) error {
	dstDir := dir.Join("profiled")
	if !t.HasProfile() {
		dstDir = dir.Join("unprofiled")
	}
	path := dstDir.Join(t.Name + ".bats")

	if paths.New("tests/bats").Join(t.Name + ".bats").Exist() {
		path = dstDir.Join("00." + t.Name + ".bats")
	}
	content := renderBatsFile(t)
	if err := path.WriteFile([]byte(content)); err != nil {
		return err
	}
	return nil
}

func renderBatsFile(data any) string {
	var res strings.Builder
	err := tmpl.Execute(&res, data)
	if err != nil {
		panic(err)
	}
	return res.String()
}

func getProfiles() []string {
	p := []string{}
	files, err := aa.MagicRoot.ReadDir(paths.FilterOutDirectories())
	if err != nil {
		panic(err)
	}
	for _, path := range files {
		p = append(p, path.Base())
	}
	return p
}
