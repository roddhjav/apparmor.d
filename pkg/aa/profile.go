// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"maps"
	"reflect"
	"slices"
	"strings"
)

const (
	tokATTRIBUTES = "xattrs"
	tokFLAGS      = "flags"
	tokPROFILE    = "profile"
)

// Profile represents a single AppArmor profile.
type Profile struct {
	RuleBase
	Header
	Rules Rules
}

// Header represents the header of a profile.
type Header struct {
	Name        string
	Attachments []string
	Attributes  map[string]string
	Flags       []string
}

func (p *Profile) Less(other any) bool {
	o, _ := other.(*Profile)
	if p.Name != o.Name {
		return p.Name < o.Name
	}
	return len(p.Attachments) < len(o.Attachments)
}

func (p *Profile) Equals(other any) bool {
	o, _ := other.(*Profile)
	return p.Name == o.Name && slices.Equal(p.Attachments, o.Attachments) &&
		maps.Equal(p.Attributes, o.Attributes) &&
		slices.Equal(p.Flags, o.Flags)
}

func (p *Profile) String() string {
	return renderTemplate(tokPROFILE, p)
}

// Merge merge similar rules together.
// Steps:
//   - Remove identical rules
//   - Merge rule access. Eg: for same path, 'r' and 'w' becomes 'rw'
//
// Note: logs.regCleanLogs helps a lot to do a first cleaning
func (p *Profile) Merge() {
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

// Sort the rules in a profile.
// Follow: https://apparmor.pujol.io/development/guidelines/#guidelines
func (p *Profile) Sort() {
	p.Rules.Sort()
}

// Format the profile for better readability before printing it.
// Follow: https://apparmor.pujol.io/development/guidelines/#the-file-block
func (p *Profile) Format() {
	const prefixOwner = "      "

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
				p.Rules = append(p.Rules[:i], append([]Rule{&RuleBase{}}, p.Rules[i:]...)...)
			}
		}
	}
}

// AddRule adds a new rule to the profile from a log map.
func (p *Profile) AddRule(log map[string]string) {

	// Generate profile flags and extra rules
	switch log["error"] {
	case "-2":
		if !slices.Contains(p.Flags, "mediate_deleted") {
			p.Flags = append(p.Flags, "mediate_deleted")
		}
	case "-13":
		if strings.Contains(log["info"], "namespace creation restricted") {
			p.Rules = append(p.Rules, newUsernsFromLog(log))
		} else if strings.Contains(log["info"], "disconnected path") && !slices.Contains(p.Flags, "attach_disconnected") {
			p.Flags = append(p.Flags, "attach_disconnected")
		}
	default:
	}

	switch log["class"] {
	case "rlimits":
		p.Rules = append(p.Rules, newRlimitFromLog(log))
	case "cap":
		p.Rules = append(p.Rules, newCapabilityFromLog(log))
	case "net":
		if log["family"] == "unix" {
			p.Rules = append(p.Rules, newUnixFromLog(log))
		} else {
			p.Rules = append(p.Rules, newNetworkFromLog(log))
		}
	case "io_uring":
		p.Rules = append(p.Rules, newIOUringFromLog(log))
	case "mount":
		if strings.Contains(log["flags"], "remount") {
			p.Rules = append(p.Rules, newRemountFromLog(log))
		} else {
			switch log["operation"] {
			case "mount":
				p.Rules = append(p.Rules, newMountFromLog(log))
			case "umount":
				p.Rules = append(p.Rules, newUmountFromLog(log))
			case "remount":
				p.Rules = append(p.Rules, newRemountFromLog(log))
			case "pivotroot":
				p.Rules = append(p.Rules, newPivotRootFromLog(log))
			}
		}
	case "posix_mqueue", "sysv_mqueue":
		p.Rules = append(p.Rules, newMqueueFromLog(log))
	case "signal":
		p.Rules = append(p.Rules, newSignalFromLog(log))
	case "ptrace":
		p.Rules = append(p.Rules, newPtraceFromLog(log))
	case "namespace":
		p.Rules = append(p.Rules, newUsernsFromLog(log))
	case "unix":
		p.Rules = append(p.Rules, newUnixFromLog(log))
	case "dbus":
		p.Rules = append(p.Rules, newDbusFromLog(log))
	case "file":
		if log["operation"] == "change_onexec" {
			p.Rules = append(p.Rules, newChangeProfileFromLog(log))
		} else {
			p.Rules = append(p.Rules, newFileFromLog(log))
		}
	default:
		if strings.Contains(log["operation"], "dbus") {
			p.Rules = append(p.Rules, newDbusFromLog(log))
		} else if log["family"] == "unix" {
			p.Rules = append(p.Rules, newUnixFromLog(log))
		}
	}
}
