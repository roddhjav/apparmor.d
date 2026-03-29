// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
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
			"attach_disconnected", "attach_disconnected.ipc", "attach_disconnected.ipc=",
			"attach_disconnected.path=",
			"audit",
			"chroot_attach", "chroot_no_attach", "chroot_relative",
			"complain", "debug",
			"default_allow", "delegate_deleted",
			"enforce", "error=",
			"interruptible", "kill", "kill.signal=",
			"mediate_deleted",
			"namespace_relative", "no_attach_disconnected",
			"prompt", "unconfined",
		},
	}
	conflicts[PROFILE] = map[string][][]string{
		tokFLAGS: {
			// Mode conflicts: enforce, complain, kill, unconfined, prompt are mutually exclusive
			{"enforce", "complain"},
			{"enforce", "kill"},
			{"enforce", "unconfined"},
			{"enforce", "prompt"},
			{"complain", "kill"},
			{"complain", "unconfined"},
			{"complain", "prompt"},
			{"kill", "unconfined"},
			{"kill", "prompt"},
			// Note: kill + interruptible is valid (flags_ok33, flags_ok36)
			{"unconfined", "prompt"},

			// default_allow conflicts with all modes except enforce
			{"default_allow", "complain"},
			{"default_allow", "kill"},
			{"default_allow", "unconfined"},
			{"default_allow", "prompt"},
			{"default_allow", "enforce"},

			// Namespace conflicts
			{"namespace_relative", "chroot_relative"},

			// Deletion conflicts
			{"mediate_deleted", "delegate_deleted"},

			// Disconnected conflicts
			{"attach_disconnected", "no_attach_disconnected"},
			{"attach_disconnected.ipc", "no_attach_disconnected"},
			{"attach_disconnected.ipc=", "no_attach_disconnected"},
			{"chroot_attach", "chroot_no_attach"},
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

	flags := rule.GetValuesAsSlice(tokFLAGS)
	for i, f := range flags {
		flags[i] = strings.TrimRight(f, ",")
	}
	return Header{
		Name:        name,
		Attachments: attachments,
		Attributes:  attributes,
		Flags:       flags,
	}, rule.ValidateMapKeys([]string{tokATTRIBUTES, tokFLAGS, "identities"})
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
	if err := validateProfileFlags(p.Flags); err != nil {
		return fmt.Errorf("profile %s: %w", p.Name, err)
	}
	return p.Rules.Validate()
}

// validateProfileFlags performs additional validation on profile flags
// beyond simple value/conflict checks.
func validateProfileFlags(flags []string) error {
	for _, f := range flags {
		switch {
		case strings.HasPrefix(f, "kill.signal="):
			sig := strings.TrimPrefix(f, "kill.signal=")
			if sig == "" || strings.ContainsAny(sig, ".=/") {
				return fmt.Errorf("invalid kill.signal value '%s'", sig)
			}
			// Signal name must be a valid identifier (letters, digits, +)
			// e.g., hup, kill, usr1, rtmin+0
			if !regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9+]*$`).MatchString(sig) {
				// Also allow pure numeric signals
				if _, err := strconv.Atoi(sig); err != nil {
					return fmt.Errorf("invalid kill.signal value '%s'", sig)
				}
			}
		case strings.HasPrefix(f, "error="):
			val := strings.TrimPrefix(f, "error=")
			if val == "" {
				return fmt.Errorf("invalid error value: empty")
			}
			// error= can be a number or an errno name (e.g., ENOENT, EISCONN)
			if _, err := strconv.Atoi(val); err != nil {
				// Must be a valid errno name: uppercase letters only
				if !regexp.MustCompile(`^E[A-Z]+$`).MatchString(val) {
					return fmt.Errorf("invalid error value '%s'", val)
				}
			}
		case strings.HasPrefix(f, "attach_disconnected.path="):
			path := strings.TrimPrefix(f, "attach_disconnected.path=")
			if path == "" || !strings.HasPrefix(path, "/") {
				return fmt.Errorf("invalid attach_disconnected.path value '%s': must be absolute", path)
			}
		case strings.HasPrefix(f, "attach_disconnected.ipc="):
			path := strings.TrimPrefix(f, "attach_disconnected.ipc=")
			if path == "" || !strings.HasPrefix(path, "/") {
				return fmt.Errorf("invalid attach_disconnected.ipc value '%s': must be absolute", path)
			}
		}
	}
	// Check for duplicate attach_disconnected.ipc= values
	ipcCount := 0
	for _, f := range flags {
		if strings.HasPrefix(f, "attach_disconnected.ipc=") {
			ipcCount++
		}
	}
	if ipcCount > 1 {
		return fmt.Errorf("duplicate attach_disconnected.ipc flags")
	}
	return nil
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
