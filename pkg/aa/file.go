// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2023 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"strings"
)

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
		Access:    maskToAccess[log["requested_mask"]],
		Target:    log["target"],
	}
}

func (r *File) Less(other any) bool {
	o, _ := other.(*File)
	letterR := ""
	letterO := ""
	for _, letter := range fileAlphabet {
		if strings.HasPrefix(r.Path, letter) {
			letterR = letter
		}
		if strings.HasPrefix(o.Path, letter) {
			letterO = letter
		}
	}

	if fileWeights[letterR] == fileWeights[letterO] || letterR == "" || letterO == "" {
		if r.Path == o.Path {
			if r.Qualifier.Equals(o.Qualifier) {
				if r.Access == o.Access {
					return r.Target < o.Target
				}
				return r.Access < o.Access
			}
			return r.Qualifier.Less(o.Qualifier)
		}
		return r.Path < o.Path
	}
	return fileWeights[letterR] < fileWeights[letterO]
}

func (r *File) Equals(other any) bool {
	o, _ := other.(*File)
	return r.Path == o.Path && r.Access == o.Access &&
		r.Target == o.Target && r.Qualifier.Equals(o.Qualifier)
}
