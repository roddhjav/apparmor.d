// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"maps"
	"slices"
	"strings"
)

// Profile represents a single AppArmor profile.
type Profile struct {
	Rule
	Header
	Rules Rules
}

// Header represents the header of a profile.
type Header struct {
	Name        string
	Attachments []string
	Attributes  map[string]string
	Flags       []string
}

func (r *Profile) Less(other any) bool {
	o, _ := other.(*Profile)
	if r.Name != o.Name {
		return r.Name < o.Name
	}
	return len(r.Attachments) < len(o.Attachments)
}

func (r *Profile) Equals(other any) bool {
	o, _ := other.(*Profile)
	return r.Name == o.Name && slices.Equal(r.Attachments, o.Attachments) &&
		maps.Equal(r.Attributes, o.Attributes) &&
		slices.Equal(r.Flags, o.Flags)
}
