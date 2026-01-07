// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package builder

import (
	"fmt"
	"strings"

	"github.com/roddhjav/apparmor.d/pkg/aa"
	"github.com/roddhjav/apparmor.d/pkg/prebuild"
	"github.com/roddhjav/apparmor.d/pkg/tasks"
)

type ReAttach struct {
	tasks.Base
}

func init() {
	RegisterBuilder(&ReAttach{
		Base: tasks.Base{
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

	isInside, err := opt.File.IsInsideDir(prebuild.RootApparmord.Join("abstractions/attached"))
	if err != nil {
		return profile, fmt.Errorf("attach: %v", err)
	}
	if isInside {
		return profile, nil // Do not re-attach twice
	}

	if strings.Contains(profile, "attach_disconnected") {
		if opt.Kind == aa.ProfileKind {
			if strings.Contains(opt.Name, ":") {
				parts := strings.Split(opt.Name, ":")
				if len(parts) != 3 {
					return profile, fmt.Errorf("attach: invalid namespaced profile name: %s", opt.Name)
				}
				insert = "@{att} = /att/" + parts[1] + "/\n"
			} else {
				insert = "@{att} = /att/" + opt.Name + "/\n"
			}
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
