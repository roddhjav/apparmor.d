// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package builder

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/roddhjav/apparmor.d/pkg/aa"
	"github.com/roddhjav/apparmor.d/pkg/prebuild/cfg"
)

const tokATTACHMENT = "@{exec_path}"

var (
	regAttachments = regexp.MustCompile(`(profile .* ` + tokATTACHMENT + `)`)
)

type Userspace struct {
	cfg.Base
}

func init() {
	RegisterBuilder(&Userspace{
		Base: cfg.Base{
			Keyword: "userspace",
			Msg:     "Bypass userspace tools restriction",
		},
	})
}

func (b Userspace) Apply(opt *Option, profile string) (string, error) {
	if ok, _ := opt.File.IsInsideDir(cfg.RootApparmord.Join("abstractions")); ok {
		return profile, nil
	}
	if ok, _ := opt.File.IsInsideDir(cfg.RootApparmord.Join("tunables")); ok {
		return profile, nil
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
		strheader := strings.Replace(matches[0], tokATTACHMENT, att, -1)
		return regAttachments.ReplaceAllLiteralString(profile, strheader), nil
	}
	return profile, nil
}
