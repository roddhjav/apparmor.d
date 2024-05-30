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

func newRule(rule []string) RuleBase {
	comment := ""
	fileInherit, noNewPrivs, optional := false, false, false

	idx := 0
	for idx < len(rule) {
		if rule[idx] == COMMENT.Tok() {
			comment = " " + strings.Join(rule[idx+1:], " ")
			break
		}
		idx++
	}
	switch {
	case strings.Contains(comment, "file_inherit"):
		fileInherit = true
		comment = strings.Replace(comment, "file_inherit ", "", 1)
	case strings.HasPrefix(comment, "no new privs"):
		noNewPrivs = true
		comment = strings.Replace(comment, "no new privs ", "", 1)
	case strings.Contains(comment, "optional:"):
		optional = true
		comment = strings.Replace(comment, "optional: ", "", 1)
	}
	return RuleBase{
		Comment:     comment,
		NoNewPrivs:  noNewPrivs,
		FileInherit: fileInherit,
		Optional:    optional,
	}
}

func newRuleFromLog(log map[string]string) RuleBase {
	comment := ""
	fileInherit, noNewPrivs, optional := false, false, false

	if log["operation"] == "file_inherit" {
		fileInherit = true
	}
	if log["error"] == "-1" {
		if strings.Contains(log["info"], "optional:") {
			optional = true
			comment = strings.Replace(log["info"], "optional: ", "", 1)
		} else {
			noNewPrivs = true
		}
	}
	if log["info"] != "" {
		comment += " " + log["info"]
	}
	return RuleBase{
		IsLineRule:  false,
		Comment:     comment,
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
	return renderTemplate(r.Kind(), r)
}

func (r RuleBase) Constraint() constraint {
	return anyKind
}

func (r RuleBase) Kind() Kind {
	return COMMENT
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
