// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2023 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa
// AppArmorProfiles represents a full set of apparmor profiles
type AppArmorProfiles map[string]*AppArmorProfile

// ApparmorProfile represents a full apparmor profile.
// Warning: close to the BNF grammar of apparmor profile but not exactly the same (yet):
//   - Some rules are not supported yet (subprofile, hat...)
//   - The structure is simplified as it only aims at writting profile, not parsing it.
type AppArmorProfile struct {
	Preamble
	Profile
}

func NewAppArmorProfile() *AppArmorProfile {
	return &AppArmorProfile{}
}
