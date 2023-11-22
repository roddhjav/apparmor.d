// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2023 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package prebuild

import (
	"regexp"
	"strings"

	"github.com/roddhjav/apparmor.d/pkg/aa"
	"github.com/roddhjav/apparmor.d/pkg/util"
	"golang.org/x/exp/slices"
)

// Build the profiles with the following build tasks
var Builds = []BuildFunc{
	BuildUserspace,
}

var (
	regAttachments      = regexp.MustCompile(`(profile .* @{exec_path})`)
	regFlags            = regexp.MustCompile(`flags=\(([^)]+)\)`)
	regProfileHeader    = regexp.MustCompile(` {`)
	regFullSystemPolicy = util.ToRegexRepl([]string{
		`r(PU|U)x,`, `rPx,`,
	})
	regAbi4To3 = util.ToRegexRepl([]string{ // Currently Abi3 -> Abi4
		`abi/3.0`, `abi/4.0`,
		`# userns,`, `userns,`,
	})
)

type BuildFunc func(string) string

// Set complain flag on all profiles
func BuildComplain(profile string) string {
	flags := []string{}
	matches := regFlags.FindStringSubmatch(profile)
	if len(matches) != 0 {
		flags = strings.Split(matches[1], ",")
		if slices.Contains(flags, "complain") {
			return profile
		}
	}
	flags = append(flags, "complain")
	strFlags := " flags=(" + strings.Join(flags, ",") + ") {"

	// Remove all flags definition, then set manifest' flags
	profile = regFlags.ReplaceAllLiteralString(profile, "")
	return regProfileHeader.ReplaceAllLiteralString(profile, strFlags)
}

// Set all profiles in enforce mode
func BuildEnforce(profile string) string {
	matches := regFlags.FindStringSubmatch(profile)
	if len(matches) == 0 {
		return profile
	}

	flags := strings.Split(matches[1], ",")
	idx := slices.Index(flags, "complain")
	if idx == -1 {
		return profile
	}
	flags = slices.Delete(flags, idx, idx+1)
	strFlags := "{"
	if len(flags) >= 1 {
		strFlags = " flags=(" + strings.Join(flags, ",") + ") {"
	}

	// Remove all flags definition, then set new flags
	profile = regFlags.ReplaceAllLiteralString(profile, "")
	return regProfileHeader.ReplaceAllLiteralString(profile, strFlags)
}

// Bypass userspace tools restriction
func BuildUserspace(profile string) string {
	p := aa.DefaultTunables()
	p.ParseVariables(profile)
	p.ResolveAttachments()
	att := p.NestAttachments()
	matches := regAttachments.FindAllString(profile, -1)
	if len(matches) > 0 {
		strheader := strings.Replace(matches[0], "@{exec_path}", att, -1)
		return regAttachments.ReplaceAllLiteralString(profile, strheader)
	}
	return profile
}

func BuildABI3(profile string) string {
	for _, abi4t3 := range regAbi4To3 {
		profile = abi4t3.Regex.ReplaceAllLiteralString(profile, abi4t3.Repl)
	}
	return profile
}

func BuildFullSystemPolicy(profile string) string {
	for _, full := range regFullSystemPolicy {
		profile = full.Regex.ReplaceAllString(profile, full.Repl)
	}
	return profile
}
