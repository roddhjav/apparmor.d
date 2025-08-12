// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"strings"
)

type Base struct {
	Comment     string
	IsLineRule  bool
	NoNewPrivs  bool
	FileInherit bool
	Optional    bool
	Paddings    []string
}

func newBase(rule rule) Base {
	comment := ""
	fileInherit, noNewPrivs, optional := false, false, false

	if len(rule) > 0 {
		if len(rule.Get(0)) > 0 && rule.Get(0)[0] == '#' {
			// Line rule is a comment
			rule = rule[1:]
			comment = rule.GetString()
		} else {
			// Comma rule, with comment at the end
			comment = rule[len(rule)-1].comment
		}
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
	return Base{
		Comment:     comment,
		NoNewPrivs:  noNewPrivs,
		FileInherit: fileInherit,
		Optional:    optional,
	}
}

func newBaseFromLog(log map[string]string) Base {
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
	return Base{
		IsLineRule:  false,
		Comment:     comment,
		NoNewPrivs:  noNewPrivs,
		FileInherit: fileInherit,
		Optional:    optional,
	}
}

func (r Base) Padding(i int) string {
	if i >= len(r.Paddings) {
		return ""
	}
	return r.Paddings[i]
}

func (r *Base) merge(other Base) bool {
	r.NoNewPrivs = r.NoNewPrivs || other.NoNewPrivs
	r.FileInherit = r.FileInherit || other.FileInherit
	r.Optional = r.Optional || other.Optional
	if other.Comment != "" {
		r.Comment += " " + other.Comment
	}
	return true
}

func (r Base) addLine(other Rule) bool {
	return false
}

type Qualifier struct {
	Audit      bool
	AccessType string
}

func newQualifierFromLog(log map[string]string) Qualifier {
	audit := log["apparmor"] == "AUDIT"
	return Qualifier{Audit: audit}
}

func (r Qualifier) Compare(o Qualifier) int {
	if r := compare(r.Audit, o.Audit); r != 0 {
		return r
	}
	return compare(r.AccessType, o.AccessType)
}

func (r Qualifier) Equal(o Qualifier) bool {
	return r.Audit == o.Audit && r.AccessType == o.AccessType
}

func (r Qualifier) getLenAudit() int {
	return length("audit", r.Audit)
}

func (r Qualifier) getLenAccess() int {
	lenAccess := 0
	if r.AccessType != "" {
		lenAccess = length("", r.AccessType)
	}
	return lenAccess
}

func (r Qualifier) setPaddings(max []int) []string {
	return setPaddings(max,
		[]string{"audit", ""},
		[]any{r.Audit, r.AccessType},
	)
}
