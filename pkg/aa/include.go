// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

type Include struct {
	IfExists bool
	Path     string
	IsMagic  bool
}

func (r *Include) Less(other any) bool {
	o, _ := other.(*Include)
	if r.Path == o.Path {
		if r.IsMagic == o.IsMagic {
			return r.IfExists
		}
		return r.IsMagic
	}
	return r.Path < o.Path
}

func (r *Include) Equals(other any) bool {
	o, _ := other.(*Include)
	return r.Path == o.Path && r.IsMagic == o.IsMagic &&
		r.IfExists == o.IfExists
}
