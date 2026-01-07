// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package configure

import (
	"fmt"
	"strings"

	"github.com/roddhjav/apparmor.d/pkg/prebuild"
	"github.com/roddhjav/apparmor.d/pkg/tasks"
)

type Configure struct {
	tasks.Base
}

func init() {
	RegisterTask(&Configure{
		Base: tasks.Base{
			Keyword: "configure",
			Msg:     "Set distribution specificities",
		},
	})
}

func removeFiles(files []string) error {
	for _, name := range files {
		if err := prebuild.RootApparmord.Join(name).RemoveAll(); err != nil {
			return err
		}
	}
	return nil
}

func (p Configure) Apply() ([]string, error) {
	res := []string{}

	switch prebuild.Distribution {
	case "arch", "opensuse":

	case "ubuntu":
		if err := prebuild.DebianHide.Init(); err != nil {
			return res, err
		}

		if prebuild.Release["VERSION_CODENAME"] == "noble" {
			remove := []string{
				"tunables/multiarch.d/base",
			}
			if err := removeFiles(remove); err != nil {
				return res, err
			}
		}

	case "debian", "whonix":
		if err := prebuild.DebianHide.Init(); err != nil {
			return res, err
		}

	default:
		return []string{}, fmt.Errorf("%s is not a supported distribution", prebuild.Distribution)

	}

	if prebuild.Version < 4.1 {
		remove := []string{
			// Require priority support
			"fbwrap",
			"fapp",
		}
		if err := removeFiles(remove); err != nil {
			return res, err
		}
	}
	if prebuild.Version >= 4.1 {
		remove := []string{
			// Remove files upstreamed in 4.1
			"abstractions/devices-usb-read",
			"abstractions/devices-usb",
			"abstractions/nameservice-strict",
			"tunables/multiarch.d/base",

			// Direct upstream contributed profiles, similar to ours
			"wg",
		}
		if err := removeFiles(remove); err != nil {
			return res, err
		}
	}
	if prebuild.Version >= 5.0 {
		remove := []string{
			// Direct upstream contributed profiles, similar to ours
			"dig",
			"free",
			"nslookup",
		}
		if err := removeFiles(remove); err != nil {
			return res, err
		}

		// @{pci_bus} was upstreamed in 5.0
		path := prebuild.RootApparmord.Join("tunables/multiarch.d/system")
		out, err := path.ReadFileAsString()
		if err != nil {
			return res, err
		}
		out = strings.ReplaceAll(out, "@{pci_bus}=pci@{hex4}:@{hex2}", "")
		return res, path.WriteFile([]byte(out))
	}
	return res, nil
}
