// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2023-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

// Package util gathers small, dependency-free helpers shared across the
// codebase: profile flag manipulation, and AppArmor userspace utilities
// such as profile reload.
package util

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
)

var (
	regFlags         = regexp.MustCompile(`flags=\(([^)]+)\)`)
	regProfileHeader = regexp.MustCompile(`(?m)^([ \t]*profile [^\n]*?) \{\n`)
	ProfileModes     = []string{
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
	profile = strings.ReplaceAll(profile, "  {\n", " {\n")
	if len(flags) == 0 {
		return profile
	}
	flagsStr := "${1} flags=(" + strings.Join(flags, ",") + ") {\n"
	return regProfileHeader.ReplaceAllString(profile, flagsStr)
}

// IsUnconfined reports whether any profile in the given content has the unconfined mode flag set.
func IsUnconfined(profile string) bool {
	for _, match := range regFlags.FindAllStringSubmatch(profile, -1) {
		for f := range strings.SplitSeq(match[1], ",") {
			if strings.TrimSpace(f) == "unconfined" {
				return true
			}
		}
	}
	return false
}

// ApplyFlags merges flags into a profile, mode-aware: non-mode flags are
// added to the ones the profile already has (e.g. attach_disconnected is
// preserved), and a mode flag (enforce, complain, ...) is applied via SetMode
// so modes never stack into an invalid combo such as (enforce,complain). The
// last mode listed wins.
func ApplyFlags(profile string, flags []string) (string, error) {
	current := GetFlags(profile)
	var mode string
	for _, f := range flags {
		if slices.Contains(ProfileModes, f) {
			mode = f
			continue
		}
		if !slices.Contains(current, f) {
			current = append(current, f)
		}
	}
	profile = SetFlags(profile, current)
	if mode == "" {
		return profile, nil
	}
	return SetMode(profile, mode)
}

// SetMode sets the given mode in the profile string, removing any conflicting mode flags.
func SetMode(profile string, mode string) (string, error) {
	if !slices.Contains(ProfileModes, mode) {
		return profile, fmt.Errorf("unknown profile mode: %s", mode)
	}

	flags := GetFlags(profile)

	// An explicitly unconfined profile stays unconfined: never replace it.
	// unconfined is never a target mode, but it can already exist in a profile.
	if slices.Contains(flags, "unconfined") {
		return profile, nil
	}

	// Remove all conflicting mode flags
	flags = slices.DeleteFunc(flags, func(f string) bool {
		return slices.Contains(ProfileModes, f)
	})

	// "enforce" is the default (no mode flag needed), otherwise add the mode
	if mode != "enforce" {
		flags = append(flags, mode)
	}

	return SetFlags(profile, flags), nil
}
