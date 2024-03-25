// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package prepare

import (
	"github.com/roddhjav/apparmor.d/pkg/prebuild/cfg"
	"github.com/roddhjav/apparmor.d/pkg/util"
)

type SystemdDefault struct {
	cfg.Base
}

type SystemdEarly struct {
	cfg.Base
}

func init() {
	RegisterTask(&SystemdDefault{
		Base: cfg.Base{
			Keyword: "systemd-default",
			Msg:     "Configure systemd unit drop in files to a profile for some units",
		},
	})
	RegisterTask(&SystemdEarly{
		Base: cfg.Base{
			Keyword: "systemd-early",
			Msg:     "Configure systemd unit drop in files to ensure some service start after apparmor",
		},
	})
}

func (p SystemdDefault) Apply() ([]string, error) {
	return []string{}, util.CopyTo(cfg.SystemdDir.Join("default"), cfg.Root.Join("systemd"))
}

func (p SystemdEarly) Apply() ([]string, error) {
	return []string{}, util.CopyTo(cfg.SystemdDir.Join("early"), cfg.Root.Join("systemd"))
}
