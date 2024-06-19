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
	LINK      Kind = "link"
	FILE      Kind = "file"
	tokOWNER       = "owner"
	tokSUBSET      = "subset"
)

func init() {
	requirements[FILE] = requirement{
		"access": {"m", "r", "w", "l", "k"},
		"transition": {
			"ix", "ux", "Ux", "px", "Px", "cx", "Cx", "pix", "Pix", "cix",
			"Cix", "pux", "PUx", "cux", "CUx", "x",
		},
	}
}

func isOwner(log map[string]string) bool {
	fsuid, hasFsUID := log["fsuid"]
	ouid, hasOuUID := log["ouid"]
	isDbus := strings.Contains(log["operation"], "dbus")
	if hasFsUID && hasOuUID && fsuid == ouid && ouid != "0" && !isDbus {
		return true
	}
	return false
}

type File struct {
	RuleBase
	Qualifier
	Owner  bool
	Path   string
	Access []string
	Target string
}

func newFileFromLog(log map[string]string) Rule {
	accesses, err := toAccess("file-log", log["requested_mask"])
	if err != nil {
		panic(fmt.Errorf("newFileFromLog(%v): %w", log, err))
	}
	if slices.Compare(accesses, []string{"l"}) == 0 {
		return newLinkFromLog(log)
	}
	return &File{
		RuleBase:  newRuleFromLog(log),
		Qualifier: newQualifierFromLog(log),
		Owner:     isOwner(log),
		Path:      log["name"],
		Access:    accesses,
		Target:    log["target"],
	}
}

func (r *File) Validate() error {
	return nil
}

func (r *File) Compare(other Rule) int {
	o, _ := other.(*File)

	letterR := getLetterIn(fileAlphabet, r.Path)
	letterO := getLetterIn(fileAlphabet, o.Path)
	if fileWeights[letterR] != fileWeights[letterO] && letterR != "" && letterO != "" {
		return fileWeights[letterR] - fileWeights[letterO]
	}
	if res := compare(r.Owner, o.Owner); res != 0 {
		return res
	}
	if res := compare(r.Path, o.Path); res != 0 {
		return res
	}
	if res := compare(r.Access, o.Access); res != 0 {
		return res
	}
	if res := compare(r.Target, o.Target); res != 0 {
		return res
	}
	return r.Qualifier.Compare(o.Qualifier)
}

func (r *File) String() string {
	return renderTemplate(r.Kind(), r)
}

func (r *File) Constraint() constraint {
	return blockKind
}

func (r *File) Kind() Kind {
	return FILE
}

type Link struct {
	RuleBase
	Qualifier
	Owner  bool
	Subset bool
	Path   string
	Target string
}

func newLinkFromLog(log map[string]string) Rule {
	return &Link{
		RuleBase:  newRuleFromLog(log),
		Qualifier: newQualifierFromLog(log),
		Owner:     isOwner(log),
		Path:      log["name"],
		Target:    log["target"],
	}
}

func (r *Link) Validate() error {
	return nil
}

func (r *Link) Compare(other Rule) int {
	o, _ := other.(*Link)

	if res := compare(r.Owner, o.Owner); res != 0 {
		return res
	}
	if res := compare(r.Path, o.Path); res != 0 {
		return res
	}
	if res := compare(r.Target, o.Target); res != 0 {
		return res
	}
	if res := compare(r.Subset, o.Subset); res != 0 {
		return res
	}
	return r.Qualifier.Compare(o.Qualifier)
}

func (r *Link) String() string {
	return renderTemplate(r.Kind(), r)
}

func (r *Link) Constraint() constraint {
	return blockKind
}

func (r *Link) Kind() Kind {
	return LINK
}
