// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

type IOUring struct {
	Qualifier
	Access string
	Label  string
}

func (r *IOUring) Less(other any) bool {
	o, _ := other.(*IOUring)
	if r.Qualifier.Equals(o.Qualifier) {
		if r.Access == o.Access {
			return r.Label < o.Label
		}
		return r.Access < o.Access
	}
	return r.Qualifier.Less(o.Qualifier)
}

func (r *IOUring) Equals(other any) bool {
	o, _ := other.(*IOUring)
	return r.Access == o.Access && r.Label == o.Label && r.Qualifier.Equals(o.Qualifier)
}
