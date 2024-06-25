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
	Base
	Qualifier
	Owner  bool
	Path   string
	Access []string
	Target string
}

func newFile(q Qualifier, rule rule) (Rule, error) {
	path, access, target, owner := "", "", "", false
	if len(rule) > 0 {
		if rule.Get(0) == tokOWNER {
			owner = true
			rule = rule[1:]
		}
		if rule.Get(0) == FILE.Tok() {
			rule = rule[1:]
		}

		r := rule.GetSlice()
		size := len(r)
		if size < 2 {
			return nil, fmt.Errorf("missing file or access in rule: %s", rule)
		}

		path, access = r[0], r[1]
		if size > 2 {
			if r[2] != tokARROW {
				return nil, fmt.Errorf("missing '%s' in rule: %s", tokARROW, rule)
			}
			target = r[3]
		}
	}
	accesses, err := toAccess(FILE, access)
	if err != nil {
		return nil, err
	}
	return &File{
		Base:      newBase(rule),
		Qualifier: q,
		Owner:     owner,
		Path:      path,
		Access:    accesses,
		Target:    target,
	}, nil
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
		Base:      newBaseFromLog(log),
		Qualifier: newQualifierFromLog(log),
		Owner:     isOwner(log),
		Path:      log["name"],
		Access:    accesses,
		Target:    log["target"],
	}
}

func (r *File) Kind() Kind {
	return FILE
}

func (r *File) Constraint() constraint {
	return blockKind
}

func (r *File) String() string {
	return renderTemplate(r.Kind(), r)
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

func (r *File) Merge(other Rule) bool {
	o, _ := other.(*File)

	if !r.Qualifier.Equal(o.Qualifier) {
		return false
	}
	if r.Owner == o.Owner && r.Path == o.Path && r.Target == o.Target {
		r.Access = merge(r.Kind(), "access", r.Access, o.Access)
		b := &r.Base
		return b.merge(o.Base)
	}
	return false
}

type Link struct {
	Base
	Qualifier
	Owner  bool
	Subset bool
	Path   string
	Target string
}

func newLink(q Qualifier, rule rule) (Rule, error) {
	owner, subset, path, target := false, false, "", ""
	if len(rule) > 0 {
		if rule.Get(0) == tokOWNER {
			owner = true
			rule = rule[1:]
		}
		if len(rule) > 0 && rule.Get(0) == tokSUBSET {
			subset = true
			rule = rule[1:]
		}

		r := rule.GetSlice()
		size := len(r)
		if size > 0 {
			path = r[0]
		}
		if size > 2 {
			if r[1] != tokARROW {
				return nil, fmt.Errorf("missing '%s' in rule: %s", tokARROW, rule)
			}
			target = r[2]
		}
	}
	return &Link{
		Base:      newBase(rule),
		Qualifier: q,
		Owner:     owner,
		Subset:    subset,
		Path:      path,
		Target:    target,
	}, nil
}

func newLinkFromLog(log map[string]string) Rule {
	return &Link{
		Base:      newBaseFromLog(log),
		Qualifier: newQualifierFromLog(log),
		Owner:     isOwner(log),
		Path:      log["name"],
		Target:    log["target"],
	}
}

func (r *Link) Kind() Kind {
	return LINK
}

func (r *Link) Constraint() constraint {
	return blockKind
}

func (r *Link) String() string {
	return renderTemplate(r.Kind(), r)
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
