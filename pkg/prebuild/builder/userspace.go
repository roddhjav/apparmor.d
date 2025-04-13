// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package builder

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/roddhjav/apparmor.d/pkg/aa"
	"github.com/roddhjav/apparmor.d/pkg/prebuild"
)

const tokATTACHMENT = "@{exec_path}"

var (
	regAttachments = regexp.MustCompile(`(profile .* ` + tokATTACHMENT + `)`)
)

type Userspace struct {
	prebuild.Base
}

func init() {
	RegisterBuilder(&Userspace{
		Base: prebuild.Base{
			Keyword: "userspace",
			Msg:     "Resolve variable in profile attachments",
		},
	})
}

func (b Userspace) Apply(opt *Option, profile string) (string, error) {
	for _, dir := range []string{"abstractions", "tunables", "local", "mappings"} {
		if ok, _ := opt.File.IsInsideDir(prebuild.RootApparmord.Join(dir)); ok {
			return profile, nil
		}
	}

	f := aa.DefaultTunables()
	if _, err := f.Parse(profile); err != nil {
		return "", err
	}
	if len(f.GetDefaultProfile().Attachments) > 0 &&
		f.GetDefaultProfile().Attachments[0] != tokATTACHMENT {
		return "", fmt.Errorf("missing '%s' attachment", tokATTACHMENT)
	}
	if err := f.Resolve(); err != nil {
		return "", err
	}

	matches := regAttachments.FindAllString(profile, -1)
	if len(matches) > 0 {
		att := f.GetDefaultProfile().GetAttachments()
		strheader := strings.ReplaceAll(matches[0], tokATTACHMENT, att)
		return regAttachments.ReplaceAllLiteralString(profile, strheader), nil
	}
	return profile, nil
}
