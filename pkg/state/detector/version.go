// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package detector

import "regexp"

func init() {
	Register(version{})
}

var regParserVersion = regexp.MustCompile(`version\s+(\d+)`)

// version detects @{VERSION} from the running `apparmor_parser -V`.
type version struct{}

func (version) Name() string { return "VERSION" }

func (version) Detect(sys *System) []string {
	out, err := sys.Run("apparmor_parser", "-V")
	if err != nil {
		return nil
	}
	if m := regParserVersion.FindStringSubmatch(out); m != nil {
		return []string{m[1]}
	}
	return nil
}
