// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2023 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"bytes"
	"reflect"
	"slices"
	"sort"
	"strings"

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
	Preamble
	Profiles []*Profile
}

// Preamble section of a profile file,
type Preamble struct {
	Abi       []*Abi
	Includes  []*Include
	Aliases   []*Alias
	Variables []*Variable
	Comments  []*RuleBase
}

func NewAppArmorProfile() *AppArmorProfileFile {
	return &AppArmorProfileFile{}
}

// String returns the formatted representation of a profile as a string
func (f *AppArmorProfileFile) String() string {
	var res bytes.Buffer
	err := tmpl["apparmor"].Execute(&res, f)
	if err != nil {
		return err.Error()
	}
	return res.String()
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
		sort.Slice(p.Rules, func(i, j int) bool {
			typeOfI := reflect.TypeOf(p.Rules[i])
			typeOfJ := reflect.TypeOf(p.Rules[j])
			if typeOfI != typeOfJ {
				valueOfI := typeToValue(typeOfI)
				valueOfJ := typeToValue(typeOfJ)
				if typeOfI == reflect.TypeOf((*Include)(nil)) && p.Rules[i].(*Include).IfExists {
					valueOfI = "include_if_exists"
				}
				if typeOfJ == reflect.TypeOf((*Include)(nil)) && p.Rules[j].(*Include).IfExists {
					valueOfJ = "include_if_exists"
				}
				return ruleWeights[valueOfI] < ruleWeights[valueOfJ]
			}
			return p.Rules[i].Less(p.Rules[j])
		})
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
		for i := 0; i < len(p.Rules); i++ {
			for j := i + 1; j < len(p.Rules); j++ {
				typeOfI := reflect.TypeOf(p.Rules[i])
				typeOfJ := reflect.TypeOf(p.Rules[j])
				if typeOfI != typeOfJ {
					continue
				}

				// If rules are identical, merge them
				if p.Rules[i].Equals(p.Rules[j]) {
					p.Rules = append(p.Rules[:j], p.Rules[j+1:]...)
					j--
				}
			}
		}
	}
}

// Format the profile for better readability before printing it.
// Follow: https://apparmor.pujol.io/development/guidelines/#the-file-block
func (f *AppArmorProfileFile) Format() {
	const prefixOwner = "      "
	for _, p := range f.Profiles {
		hasOwnerRule := false
		for i := len(p.Rules) - 1; i > 0; i-- {
			j := i - 1
			typeOfI := reflect.TypeOf(p.Rules[i])
			typeOfJ := reflect.TypeOf(p.Rules[j])

			// File rule
			if typeOfI == reflect.TypeOf((*File)(nil)) && typeOfJ == reflect.TypeOf((*File)(nil)) {
				letterI := getLetterIn(fileAlphabet, p.Rules[i].(*File).Path)
				letterJ := getLetterIn(fileAlphabet, p.Rules[j].(*File).Path)

				// Add prefix before rule path to align with other rule
				if p.Rules[i].(*File).Owner {
					hasOwnerRule = true
				} else if hasOwnerRule {
					p.Rules[i].(*File).Prefix = prefixOwner
				}

				if letterI != letterJ {
					// Add a new empty line between Files rule of different type
					hasOwnerRule = false
					p.Rules = append(p.Rules[:i], append([]ApparmorRule{&Rule{}}, p.Rules[i:]...)...)
				}
			}
		}
	}
}
