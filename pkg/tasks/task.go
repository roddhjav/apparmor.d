// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package tasks

import (
	"github.com/roddhjav/apparmor.d/pkg/paths"
)

type TaskConfig struct {
	Root         *paths.Path // Root is the root directory for the runner (e.g. .build)
	RootApparmor *paths.Path // RootApparmor is the source apparmor.d directory (e.g. .build/apparmor.d)
}

func NewTaskConfig(root *paths.Path) TaskConfig {
	return TaskConfig{
		Root:         root,
		RootApparmor: root.Join("apparmor.d"),
	}
}

type BaseTaskInterface interface {
	Message() string
	Name() string
	Usage() []string
	SetConfig(c TaskConfig)
}

type BaseTask struct {
	TaskConfig
	Msg     string
	Keyword string
	Help    []string
}

func (b BaseTask) Name() string {
	return b.Keyword
}

func (b *BaseTask) SetConfig(c TaskConfig) {
	b.TaskConfig = c
}

func (b BaseTask) Usage() []string {
	return b.Help
}

func (b BaseTask) Message() string {
	return b.Msg
}
