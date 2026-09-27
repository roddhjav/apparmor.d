// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package builder

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/roddhjav/apparmor.d/pkg/aa"
	"github.com/roddhjav/apparmor.d/pkg/tasks"
	"github.com/roddhjav/apparmor.d/pkg/util"
)

var (
	regAttDefinition = regexp.MustCompile(`(?m)^@\{att\} = .*\n`)
)

type ReAttach struct {
	tasks.BaseTask
}

type Disconnected struct {
	tasks.BaseTask
}

// NewAttach creates a new ReAttach builder.
func NewAttach() *ReAttach {
	return &ReAttach{
		BaseTask: tasks.BaseTask{
			Keyword: "attach",
			Msg:     "Feat: re-attach disconnected path",
		},
	}
}

// Apply will re-attach the disconnected path
//   - Add the attach_disconnected.path flag on all profile with the attach_disconnected flag
//   - Replace the base abstraction by attached/base
//   - Replace the consoles abstraction by attached/consoles
//   - For compatibility, non disconnected profile will have the @{att} variable set to /
//
// An already re-attached profile (e.g. a prebuilt profile processed again by
// aa-install) is re-attached from scratch, so a profile that gained the
// attach_disconnected flag since is re-attached too.
func (b ReAttach) Apply(opt *Option, profile string) (string, error) {
	var insert string

	isInside, err := opt.File.IsInsideDir(b.RootApparmor.Join("abstractions/attached"))
	if err != nil {
		return profile, fmt.Errorf("attach: %v", err)
	}
	if isInside {
		return profile, nil // Do not re-attach twice
	}

	// Undo a previous re-attach pass, the included abstractions are already idempotent
	profile = regAttDefinition.ReplaceAllString(profile, "")
	profile = strings.ReplaceAll(profile, ",attach_disconnected.path=@{att}", "")

	// The header name is authoritative: for namespaced profiles the filename
	// (opt.Name) is e.g. "bwrap" while the header is ":glycin:bwrap".
	name := opt.Name
	if matches := regProfileName.FindStringSubmatch(profile); matches != nil {
		name = matches[1]
	}
	origin := "profile " + name

	if strings.Contains(profile, "attach_disconnected") {
		if opt.Kind == aa.ProfileKind {
			if strings.Contains(name, ":") {
				parts := strings.Split(name, ":")
				if len(parts) != 3 {
					return profile, fmt.Errorf("attach: invalid namespaced profile name: %s", name)
				}
				insert = "@{att} = /att/" + parts[1] + "/\n"
			} else {
				insert = "@{att} = /att/" + name + "/\n"
			}
		}
		replacer := strings.NewReplacer(
			"attach_disconnected", "attach_disconnected,attach_disconnected.path=@{att}",
			"include <abstractions/base-strict>", "include <abstractions/attached/base>",
			"include <abstractions/base>", "include <abstractions/attached/base>",
			"include <abstractions/consoles>", "include <abstractions/attached/consoles>",
			"include <abstractions/nameservice-strict>", "include <abstractions/attached/nameservice-strict>",
		)
		profile = replacer.Replace(profile)

	} else {
		if opt.Kind == aa.ProfileKind {
			insert = "@{att} = \"\"\n"
		}

	}

	return strings.Replace(profile, origin, insert+origin, 1), nil
}

// NewDisconnected creates a new Disconnected builder.
func NewDisconnected() *Disconnected {
	return &Disconnected{
		BaseTask: tasks.BaseTask{
			Keyword: "disconnected",
			Msg:     "Feat: set the attach_disconnected flag on all profiles",
		},
	}
}

// Apply adds the attach_disconnected flag to every profile. It must be
// followed by the ReAttach builder to re-attach the newly disconnected paths.
func (b Disconnected) Apply(opt *Option, profile string) (string, error) {
	if opt.Kind != aa.ProfileKind {
		return profile, nil
	}
	if slices.ContainsFunc(util.GetFlags(profile), func(f string) bool {
		return strings.TrimSpace(f) == "attach_disconnected"
	}) {
		return profile, nil
	}
	return util.ApplyFlags(profile, []string{"attach_disconnected"})
}
