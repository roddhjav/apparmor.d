// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2023-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package prebuild

import (
	"strings"

	"github.com/arduino/go-paths-helper"
	"github.com/roddhjav/apparmor.d/pkg/logging"
	"github.com/roddhjav/apparmor.d/pkg/prebuild/builder"
	"github.com/roddhjav/apparmor.d/pkg/prebuild/cfg"
	"github.com/roddhjav/apparmor.d/pkg/prebuild/directive"
	"github.com/roddhjav/apparmor.d/pkg/prebuild/prepare"
)

func init() {
	// Define the tasks applied by default
	prepare.Register(
		"synchronise",
		"ignore",
		"merge",
		"configure",
		"setflags",
		"systemd-default",
	)

	// Build tasks applied by default
	builder.Register("userspace")
	builder.Register("dev")

	switch cfg.Distribution {
	case "ubuntu":
		if cfg.Release["VERSION_CODENAME"] == "noble" {
			builder.Register("abi3")
			cfg.Overwrite.Enabled = true
		}
	}
}

func Prepare() error {
	for _, task := range prepare.Prepares {
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

func Build() error {
	files, _ := cfg.RootApparmord.ReadDirRecursiveFiltered(nil, paths.FilterOutDirectories())
	for _, file := range files {
		if !file.Exist() {
			continue
		}
		content, err := file.ReadFile()
		if err != nil {
			return err
		}
		profile := string(content)
		for _, b := range builder.Builds {
			profile = b.Apply(profile)
		}
		profile = directive.Run(file, profile)
		if err := file.WriteFile([]byte(profile)); err != nil {
			return err
		}
	}

	logging.Success("Build tasks:")
	logging.Indent = "   "
	for _, task := range builder.Builds {
		logging.Bullet("%s", task.Message())
	}
	logging.Indent = ""
	logging.Success("Directives processed:")
	logging.Indent = "   "
	for _, dir := range directive.Directives {
		logging.Bullet("%s%s", directive.Keyword, dir.Name())
	}
	logging.Indent = ""
	return nil
}
