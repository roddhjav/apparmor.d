// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package tasks

type BaseTaskInterface interface {
	Message() string
	Name() string
	Usage() []string
	SetConfig(c *TaskConfig)
}

type BaseTask struct {
	*TaskConfig
	Msg     string
	Keyword string
	Help    []string
}

func (b BaseTask) Name() string {
	return b.Keyword
}

func (b *BaseTask) SetConfig(c *TaskConfig) {
	b.TaskConfig = c
}

func (b BaseTask) Usage() []string {
	return b.Help
}

func (b BaseTask) Message() string {
	return b.Msg
}
