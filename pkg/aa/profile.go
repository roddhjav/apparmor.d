// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"fmt"
	"slices"
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
			"attach_disconneced.path=", "attach_disconnected", "audit",
			"chroot_relative", "complain", "debug", "default_allow", "enforce",
			"interruptible", "kill", "mediate_deleted",
			"prompt", "unconfined", "namespace_relative", "delegate_deleted", "chroot_attach",
			"chroot_no_attach", "no_attach_disconnected",
		},
	}
	conflicts[PROFILE] = map[string][][]string{
		tokFLAGS: {
			{"enforce", "complain"},
			{"enforce", "unconfined"},
			{"enforce", "prompt"},
			{"complain", "unconfined"},
			{"default_allow", "kill"},
			{"default_allow", "enforce"},
		},
	}
}

// Profile represents a single AppArmor profile.
type Profile struct {
	Base
	Header
	Rules Rules
}

// Header represents the header of a profile.
type Header struct {
	Name        string
	NameSpace   string
	Attachments []string
	Attributes  map[string]string
	Flags       []string
}

func newHeader(rule rule) (Header, error) {
	if len(rule) == 0 {
		return Header{}, nil
	}
	if rule.Get(0) == PROFILE.Tok() {
		rule = rule[1:]
	}
	name, attachments := "", []string{}
	if len(rule) >= 1 {
		name = rule.Get(0)
		if len(rule) > 1 {
			attachments = rule[1:].GetSlice()
		}
	}
	attributes := make(map[string]string)
	for k, v := range rule.GetValues(tokATTRIBUTES).GetAsMap() {
		attributes[k] = strings.Join(v, "")
	}
	return Header{
		Name:        name,
		Attachments: attachments,
		Attributes:  attributes,
		Flags:       rule.GetValuesAsSlice(tokFLAGS),
	}, rule.ValidateMapKeys([]string{tokATTRIBUTES, tokFLAGS})
}

func (p *Profile) Kind() Kind {
	return PROFILE
}

func (p *Profile) Constraint() Constraint {
	return BlockRule
}

func (p *Profile) String() string {
	return renderTemplate(p.Kind(), p)
}

func (p *Profile) Validate() error {
	if err := validateValues(p.Kind(), tokFLAGS, p.Flags); err != nil {
		return fmt.Errorf("profile %s: %w", p.Name, err)
	}
	if err := validateConflicts(p.Kind(), tokFLAGS, p.Flags); err != nil {
		return fmt.Errorf("profile %s: %w", p.Name, err)
	}
	return p.Rules.Validate()
}

func (p *Profile) Compare(other Rule) int {
	o, _ := other.(*Profile)
	if res := compare(p.Name, o.Name); res != 0 {
		return res
	}
	return compare(p.Attachments, o.Attachments)
}

func (p *Profile) Merge(other Rule) bool {
	slices.Sort(p.Flags)
	p.Flags = slices.Compact(p.Flags)
	p.Rules = p.Rules.Merge()
	return false
}

func (p *Profile) Lengths() []int {
	return []int{} // No len for profile
}

func (p *Profile) setPaddings(max []int) {} // No paddings for profile

func (p *Profile) Sort() {
	p.Rules = p.Rules.Sort()
}

func (p *Profile) Format() {
	p.Rules = p.Rules.Format()
}

// GetAttachments return a nested attachment string
func (p *Profile) GetAttachments() string {
	switch len(p.Attachments) {
	case 0:
		return ""
	case 1:
		return p.Attachments[0]
	default:
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
		// class
		"rlimits":   newRlimitFromLog,
		"namespace": newUsernsFromLog,
		"cap":       newCapabilityFromLog,
		"net": func(log map[string]string) Rule {
			if log["family"] == "unix" {
				return newUnixFromLog(log)
			} else {
				return newNetworkFromLog(log)
			}
		},
		"posix_mqueue": newMqueueFromLog,
		"sysv_mqueue":  newMqueueFromLog,
		"signal":       newSignalFromLog,
		"ptrace":       newPtraceFromLog,
		"unix":         newUnixFromLog,
		"io_uring":     newIOUringFromLog,
		"dbus":         newDbusFromLog,
		"mount": func(log map[string]string) Rule {
			if strings.Contains(log["flags"], "remount") {
				return newRemountFromLog(log)
			}
			newRule := newLogMountMap[log["operation"]]
			return newRule(log)
		},
		"file": func(log map[string]string) Rule {
			if log["operation"] == "change_onexec" {
				return newChangeProfileFromLog(log)
			} else {
				return newFileFromLog(log)
			}
		},
		// operation
		"capable":     newCapabilityFromLog,
		"chmod":       newFileFromLog,
		"chown":       newFileFromLog,
		"exec":        newFileFromLog,
		"getattr":     newFileFromLog,
		"link":        newFileFromLog,
		"mkdir":       newFileFromLog,
		"mknod":       newFileFromLog,
		"open":        newFileFromLog,
		"rename_dest": newFileFromLog,
		"rename_src":  newFileFromLog,
		"rmdir":       newFileFromLog,
		"truncate":    newFileFromLog,
		"unlink":      newFileFromLog,
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

	done := false
	for _, key := range []string{"class", "family", "operation"} {
		if newRule, ok := newLogMap[log[key]]; ok {
			p.Rules = append(p.Rules, newRule(log))
			done = true
			break
		}
	}

	if !done {
		switch {
		case strings.HasPrefix(log["operation"], "file_"):
			p.Rules = append(p.Rules, newFileFromLog(log))
		case strings.Contains(log["operation"], "dbus"):
			p.Rules = append(p.Rules, newDbusFromLog(log))
		default:
			fmt.Printf("unknown log type: %s:%v\n", log["operation"], log)
		}
	}
}
