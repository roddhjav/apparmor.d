// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package detector

import (
	"path"
	"strings"

	"github.com/roddhjav/apparmor.d/pkg/aa"
)

func init() {
	Register(abi{})
}

// abi detects @{ABI} from the abi rule of the abstractions/base-strict file
// in the tree being written.
type abi struct{}

func (abi) Name() string { return "ABI" }

func (abi) Detect(sys *System) []string {
	lines, err := sys.Apparmor.Join("abstractions/base-strict").ReadFileAsLines()
	if err != nil {
		return nil
	}
	// The aa parser only consumes preamble rules, not an abstraction
	// body: feed it the abi lines only.
	var abiLines []string
	for _, line := range lines {
		if line = strings.TrimSpace(line); strings.HasPrefix(line, "abi ") {
			abiLines = append(abiLines, line)
		}
	}
	f := &aa.AppArmorProfileFile{}
	if _, err := f.Parse(strings.Join(abiLines, "\n") + "\n"); err != nil {
		return nil
	}
	for _, rule := range f.Preamble.Filter(aa.ABI) {
		major, _, _ := strings.Cut(path.Base(rule.(*aa.Abi).Path), ".")
		return []string{major}
	}
	return nil
}
