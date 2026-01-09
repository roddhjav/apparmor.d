// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package configure

import (
	"strings"

	"github.com/roddhjav/apparmor.d/pkg/tasks"
)

type ReAttach struct {
	tasks.BaseTask
}

// NewAttach creates a new ReAttach task.
func NewAttach() *ReAttach {
	return &ReAttach{
		BaseTask: tasks.BaseTask{
			Keyword: "attach",
			Msg:     "Configure tunable for re-attached path",
		},
	}
}

func (p ReAttach) Apply() ([]string, error) {
	res := []string{}

	// Remove the @{att} tunable that is going to be defined in profile header
	path := p.RootApparmor.Join("tunables/multiarch.d/system")
	out, err := path.ReadFileAsString()
	if err != nil {
		return res, err
	}
	out = strings.ReplaceAll(out, `@{att}=""`, `# @{att}=""`)
	return res, path.WriteFile([]byte(out))
}
