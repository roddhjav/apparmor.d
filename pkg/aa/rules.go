// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

const (
	tokALLOW = "allow"
	tokAUDIT = "audit"
	tokDENY  = "deny"
)

type constraint uint

const (
	anyKind      constraint = iota // The rule can be found in either preamble or profile
	preambleKind                   // The rule can only be found in the preamble
	blockKind                      // The rule can only be found in a profile
)

// Rule generic interface for all AppArmor rules
type Rule interface {
	Less(other any) bool
	Equals(other any) bool
	String() string
	Constraint() constraint
	Kind() string
}

type Rules []Rule

func (r Rules) String() string {
	return renderTemplate("rules", r)
}
