// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package builder

import (
	"regexp"
	"strings"

	"github.com/roddhjav/apparmor.d/pkg/tasks"
)

var (
	regDebug = regexp.MustCompile(`(?m)^([ \t]*)(.*)(pi|Pi|pu|PU|p|P|C|c)x(.*),(.*)$`)
)

type Debug struct {
	tasks.BaseTask
}

// NewDebug creates a new Debug builder.
func NewDebug() *Debug {
	return &Debug{
		BaseTask: tasks.BaseTask{
			Keyword: "debug",
			Msg:     "Build: debug mode enabled",
		},
	}
}

func (b Debug) Apply(opt *Option, profile string) (string, error) {
	for _, dir := range []string{"tunables"} {
		if ok, _ := opt.File.IsInsideDir(b.RootApparmor.Join(dir)); ok {
			return profile, nil
		}
	}

	lines := strings.Split(profile, "\n")
	for i, line := range lines {
		trimmed := strings.TrimLeft(line, " \t")
		if strings.HasPrefix(trimmed, "#") {
			continue
		}
		if strings.Contains(trimmed, "=") {
			continue
		}
		if strings.HasPrefix(trimmed, "audit") {
			continue
		}
		lines[i] = regDebug.ReplaceAllString(line, `${1}audit ${2}${3}x${4},${5}`)
	}
	profile = strings.Join(lines, "\n")
	return profile, nil
}
