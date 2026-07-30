// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package detector

import (
	"os"
	"path/filepath"
	"strings"
)

func init() {
	Register(displayManager{})
}

var dmAliases = map[string]string{
	"gdm3": "gdm",
}

// displayManager detects @{DM} from the enabled display-manager.service
// systemd unit.
type displayManager struct{}

func (displayManager) Name() string { return "DM" }

func (displayManager) Detect(sys *System) []string {
	target, err := os.Readlink(sys.Root.Join("etc/systemd/system/display-manager.service").String())
	if err != nil {
		return []string{"none"}
	}
	name := strings.TrimSuffix(filepath.Base(target), ".service")
	if alias, ok := dmAliases[name]; ok {
		name = alias
	}
	return []string{name}
}
