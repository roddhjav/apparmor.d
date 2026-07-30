// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package detector

import (
	"os"
	"strings"
)

func init() {
	Register(desktop{})
}

var (
	// deNames normalizes desktop names as reported by logind
	deNames = []struct{ substr, de string }{
		{"gnome", "gnome"},
		{"plasma", "kde"},
		{"kde", "kde"},
		{"xfce", "xfce"},
		{"cosmic", "cosmic"},
	}

	// deShells maps a desktop shell process name to its desktop environment.
	deShells = map[string]string{
		"gnome-shell":    "gnome",
		"plasmashell":    "kde",
		"xfce4-session":  "xfce",
		"cosmic-session": "cosmic",
	}
)

// desktop detects @{DE} from the desktops the logind sessions report,
// falling back to the desktop shell processes running on the system.
type desktop struct{}

func (desktop) Name() string { return "DE" }

func (desktop) Detect(sys *System) []string {
	var values []string
	for _, name := range sys.sessionProperty("Desktop") {
		values = append(values, desktopNames(name)...)
	}
	if len(values) == 0 {
		values = runningDesktops(sys)
	}
	return dedup(values)
}

// desktopNames normalizes a logind Desktop value, a colon separated list of
// desktop names, to the documented @{DE} values.
func desktopNames(value string) []string {
	var des []string
	for field := range strings.SplitSeq(strings.ToLower(value), ":") {
		for _, m := range deNames {
			if strings.Contains(field, m.substr) {
				des = append(des, m.de)
				break
			}
		}
	}
	return des
}

// runningDesktops identifies the running desktop shell processes. Not every
// session reports its desktop to logind.
func runningDesktops(sys *System) []string {
	var des []string
	for _, comm := range sortedGlob(sys.Root.Join("proc/*/comm").String()) {
		data, err := os.ReadFile(comm)
		if err != nil {
			continue
		}
		if de := deShells[strings.TrimSpace(string(data))]; de != "" {
			des = append(des, de)
		}
	}
	return des
}
