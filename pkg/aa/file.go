// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

type File struct {
	RuleBase
	Qualifier
	Owner  bool
	Path   string
	Access []string
	Target string
}

func newFileFromLog(log map[string]string) Rule {
	owner := false
	fsuid, hasFsUID := log["fsuid"]
	ouid, hasOuUID := log["ouid"]
	isDbus := strings.Contains(log["operation"], "dbus")
	if hasFsUID && hasOuUID && fsuid == ouid && ouid != "0" && !isDbus {
		owner = true
	}
	return &File{
		RuleBase:  newRuleFromLog(log),
		Qualifier: newQualifierFromLog(log),
		Owner:     owner,
		Path:      log["name"],
		Access:    toAccess("file-log", log["requested_mask"]),
		Target:    log["target"],
	}
}

func (r *File) Less(other any) bool {
	o, _ := other.(*File)
	letterR := getLetterIn(fileAlphabet, r.Path)
	letterO := getLetterIn(fileAlphabet, o.Path)
	if fileWeights[letterR] != fileWeights[letterO] && letterR != "" && letterO != "" {
		return fileWeights[letterR] < fileWeights[letterO]
	}
	if r.Path != o.Path {
		return r.Path < o.Path
	}
	if len(r.Access) != len(o.Access) {
		return len(r.Access) < len(o.Access)
	}
	if r.Target != o.Target {
		return r.Target < o.Target
	}
	if o.Owner != r.Owner {
		return r.Owner
	}
	return r.Qualifier.Less(o.Qualifier)
}

func (r *File) Equals(other any) bool {
	o, _ := other.(*File)
	return r.Path == o.Path && slices.Equal(r.Access, o.Access) && r.Owner == o.Owner &&
		r.Target == o.Target && r.Qualifier.Equals(o.Qualifier)
}

func (r *File) String() string {
	return renderTemplate(r.Kind(), r)
}

func (r *File) Constraint() constraint {
	return blockKind
}

func (r *File) Kind() string {
	return "file"
}
