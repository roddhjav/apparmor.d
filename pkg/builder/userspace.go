// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package builder

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/roddhjav/apparmor.d/pkg/aa"
	"github.com/roddhjav/apparmor.d/pkg/prebuild"
	"github.com/roddhjav/apparmor.d/pkg/tasks"
)

const tokATTACHMENT = "@{exec_path}"

var (
	regAttachments = regexp.MustCompile(`(profile .* ` + tokATTACHMENT + `)`)
)

type Userspace struct {
	tasks.Base
}

func init() {
	RegisterBuilder(&Userspace{
		Base: tasks.Base{
			Keyword: "userspace",
			Msg:     "Fix: resolve variable in profile attachments",
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
	if prebuild.Distribution == "arch" {
		f.Preamble = append(f.Preamble, &aa.Variable{
			Name: "sbin", Values: []string{"/{,usr/}{,s}bin"}, Define: true,
		})
	} else {
		f.Preamble = append(f.Preamble, &aa.Variable{
			Name: "sbin", Values: []string{"/{,usr/}sbin"}, Define: true,
		})
	}

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
