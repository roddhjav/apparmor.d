// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2023 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"github.com/arduino/go-paths-helper"
)

// Default Apparmor magic directory: /etc/apparmor.d/.
var MagicRoot = paths.New("/etc/apparmor.d")

// AppArmorProfileFiles represents a full set of apparmor profiles
type AppArmorProfileFiles map[string]*AppArmorProfileFile

// AppArmorProfileFile represents a full apparmor profile file.
// Warning: close to the BNF grammar of apparmor profile but not exactly the same (yet):
//   - Some rules are not supported yet (subprofile, hat...)
//   - The structure is simplified as it only aims at writing profile, not parsing it.
type AppArmorProfileFile struct {
	Preamble Rules
	Profiles []*Profile
}

func NewAppArmorProfile() *AppArmorProfileFile {
	return &AppArmorProfileFile{}
}

// String returns the formatted representation of a profile file as a string
func (f *AppArmorProfileFile) String() string {
	return renderTemplate("apparmor", f)
}

// GetDefaultProfile ensure a profile is always present in the profile file and
// return it, as a default profile.
func (f *AppArmorProfileFile) GetDefaultProfile() *Profile {
	if len(f.Profiles) == 0 {
		f.Profiles = append(f.Profiles, &Profile{})
	}
	return f.Profiles[0]
}

// Sort the rules in the profile
// Follow: https://apparmor.pujol.io/development/guidelines/#guidelines
func (f *AppArmorProfileFile) Sort() {
	for _, p := range f.Profiles {
		p.Sort()
	}
}

// MergeRules merge similar rules together.
// Steps:
//   - Remove identical rules
//   - Merge rule access. Eg: for same path, 'r' and 'w' becomes 'rw'
//
// Note: logs.regCleanLogs helps a lot to do a first cleaning
func (f *AppArmorProfileFile) MergeRules() {
	for _, p := range f.Profiles {
		p.Merge()
	}
}

// Format the profile for better readability before printing it.
// Follow: https://apparmor.pujol.io/development/guidelines/#the-file-block
func (f *AppArmorProfileFile) Format() {
	for _, p := range f.Profiles {
		p.Format()
	}
}
