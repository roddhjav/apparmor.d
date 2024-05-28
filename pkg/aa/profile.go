// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"fmt"
	"maps"
	"reflect"
	"slices"
	"sort"
	"strings"
)

const (
	PROFILE Kind = "profile"

	tokATTRIBUTES = "xattrs"
	tokFLAGS      = "flags"
)

func init() {
	requirements[PROFILE] = requirement{
		tokFLAGS: {
			"enforce", "complain", "kill", "default_allow", "unconfined",
			"prompt", "audit", "mediate_deleted", "attach_disconnected",
			"attach_disconneced.path=", "chroot_relative", "debug",
			"interruptible", "kill", "kill.signal=",
		},
	}
}

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

func newHeader(rule []string) (Header, error) {
	if len(rule) == 0 {
		return Header{}, nil
	}
	if rule[len(rule)-1] == "{" {
		rule = rule[:len(rule)-1]
	}
	if rule[0] == PROFILE.Tok() {
		rule = rule[1:]
	}

	delete := []int{}
	flags := []string{}
	attributes := make(map[string]string)
	for idx, token := range rule {
		if item, ok := strings.CutPrefix(token, tokFLAGS+"="); ok {
			flags = tokenToSlice(item)
			delete = append(delete, idx)
		} else if item, ok := strings.CutPrefix(token, tokATTRIBUTES+"="); ok {
			for _, m := range tokenToSlice(item) {
				kv := strings.SplitN(m, "=", 2)
				attributes[kv[0]] = kv[1]
			}
			delete = append(delete, idx)
		}
	}
	for i := len(delete) - 1; i >= 0; i-- {
		rule = slices.Delete(rule, delete[i], delete[i]+1)
	}

	name, attachments := "", []string{}
	if len(rule) >= 1 {
		name = rule[0]
		if len(rule) > 1 {
			attachments = rule[1:]
		}
	}
	return Header{
		Name:        name,
		Attachments: attachments,
		Attributes:  attributes,
		Flags:       flags,
	}, nil
}

func (r *Profile) Validate() error {
	if err := validateValues(r.Kind(), tokFLAGS, r.Flags); err != nil {
		return fmt.Errorf("profile %s: %w", r.Name, err)
	}
	return r.Rules.Validate()
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
	return renderTemplate(p.Kind(), p)
}

func (p *Profile) Constraint() constraint {
	return blockKind
}

func (p *Profile) Kind() Kind {
	return PROFILE
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
				continue
			}

			// File rule
			if typeOfI == reflect.TypeFor[*File]() && typeOfJ == reflect.TypeFor[*File]() {
				// Merge access
				fileI := p.Rules[i].(*File)
				fileJ := p.Rules[j].(*File)
				if fileI.Path == fileJ.Path {
					fileI.Access = append(fileI.Access, fileJ.Access...)
					slices.SortFunc(fileI.Access, cmpFileAccess)
					fileI.Access = slices.Compact(fileI.Access)
					p.Rules = append(p.Rules[:j], p.Rules[j+1:]...)
					j--
				}
			}
		}
	}
}

// Sort the rules in a profile.
// Follow: https://apparmor.pujol.io/development/guidelines/#guidelines
func (p *Profile) Sort() {
	sort.Slice(p.Rules, func(i, j int) bool {
		typeOfI := reflect.TypeOf(p.Rules[i])
		typeOfJ := reflect.TypeOf(p.Rules[j])
		if typeOfI != typeOfJ {
			valueOfI := typeToValue(typeOfI)
			valueOfJ := typeToValue(typeOfJ)
			if typeOfI == reflect.TypeFor[*Include]() && p.Rules[i].(*Include).IfExists {
				valueOfI = "include_if_exists"
			}
			if typeOfJ == reflect.TypeFor[*Include]() && p.Rules[j].(*Include).IfExists {
				valueOfJ = "include_if_exists"
			}
			return ruleWeights[valueOfI] < ruleWeights[valueOfJ]
		}
		return p.Rules[i].Less(p.Rules[j])
	})
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
		if typeOfI == reflect.TypeFor[*File]() && typeOfJ == reflect.TypeFor[*File]() {
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
				p.Rules = append(p.Rules[:i], append(Rules{nil}, p.Rules[i:]...)...)
			}
		}
	}
}

// GetAttachments return a nested attachment string
func (p *Profile) GetAttachments() string {
	if len(p.Attachments) == 0 {
		return ""
	} else if len(p.Attachments) == 1 {
		return p.Attachments[0]
	} else {
		res := []string{}
		for _, attachment := range p.Attachments {
			if strings.HasPrefix(attachment, "/") {
				res = append(res, attachment[1:])
			} else {
				res = append(res, attachment)
			}
		}
		return "/{" + strings.Join(res, ",") + "}"
	}
}

var (
	newLogMap = map[string]func(log map[string]string) Rule{
		"rlimits":      newRlimitFromLog,
		"cap":          newCapabilityFromLog,
		"io_uring":     newIOUringFromLog,
		"signal":       newSignalFromLog,
		"ptrace":       newPtraceFromLog,
		"namespace":    newUsernsFromLog,
		"unix":         newUnixFromLog,
		"dbus":         newDbusFromLog,
		"posix_mqueue": newMqueueFromLog,
		"sysv_mqueue":  newMqueueFromLog,
		"mount": func(log map[string]string) Rule {
			if strings.Contains(log["flags"], "remount") {
				return newRemountFromLog(log)
			}
			newRule := newLogMountMap[log["operation"]]
			return newRule(log)
		},
		"net": func(log map[string]string) Rule {
			if log["family"] == "unix" {
				return newUnixFromLog(log)
			} else {
				return newNetworkFromLog(log)
			}
		},
		"file": func(log map[string]string) Rule {
			if log["operation"] == "change_onexec" {
				return newChangeProfileFromLog(log)
			} else {
				return newFileFromLog(log)
			}
		},
	}
	newLogMountMap = map[string]func(log map[string]string) Rule{
		"mount":     newMountFromLog,
		"umount":    newUmountFromLog,
		"remount":   newRemountFromLog,
		"pivotroot": newPivotRootFromLog,
	}
)

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

	if newRule, ok := newLogMap[log["class"]]; ok {
		p.Rules = append(p.Rules, newRule(log))
	} else {
		if strings.Contains(log["operation"], "dbus") {
			p.Rules = append(p.Rules, newDbusFromLog(log))
		} else if log["family"] == "unix" {
			p.Rules = append(p.Rules, newUnixFromLog(log))
		} else {
			panic("unknown class: " + log["class"])
		}
	}
}
