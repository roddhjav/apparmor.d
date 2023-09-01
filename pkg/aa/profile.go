// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2023 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"bytes"
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
		if !slices.Contains(p.Flags, "attach_disconnected") {
			p.Flags = append(p.Flags, "attach_disconnected")
		}
	default:
	}

	switch log["class"] {
	case "cap":
		p.Capability = append(p.Capability, NewCapability(log, noNewPrivs, fileInherit))
	case "file":
		p.File = append(p.File, NewFile(log, noNewPrivs, fileInherit))
	case "net":
		if log["family"] == "unix" {
			p.Unix = append(p.Unix, NewUnix(log, noNewPrivs, fileInherit))
		} else {
			p.Network = append(p.Network, NewNetwork(log, noNewPrivs, fileInherit))
		}
	case "signal":
		p.Signal = append(p.Signal, NewSignal(log, noNewPrivs, fileInherit))
	case "ptrace":
		p.Ptrace = append(p.Ptrace, NewPtrace(log, noNewPrivs, fileInherit))
	case "unix":
		p.Unix = append(p.Unix, NewUnix(log, noNewPrivs, fileInherit))
	case "mount":
		p.Mount = append(p.Mount, NewMount(log, noNewPrivs, fileInherit))
	default:
		if strings.Contains(log["operation"], "dbus") {
			p.Dbus = append(p.Dbus, NewDbus(log, noNewPrivs, fileInherit))
		} else if log["family"] == "unix" {
			p.Unix = append(p.Unix, NewUnix(log, noNewPrivs, fileInherit))
		}
	}
}

