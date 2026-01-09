// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package tasks

import (
	"fmt"
)

// Runner provides a fluent interface for building and executing task pipelines.
type Runner[T BaseTaskInterface] interface {
	// Add appends a task to the execution pipeline
	Add(task T) Runner[T]

	// Run executes all tasks in order, stopping on first error
	Run() error

	// Help returns usage information for all registered tasks
	Help(name string) string
}

// BaseRunner provides common runner implementation for task pipelines.
type BaseRunner[T BaseTaskInterface] struct {
	TaskConfig
	Tasks []T
}

// NewBaseRunner creates a new BaseRunner instance.
func NewBaseRunner[T BaseTaskInterface](config TaskConfig) *BaseRunner[T] {
	r := &BaseRunner[T]{
		TaskConfig: config,
		Tasks:      make([]T, 0),
	}
	return r
}

// Add appends a task to the execution pipeline.
func (r *BaseRunner[T]) Add(task T) *BaseRunner[T] {
	task.SetConfig(r.TaskConfig)
	r.Tasks = append(r.Tasks, task)
	return r
}

func (r *BaseRunner[T]) Help(name string) string {
	res := fmt.Sprintf("%s tasks:\n", name)
	for _, t := range r.Tasks {
		res += fmt.Sprintf("    %s - %s\n", t.Name(), t.Message())
	}
	return res
}

// Run is not implemented in BaseRunner - concrete runners must implement their own Run method.
