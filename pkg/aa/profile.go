// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"maps"
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
