// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

type File struct {
	Qualifier
	Path   string
	Access string
	Target string
}

func FileFromLog(log map[string]string) ApparmorRule {
	return &File{
		Qualifier: NewQualifierFromLog(log),
		Path:      log["name"],
		Access:    toAccess(log["requested_mask"]),
		Target:    log["target"],
	}
}

func (r *File) Less(other any) bool {
	o, _ := other.(*File)
	letterR := getLetterIn(fileAlphabet, r.Path)
	letterO := getLetterIn(fileAlphabet, o.Path)
	if fileWeights[letterR] == fileWeights[letterO] || letterR == "" || letterO == "" {
		if r.Qualifier.Equals(o.Qualifier) {
			if r.Path == o.Path {
				if r.Access == o.Access {
					return r.Target < o.Target
				}
				return r.Access < o.Access
			}
			return r.Path < o.Path
		}
		return r.Qualifier.Less(o.Qualifier)
	}
	return fileWeights[letterR] < fileWeights[letterO]
}

func (r *File) Equals(other any) bool {
	o, _ := other.(*File)
	return r.Path == o.Path && r.Access == o.Access &&
		r.Target == o.Target && r.Qualifier.Equals(o.Qualifier)
}
