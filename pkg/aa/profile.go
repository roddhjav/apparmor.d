// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"bytes"
	"reflect"
	"sort"
	"strings"

	"golang.org/x/exp/slices"
)

// AppArmorProfiles represents a full set of apparmor profiles
type AppArmorProfiles map[string]*AppArmorProfile

// ApparmorProfile represents a full apparmor profile.
// Warning: close to the BNF grammar of apparmor profile but not exactly the same (yet):
//   - Some rules are not supported yet (subprofile, hat...)
//   - The structure is simplified as it only aims at writing profile, not parsing it.
type AppArmorProfile struct {
	Preamble
	Profile
}

// Preamble section of a profile
type Preamble struct {
	Abi       []Abi
	Includes  []Include
	Aliases   []Alias
	Variables []Variable
}

// Profile section of a profile
type Profile struct {
	Name        string
	Attachments []string
	Attributes  map[string]string
	Flags       []string
	Rules       Rules
}

// ApparmorRule generic interface
type ApparmorRule interface {
	Less(other any) bool
	Equals(other any) bool
}

type Rules []ApparmorRule

func NewAppArmorProfile() *AppArmorProfile {
	return &AppArmorProfile{}
}

// String returns the formatted representation of a profile as a string
func (p *AppArmorProfile) String() string {
	var res bytes.Buffer
	err := tmplAppArmorProfile.Execute(&res, p)
	if err != nil {
		return err.Error()
	}
	return res.String()
}

// AddRule adds a new rule to the profile from a log map
func (p *AppArmorProfile) AddRule(log map[string]string) {
	// Generate profile flags and extra rules
	switch log["error"] {
	case "-2":
		if !slices.Contains(p.Flags, "mediate_deleted") {
			p.Flags = append(p.Flags, "mediate_deleted")
		}
	case "-13":
		if strings.Contains(log["info"], "namespace creation restricted") {
			p.Rules = append(p.Rules, UsernsFromLog(log))
		} else if strings.Contains(log["info"], "disconnected path") && !slices.Contains(p.Flags, "attach_disconnected") {
			p.Flags = append(p.Flags, "attach_disconnected")
		}
	default:
	}

	switch log["class"] {
	case "cap":
		p.Rules = append(p.Rules, CapabilityFromLog(log))
	case "net":
		if log["family"] == "unix" {
			p.Rules = append(p.Rules, UnixFromLog(log))
		} else {
			p.Rules = append(p.Rules, NetworkFromLog(log))
		}
	case "mount":
		switch log["operation"] {
		case "mount":
			p.Rules = append(p.Rules, MountFromLog(log))
		case "umount":
			p.Rules = append(p.Rules, UmountFromLog(log))
		case "remount":
			p.Rules = append(p.Rules, RemountFromLog(log))
		case "pivotroot":
			p.Rules = append(p.Rules, PivotRootFromLog(log))
		}
	case "posix_mqueue", "sysv_mqueue":
		p.Rules = append(p.Rules, MqueueFromLog(log))
	case "signal":
		p.Rules = append(p.Rules, SignalFromLog(log))
	case "ptrace":
		p.Rules = append(p.Rules, PtraceFromLog(log))
	case "namespace":
		p.Rules = append(p.Rules, UsernsFromLog(log))
	case "unix":
		p.Rules = append(p.Rules, UnixFromLog(log))
	case "file":
		if log["operation"] == "change_onexec" {
			p.Rules = append(p.Rules, ChangeProfileFromLog(log))
		} else {
			p.Rules = append(p.Rules, FileFromLog(log))
		}
	default:
		if strings.Contains(log["operation"], "dbus") {
			p.Rules = append(p.Rules, DbusFromLog(log))
		} else if log["family"] == "unix" {
			p.Rules = append(p.Rules, UnixFromLog(log))
		}
	}
}

// Sort the rules in the profile
// Follow: https://apparmor.pujol.io/development/guidelines/#guidelines
func (p *AppArmorProfile) Sort() {
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

// MergeRules merge similar rules together
// Steps:
//   - Remove identical rules
//   - Merge rule access. Eg: for same path, 'r' and 'w' becomes 'rw'
//
// Note: logs.regCleanLogs helps a lot to do a first cleaning
func (p *AppArmorProfile) MergeRules() {
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

// Format the profile for better readability before printing it
// Follow: https://apparmor.pujol.io/development/guidelines/#the-file-block
func (p *AppArmorProfile) Format() {
	hasOwnedRule := false
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
				hasOwnedRule = true
			} else if hasOwnedRule {
				p.Rules[i].(*File).Prefix = "      "
			}

			if letterI != letterJ {
				// Add a new empty line between Files rule of different type
				hasOwnedRule = false
				p.Rules = append(p.Rules[:i], append([]ApparmorRule{&Rule{}}, p.Rules[i:]...)...)
			}
		}
	}
}
