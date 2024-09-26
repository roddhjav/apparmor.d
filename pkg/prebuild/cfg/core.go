// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package cfg

import "fmt"

type BaseInterface interface {
	Message() string
	Name() string
	Usage() []string
}

type Base struct {
	Msg     string
	Keyword string
	Help    []string
}

func (b Base) Name() string {
	return b.Keyword
}

func (b Base) Usage() []string {
	return b.Help
}

func (b Base) Message() string {
	return b.Msg
}

func Help[T BaseInterface](name string, tasks map[string]T) string {
	res := fmt.Sprintf("%s tasks:\n", name)
	for _, t := range tasks {
		res += fmt.Sprintf("    %s - %s\n", t.Name(), t.Message())
	}
	return res
}

func Usage[T BaseInterface](name string, tasks map[string]T) string {
	res := fmt.Sprintf("%s\n", name)
	for _, t := range tasks {
		for _, h := range t.Usage() {
			res += fmt.Sprintf("    #aa:%s %s\n", t.Name(), h)
		}
	}
	return res
}
