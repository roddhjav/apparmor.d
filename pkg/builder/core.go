// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package builder

import (
	"fmt"
	"strings"

	"github.com/roddhjav/apparmor.d/pkg/aa"
	"github.com/roddhjav/apparmor.d/pkg/paths"
	"github.com/roddhjav/apparmor.d/pkg/tasks"
)

// Builder main directive interface
type Builder interface {
	tasks.BaseTaskInterface
	Apply(opt *Option, profile string) (string, error)
}

// Option for a builder
type Option struct {
	Name string
	File *paths.Path
	Kind aa.FileKind
}

func NewOption(file *paths.Path) *Option {
	return &Option{
		Name: strings.TrimSuffix(file.Base(), ".apparmor.d"),
		File: file,
		Kind: aa.KindFromPath(file),
	}
}

// Builders executes builders on profile strings in a pipeline.
type Builders struct {
	*tasks.BaseRunner[Builder]
}

// NewRunner creates a new Builders instance.
func NewRunner(t *tasks.TaskConfig) *Builders {
	return &Builders{
		BaseRunner: tasks.NewBaseRunner[Builder](t),
	}
}

// Run executes all builders on a profile string.
func (r *Builders) Run(file *paths.Path, profile string) (string, error) {
	opt := NewOption(file)
	var err error

	for _, b := range r.Tasks {
		profile, err = b.Apply(opt, profile)
		if err != nil {
			return "", fmt.Errorf("%s %s: %w", b.Name(), opt.File, err)
		}
	}
	return profile, nil
}

// Add appends a builder to the pipeline with fluent interface.
func (r *Builders) Add(builder Builder) *Builders {
	r.BaseRunner.Add(builder)
	return r
}
