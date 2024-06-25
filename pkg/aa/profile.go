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
			"enforce", "complain", "kill", "default_allow", "unconfined",
			"prompt", "audit", "mediate_deleted", "attach_disconnected",
			"attach_disconneced.path=", "chroot_relative", "debug",
			"interruptible", "kill", "kill.signal=",
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
	}, nil
}

func (p *Profile) Kind() Kind {
	return PROFILE
}

func (p *Profile) Constraint() constraint {
	return blockKind
}

func (p *Profile) String() string {
	return renderTemplate(p.Kind(), p)
}

func (r *Profile) Validate() error {
	if err := validateValues(r.Kind(), tokFLAGS, r.Flags); err != nil {
		return fmt.Errorf("profile %s: %w", r.Name, err)
	}
	return r.Rules.Validate()
}

func (r *Profile) Compare(other Rule) int {
	o, _ := other.(*Profile)
	if res := compare(r.Name, o.Name); res != 0 {
		return res
	}
	return compare(r.Attachments, o.Attachments)
}

func (p *Profile) Merge(other Rule) bool {
	slices.Sort(p.Flags)
	p.Flags = slices.Compact(p.Flags)
	p.Rules = p.Rules.Merge()
	return false
}

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
		"net": newNetworkFromLog,
		"file": func(log map[string]string) Rule {
			if log["operation"] == "change_onexec" {
				return newChangeProfileFromLog(log)
			} else {
				return newFileFromLog(log)
			}
		},
		"exec":       newFileFromLog,
		"getattr":    newFileFromLog,
		"mkdir":      newFileFromLog,
		"mknod":      newFileFromLog,
		"open":       newFileFromLog,
		"rename_src": newFileFromLog,
		"truncate":   newFileFromLog,
		"unlink":     newFileFromLog,
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
			fmt.Printf("unknown log type: %s", log["operation"])
		}
	}
}
