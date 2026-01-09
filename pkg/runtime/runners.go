// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package runtime

import (
	"strings"

	"github.com/roddhjav/apparmor.d/pkg/builder"
	"github.com/roddhjav/apparmor.d/pkg/configure"
	"github.com/roddhjav/apparmor.d/pkg/directive"
	"github.com/roddhjav/apparmor.d/pkg/logging"
	"github.com/roddhjav/apparmor.d/pkg/paths"
	"github.com/roddhjav/apparmor.d/pkg/tasks"
)

// Runners groups all runners used during install or prebuild jobs.
type Runners struct {
	Configures *configure.Configures
	Builders   *builder.Builders
	Directives *directive.Directives
}

// NewRunners groups all runners used during install.
func NewRunners(c tasks.TaskConfig) *Runners {
	return &Runners{
		Configures: configure.NewRunner(c),
		Builders:   builder.NewRunner(c),
		Directives: directive.NewRunner(c),
	}
}

// Configure runs all configure tasks.
func (r *Runners) Configure() error {
	for _, task := range r.Configures.Tasks {
		msg, err := task.Apply()
		if err != nil {
			return err
		}
		logging.Success("%s", task.Message())
		logging.Indent = "   "
		for _, line := range msg {
			if strings.Contains(line, "not found") {
				logging.Warning("%s", line)
			} else {
				logging.Bullet("%s", line)
			}
		}
		logging.Indent = ""
	}
	return nil
}

// Build runs all build tasks and processes all directives.
func (r *Runners) Build() error {
	files, _ := r.Builders.RootApparmor.ReadDirRecursiveFiltered(nil, paths.FilterOutDirectories())
	for _, file := range files {
		if !file.Exist() {
			continue
		}
		profile, err := file.ReadFileAsString()
		if err != nil {
			return err
		}
		profile, err = r.Builders.Run(file, profile)
		if err != nil {
			return err
		}
		profile, err = r.Directives.Run(file, profile)
		if err != nil {
			return err
		}
		if err := file.WriteFile([]byte(profile)); err != nil {
			return err
		}
	}

	logging.Success("Build tasks:")
	logging.Indent = "   "
	for _, task := range r.Builders.Tasks {
		logging.Bullet("%s", task.Message())
	}
	if len(r.Directives.Directives) > 0 {
		logging.Indent = ""
		logging.Success("Directives processed:")
		logging.Indent = "   "
		for _, d := range r.Directives.Directives {
			logging.Bullet("%s%s", directive.Keyword, d.Name())
		}
		logging.Indent = ""
	}
	return nil
}
