// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package configure

import (
	"strings"

	"github.com/roddhjav/apparmor.d/pkg/prebuild"
	"github.com/roddhjav/apparmor.d/pkg/tasks"
)

type ReAttach struct {
	tasks.Base
}

func init() {
	RegisterTask(&ReAttach{
		Base: tasks.Base{
			Keyword: "attach",
			Msg:     "Configure tunable for re-attached path",
		},
	})
}

func (p ReAttach) Apply() ([]string, error) {
	res := []string{}

	// Remove the @{att} tunable that is going to be defined in profile header
	path := prebuild.RootApparmord.Join("tunables/multiarch.d/system")
	out, err := path.ReadFileAsString()
	if err != nil {
		return res, err
	}
	out = strings.ReplaceAll(out, `@{att}=""`, `# @{att}=""`)
	return res, path.WriteFile([]byte(out))
}
