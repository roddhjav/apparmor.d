// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package builder

import (
	"strings"

	"github.com/roddhjav/apparmor.d/pkg/aa"
	"github.com/roddhjav/apparmor.d/pkg/prebuild"
)

type ReAttach struct {
	prebuild.Base
}

func init() {
	RegisterBuilder(&ReAttach{
		Base: prebuild.Base{
			Keyword: "attach",
			Msg:     "Feat: re-attach disconnected path",
		},
	})
}

// Apply will re-attach the disconnected path
//   - Add the attach_disconnected.path flag on all frofile with the attach_disconnected flag
//   - Replace the base abstraction by attached/base
//   - Replace the consoles abstraction by attached/consoles
//   - For compatibility, non disconnected profile will have the @{att} variable set to /
func (b ReAttach) Apply(opt *Option, profile string) (string, error) {
	var insert string
	var origin = "profile " + opt.Name
	if opt.File.HasSuffix("attached/base") {
		return profile, nil // Do not re-attach twice
	}

	if strings.Contains(profile, "attach_disconnected") {
		if opt.Kind == aa.ProfileKind {
			insert = "@{att} = /att/" + opt.Name + "/\n"
		}
		profile = strings.ReplaceAll(profile,
			"attach_disconnected",
			"attach_disconnected,attach_disconnected.path=@{att}",
		)
		profile = strings.ReplaceAll(profile,
			"include <abstractions/base>",
			"include <abstractions/attached/base>",
		)
		profile = strings.ReplaceAll(profile,
			"include <abstractions/base-strict>",
			"include <abstractions/attached/base>",
		)
		profile = strings.ReplaceAll(profile,
			"include <abstractions/consoles>",
			"include <abstractions/attached/consoles>",
		)
		profile = strings.ReplaceAll(profile,
			"include <abstractions/nameservice-strict>",
			"include <abstractions/attached/nameservice-strict>",
		)

	} else {
		if opt.Kind == aa.ProfileKind {
			insert = "@{att} = \"\"\n"
		}

	}

	return strings.Replace(profile, origin, insert+origin, 1), nil
}
