// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package configure

import (
	"fmt"

	"github.com/roddhjav/apparmor.d/pkg/logging"
	"github.com/roddhjav/apparmor.d/pkg/tasks"
)

// Task main directive interface
type Task interface {
	tasks.BaseTaskInterface
	Apply() ([]string, error)
}

// Configures executes configure tasks in a pipeline.
type Configures struct {
	*tasks.BaseRunner[Task]
}

// NewRunner creates a new Configures instance.
func NewRunner(t *tasks.TaskConfig) *Configures {
	return &Configures{
		BaseRunner: tasks.NewBaseRunner[Task](t),
	}
}

// Run executes all tasks in the pipeline, logging their output.
func (r *Configures) Run() error {
	for _, task := range r.Tasks {
		msg, err := task.Apply()
		if err != nil {
			return fmt.Errorf("%s: %w", task.Name(), err)
		}
		logging.Success("%s", task.Message())
		for _, m := range msg {
			logging.Bullet("%s", m)
		}
	}
	return nil
}

// Add appends a task to the pipeline with fluent interface.
func (r *Configures) Add(task Task) *Configures {
	r.BaseRunner.Add(task)
	return r
}
