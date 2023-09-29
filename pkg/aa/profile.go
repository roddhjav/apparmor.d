// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2023 Alexandre Pujol <alexandre@pujol.io>
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
	noNewPrivs := false
	fileInherit := false
	if log["operation"] == "file_inherit" {
		fileInherit = true
	}

	switch log["error"] {
	case "-1":
		noNewPrivs = true
	case "-2":
		if !slices.Contains(p.Flags, "mediate_deleted") {
			p.Flags = append(p.Flags, "mediate_deleted")
		}
	case "-13":
		// FIXME: -13 can be a lot of things, not only attach_disconnected
		// Eg: info="User namespace creation restricted"
		if !slices.Contains(p.Flags, "attach_disconnected") {
			p.Flags = append(p.Flags, "attach_disconnected")
		}
	default:
	}

	switch log["class"] {
	case "cap":
		p.Rules = append(p.Rules, CapabilityFromLog(log, noNewPrivs, fileInherit))
	case "net":
		p.Rules = append(p.Rules, NetworkFromLog(log, noNewPrivs, fileInherit))
	case "mount":
		p.Rules = append(p.Rules, MountFromLog(log, noNewPrivs, fileInherit))
	case "remount":
		p.Rules = append(p.Rules, RemountFromLog(log, noNewPrivs, fileInherit))
	case "umount":
		p.Rules = append(p.Rules, UmountFromLog(log, noNewPrivs, fileInherit))
	case "pivot_root":
		p.Rules = append(p.Rules, PivotRootFromLog(log, noNewPrivs, fileInherit))
	case "change_profile":
		p.Rules = append(p.Rules, RemountFromLog(log, noNewPrivs, fileInherit))
	case "mqueue":
		p.Rules = append(p.Rules, MqueueFromLog(log, noNewPrivs, fileInherit))
	case "signal":
		p.Rules = append(p.Rules, SignalFromLog(log, noNewPrivs, fileInherit))
	case "ptrace":
		p.Rules = append(p.Rules, PtraceFromLog(log, noNewPrivs, fileInherit))
	case "namespace":
		p.Rules = append(p.Rules, UsernsFromLog(log, noNewPrivs, fileInherit))
	case "unix":
		p.Rules = append(p.Rules, UnixFromLog(log, noNewPrivs, fileInherit))
	case "file":
		p.Rules = append(p.Rules, FileFromLog(log, noNewPrivs, fileInherit))
	default:
		if strings.Contains(log["operation"], "dbus") {
			p.Rules = append(p.Rules, DbusFromLog(log, noNewPrivs, fileInherit))
		} else if log["family"] == "unix" {
			p.Rules = append(p.Rules, UnixFromLog(log, noNewPrivs, fileInherit))
		}
	}
}

func typeToValue(i reflect.Type) string {
	return strings.ToLower(strings.TrimPrefix(i.String(), "*aa."))
}

// Sort the rules in the profile
func (p *AppArmorProfile) Sort() {
	sort.Slice(p.Rules, func(i, j int) bool {
		typeOfI := reflect.TypeOf(p.Rules[i])
		typeOfJ := reflect.TypeOf(p.Rules[j])
		if typeOfI != typeOfJ {
			valueOfI := typeToValue(typeOfI)
			valueOfJ := typeToValue(typeOfJ)
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
