// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package detector

func init() {
	Register(displayServer{})
}

// displayServer detects @{DS} from the type of the logind sessions.
type displayServer struct{}

func (displayServer) Name() string { return "DS" }

func (displayServer) Detect(sys *System) []string {
	var values []string
	for _, value := range sys.sessionProperty("Type") {
		if value == "wayland" || value == "x11" {
			values = append(values, value)
		}
	}
	return dedup(values)
}
