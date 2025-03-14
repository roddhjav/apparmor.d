// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package builder

import (
	"regexp"
	"strings"

	"github.com/roddhjav/apparmor.d/pkg/prebuild"
)

var (
	regProfile = regexp.MustCompile(`profile ([^ ]+)`)
)

type ReAttach struct {
	prebuild.Base
}

func init() {
	RegisterBuilder(&ReAttach{
		Base: prebuild.Base{
			Keyword: "attach",
			Msg:     "Re-attach disconnected path",
		},
	})
}

// Apply will re-attach the disconnected path
//   - Add the attach_disconnected.path flag on all frofile with the attach_disconnected flag
//   - Add the attached/base abstraction in the profile
//   - For compatibility, non disconnected profile will have the @{att} variable set to /
func (b ReAttach) Apply(opt *Option, profile string) (string, error) {
	var insert string
	var origin = "profile " + opt.Name

	if strings.Contains(profile, "attach_disconnected") {
		insert = "@{att} = /att/" + opt.Name + "/\n"
		profile = strings.Replace(profile,
			"attach_disconnected",
			"attach_disconnected,attach_disconnected.path=@{att}", -1,
		)

		old := "include if exists <local/" + opt.Name + ">"
		new := "include <abstractions/attached/base>\n  " + old
		profile = strings.Replace(profile, old, new, 1)

		for _, match := range regProfile.FindAllStringSubmatch(profile, -1) {
			name := match[1]
			if name == opt.Name {
				continue
			}
			old = "include if exists <local/" + opt.Name + "_" + name + ">"
			new = "include <abstractions/attached/base>\n    " + old
			profile = strings.Replace(profile, old, new, 1)
		}

	} else {
		insert = "@{att} = /\n"
	}

	return strings.Replace(profile, origin, insert+origin, 1), nil
}
