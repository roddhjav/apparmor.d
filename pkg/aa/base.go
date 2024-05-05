// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"strings"
)

type RuleBase struct {
	IsLineRule  bool
	Comment     string
	NoNewPrivs  bool
	FileInherit bool
	Prefix      string
	Padding     string
	Optional    bool
}


func newRuleFromLog(log map[string]string) RuleBase {
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

	return RuleBase{
		IsLineRule:  false,
		Comment:     msg,
		NoNewPrivs:  noNewPrivs,
		FileInherit: fileInherit,
		Optional:    optional,
	}
}

func (r RuleBase) Less(other any) bool {
	return false
}

func (r RuleBase) Equals(other any) bool {
	return false
}

func (r RuleBase) String() string {
	return renderTemplate("comment", r)
}

func (r RuleBase) Constraint() constraint {
	return anyKind
}

func (r RuleBase) Kind() string {
	return "base"
}

type Qualifier struct {
	Audit      bool
	AccessType string
}

func newQualifierFromLog(log map[string]string) Qualifier {
	audit := false
	if log["apparmor"] == "AUDIT" {
		audit = true
	}
	return Qualifier{Audit: audit}
}

func (r Qualifier) Less(other Qualifier) bool {
	if r.Audit != other.Audit {
		return r.Audit
	}
	return r.AccessType < other.AccessType
}

func (r Qualifier) Equals(other Qualifier) bool {
	return r.Audit == other.Audit && r.AccessType == other.AccessType
}
