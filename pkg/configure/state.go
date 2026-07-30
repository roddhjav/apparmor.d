// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package configure

import (
	"fmt"

	"github.com/roddhjav/apparmor.d/pkg/state"
	"github.com/roddhjav/apparmor.d/pkg/state/detector"
	"github.com/roddhjav/apparmor.d/pkg/tasks"
)

const (
	// stateRel is the state tunable, relative to relative to aa.Magic.
	stateRel = "tunables/multiarch.d/state"

	// stateDropinRel is the generated detected state drop-in, relative to aa.Magic.
	stateDropinRel = "tunables/multiarch.d/state.d/aa-install"
)

type SetState struct {
	tasks.BaseTask
}

// NewSetState creates a new SetState task.
func NewSetState() *SetState {
	return &SetState{
		BaseTask: tasks.BaseTask{
			Keyword: "state",
			Msg:     "Set detected system state in " + stateDropinRel,
		},
	}
}

func (p SetState) Apply() ([]string, error) {
	res := []string{}
	src := p.RootApparmor.Join(stateRel)
	if !src.Exist() {
		return res, nil
	}
	sys := detector.NewSystem(p.RootApparmor)
	if !state.Enabled(sys) {
		return res, nil
	}
	file, err := state.Detect(sys, src, p.RootApparmor.Join(stateDropinRel))
	if err != nil {
		return res, err
	}
	for _, name := range file.Names() {
		value, _ := file.Get(name)
		res = append(res, fmt.Sprintf("%s = %s", name, value))
	}
	return res, file.Save()
}
