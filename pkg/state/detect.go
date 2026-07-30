// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package state

import (
	"strconv"
	"strings"

	"github.com/roddhjav/apparmor.d/pkg/logging"
	"github.com/roddhjav/apparmor.d/pkg/paths"
	"github.com/roddhjav/apparmor.d/pkg/state/detector"
)

// Enabled reports whether the state system applies to the tree sys inspects:
// it is only enabled on abi 5 and over.
func Enabled(sys *detector.System) bool {
	values := detector.Registry["ABI"].Detect(sys)
	if len(values) == 0 {
		return false
	}
	abi, err := strconv.Atoi(values[0])
	return err == nil && abi >= 5
}

// Detect runs the detector named after every variable defined in the src
// state file and returns the detected state as a file at dst, to be saved
// as a state.d drop-in. Variables without a detector are logged and
// skipped; detectors without a detected value are skipped.
func Detect(sys *detector.System, src *paths.Path, dst *paths.Path) (*File, error) {
	stateFile, err := Load(src)
	if err != nil {
		return nil, err
	}
	res := NewFile(dst)
	for _, key := range stateFile.Names() {
		d, ok := detector.Registry[varName(key)]
		if !ok {
			logging.Warning("No state detector for %s", key)
			continue
		}
		if values := d.Detect(sys); len(values) > 0 {
			res.Add(key, strings.Join(values, " "))
		}
	}
	return res, nil
}
