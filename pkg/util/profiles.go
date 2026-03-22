// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2023-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package util

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
)

var (
	regFlags         = regexp.MustCompile(`flags=\(([^)]+)\)`)
	regProfileHeader = regexp.MustCompile(` {\n`)
	profileModes     = []string{
		"enforce", "complain", "kill", "default_allow", "unconfined", "prompt",
	}
)

// GetFlags parses the flags from a profile string.
func GetFlags(profile string) []string {
	matches := regFlags.FindStringSubmatch(profile)
	if len(matches) == 0 {
		return nil
	}
	return strings.Split(matches[1], ",")
}

// SetFlags replaces flags in a profile string. If flags is empty, removes the flags clause.
func SetFlags(profile string, flags []string) string {
	profile = regFlags.ReplaceAllLiteralString(profile, "")
	if len(flags) == 0 {
		// Clean up any extra space left after removing flags
		profile = strings.ReplaceAll(profile, "  {\n", " {\n")
		return profile
	}
	flagsStr := " flags=(" + strings.Join(flags, ",") + ") {\n"
	return regProfileHeader.ReplaceAllLiteralString(profile, flagsStr)
}

// SetMode sets the given mode in the profile string, removing any conflicting mode flags.
func SetMode(profile string, mode string) (string, error) {
	if !slices.Contains(profileModes, mode) {
		return profile, fmt.Errorf("unknown profile mode: %s", mode)
	}

	flags := GetFlags(profile)

	// Remove all conflicting mode flags
	flags = slices.DeleteFunc(flags, func(f string) bool {
		return slices.Contains(profileModes, f)
	})

	// "enforce" is the default (no mode flag needed), otherwise add the mode
	if mode != "enforce" {
		flags = append(flags, mode)
	}

	return SetFlags(profile, flags), nil
}
