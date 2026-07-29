// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

// Package detector provides the state variable detectors. Each detector
// lives in its own file and registers itself under the name of the state
// variable it detects.
package detector

import (
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/roddhjav/apparmor.d/pkg/paths"
)

var (
	// Root is the default filesystem root of the system to inspect. Tests
	// override it.
	Root = paths.New("/")

	// Run is the default command runner, returning the command standard
	// output. Tests override it.
	Run = func(name string, arg ...string) (string, error) {
		out, err := exec.Command(name, arg...).Output()
		return string(out), err
	}
)

// System is the handle detectors inspect the running system through.
type System struct {
	// Root is the filesystem root of the system.
	Root *paths.Path

	// Apparmor is the apparmor.d tree being written: the build tree for
	// aa-install, the installed tree for aa-config.
	Apparmor *paths.Path

	// Run executes a command and returns its standard output.
	Run func(name string, arg ...string) (string, error)
}

// NewSystem returns a System for the running host.
func NewSystem(rootApparmor *paths.Path) *System {
	return &System{Root: Root, Apparmor: rootApparmor, Run: Run}
}

// sessions returns the logind session ids.
func (sys *System) sessions() []string {
	out, err := sys.Run("loginctl", "list-sessions", "--no-legend")
	if err != nil {
		return nil
	}
	var ids []string
	for line := range strings.SplitSeq(out, "\n") {
		if fields := strings.Fields(line); len(fields) > 0 {
			ids = append(ids, fields[0])
		}
	}
	return ids
}

// sessionProperty returns the non-empty values of a logind session property
// across all sessions.
func (sys *System) sessionProperty(prop string) []string {
	var values []string
	for _, id := range sys.sessions() {
		out, err := sys.Run("loginctl", "show-session", id, "-p", prop, "--value")
		if err != nil {
			continue
		}
		if value := strings.TrimSpace(out); value != "" {
			values = append(values, value)
		}
	}
	return values
}

// Detector detects the values of one state variable from the enabled system
// state only, never from installed files or distribution specific paths.
type Detector interface {
	// Name is the name of the state variable, without `@{}`.
	Name() string

	// Detect returns the detected values, or nil when they cannot be
	// determined.
	Detect(sys *System) []string
}

// Registry maps a state variable name to its detector. Each detector file
// registers itself on init.
var Registry = map[string]Detector{}

// Register adds a detector to the registry.
func Register(detectors ...Detector) {
	for _, d := range detectors {
		Registry[d.Name()] = d
	}
}

// dedup removes duplicates, preserving order.
func dedup(values []string) []string {
	var res []string
	seen := map[string]bool{}
	for _, value := range values {
		if !seen[value] {
			seen[value] = true
			res = append(res, value)
		}
	}
	return res
}

func sortedGlob(pattern string) []string {
	files, _ := filepath.Glob(pattern)
	sort.Strings(files)
	return files
}
