// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2023-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package prebuild

import (
	"strings"

	"github.com/roddhjav/apparmor.d/pkg/logging"
	"github.com/roddhjav/apparmor.d/pkg/paths"
	"github.com/roddhjav/apparmor.d/pkg/prebuild/builder"
	"github.com/roddhjav/apparmor.d/pkg/prebuild/cfg"
	"github.com/roddhjav/apparmor.d/pkg/prebuild/directive"
	"github.com/roddhjav/apparmor.d/pkg/prebuild/prepare"
	"github.com/roddhjav/apparmor.d/pkg/util"
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
	case "whonix":
		cfg.Hide += `/etc/apparmor.d/abstractions/base.d/kicksecure
/etc/apparmor.d/home.tor-browser.firefox
/etc/apparmor.d/tunables/home.d/anondist
/etc/apparmor.d/tunables/home.d/live-mode
/etc/apparmor.d/tunables/home.d/qubes-whonix-anondist
/etc/apparmor.d/usr.bin.hexchat
/etc/apparmor.d/usr.bin.sdwdate
/etc/apparmor.d/usr.bin.systemcheck
/etc/apparmor.d/usr.bin.timesanitycheck
/etc/apparmor.d/usr.bin.url_to_unixtime
/etc/apparmor.d/whonix-firewall
`
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
		profile, err := util.ReadFile(file)
		if err != nil {
			return err
		}
		profile, err = builder.Run(file, profile)
		if err != nil {
			return err
		}
		profile, err = directive.Run(file, profile)
		if err != nil {
			return err
		}
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
