// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"reflect"
	"sort"
)

const (
	tokALL   = "all"
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

// Sort the rules in a profile.
// Follow: https://apparmor.pujol.io/development/guidelines/#guidelines
func (r Rules) Sort() {
	sort.Slice(r, func(i, j int) bool {
		typeOfI := reflect.TypeOf(r[i])
		typeOfJ := reflect.TypeOf(r[j])
		if typeOfI != typeOfJ {
			valueOfI := typeToValue(typeOfI)
			valueOfJ := typeToValue(typeOfJ)
			if typeOfI == reflect.TypeOf((*Include)(nil)) && r[i].(*Include).IfExists {
				valueOfI = "include_if_exists"
			}
			if typeOfJ == reflect.TypeOf((*Include)(nil)) && r[j].(*Include).IfExists {
				valueOfJ = "include_if_exists"
			}
			return ruleWeights[valueOfI] < ruleWeights[valueOfJ]
		}
		return r[i].Less(r[j])
	})
}
