// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package detector

import (
	"os"
	"strings"
)

func init() {
	Register(gpu{})
}

var gpuVendors = map[string]string{
	"0x10de": "nvidia",
	"0x1002": "amd",
	"0x8086": "intel",
}

// gpu detects @{GPU} from the vendor of the cards in /sys/class/drm.
type gpu struct{}

func (gpu) Name() string { return "GPU" }

func (gpu) Detect(sys *System) []string {
	var vendors []string
	for _, card := range sortedGlob(sys.Root.Join("sys/class/drm/*/device/vendor").String()) {
		data, err := os.ReadFile(card)
		if err != nil {
			continue
		}
		if vendor := gpuVendors[strings.TrimSpace(string(data))]; vendor != "" {
			vendors = append(vendors, vendor)
		}
	}
	if len(vendors) == 0 {
		return []string{"none"}
	}
	return dedup(vendors)
}
