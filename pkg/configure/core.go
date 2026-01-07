// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package configure

import (
	"fmt"

	"github.com/roddhjav/apparmor.d/pkg/tasks"
)

var (
	// Prepare the build directory with the following tasks
	Prepares = []Task{}

	// Available prepare tasks
	Tasks = map[string]Task{}
)

// Task main directive interface
type Task interface {
	tasks.BaseInterface
	Apply() ([]string, error)
}

func Register(names ...string) {
	for _, name := range names {
		if b, present := Tasks[name]; present {
			Prepares = append(Prepares, b)
		} else {
			panic(fmt.Sprintf("Unknown task: %s", name))
		}
	}
}

func RegisterTask(t Task) {
	Tasks[t.Name()] = t
}
