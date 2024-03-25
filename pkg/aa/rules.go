// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"strings"
)

type Rule struct {
	Comment     string
	NoNewPrivs  bool
	FileInherit bool
}

func (r *Rule) Less(other any) bool {
	return false
}

func (r *Rule) Equals(other any) bool {
	return false
}

// Qualifier to apply extra settings to a rule
type Qualifier struct {
	Audit       bool
	AccessType  string
	Owner       bool
	NoNewPrivs  bool
	FileInherit bool
	Optional    bool
	Comment     string
	Prefix      string
	Padding     string
}

func NewQualifierFromLog(log map[string]string) Qualifier {
	owner := false
	fsuid, hasFsUID := log["fsuid"]
	ouid, hasOuUID := log["ouid"]
	isDbus := strings.Contains(log["operation"], "dbus")
	if hasFsUID && hasOuUID && fsuid == ouid && ouid != "0" && !isDbus {
		owner = true
	}

	audit := false
	if log["apparmor"] == "AUDIT" {
		audit = true
	}

	fileInherit := false
	if log["operation"] == "file_inherit" {
		fileInherit = true
	}

	noNewPrivs := false
	optional := false
	msg := ""
	switch log["error"] {
	case "-1":
		if strings.Contains(log["info"], "optional:") {
			optional = true
			msg = strings.Replace(log["info"], "optional: ", "", 1)
		} else {
			noNewPrivs = true
		}
	case "-13":
		ignoreProfileInfo := []string{"namespace", "disconnected path"}
		for _, info := range ignoreProfileInfo {
			if strings.Contains(log["info"], info) {
				break
			}
		}
		msg = log["info"]
	default:
	}

	return Qualifier{
		Audit:       audit,
		Owner:       owner,
		NoNewPrivs:  noNewPrivs,
		FileInherit: fileInherit,
		Optional:    optional,
		Comment:     msg,
	}
}

func (r Qualifier) Less(other Qualifier) bool {
	if r.Owner == other.Owner {
		if r.Audit == other.Audit {
			return r.AccessType < other.AccessType
		}
		return r.Audit
	}
	return other.Owner
}

func (r Qualifier) Equals(other Qualifier) bool {
	return r.Audit == other.Audit && r.AccessType == other.AccessType &&
		r.Owner == other.Owner && r.NoNewPrivs == other.NoNewPrivs &&
		r.FileInherit == other.FileInherit
}

// Preamble specific rules

type Abi struct {
	Path    string
	IsMagic bool
}

func (r Abi) Less(other Abi) bool {
	if r.Path == other.Path {
		return r.IsMagic == other.IsMagic
	}
	return r.Path < other.Path
}

func (r Abi) Equals(other Abi) bool {
	return r.Path == other.Path && r.IsMagic == other.IsMagic
}

type Alias struct {
	Path          string
	RewrittenPath string
}

func (r Alias) Less(other Alias) bool {
	if r.Path == other.Path {
		return r.RewrittenPath < other.RewrittenPath
	}
	return r.Path < other.Path
}

func (r Alias) Equals(other Alias) bool {
	return r.Path == other.Path && r.RewrittenPath == other.RewrittenPath
}
