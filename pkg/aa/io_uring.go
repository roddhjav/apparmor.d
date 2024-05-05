// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import "slices"

const tokIOURING = "io_uring"


type IOUring struct {
	RuleBase
	Qualifier
	Access []string
	Label  string
}

func newIOUringFromLog(log map[string]string) Rule {
	return &IOUring{
		RuleBase:  newRuleFromLog(log),
		Qualifier: newQualifierFromLog(log),
		Access:    toAccess(tokIOURING, log["requested"]),
		Label:     log["label"],
	}
}

func (r *IOUring) Less(other any) bool {
	o, _ := other.(*IOUring)
	if len(r.Access) != len(o.Access) {
		return len(r.Access) < len(o.Access)
	}
	if r.Label != o.Label {
		return r.Label < o.Label
	}
	return r.Qualifier.Less(o.Qualifier)
}

func (r *IOUring) Equals(other any) bool {
	o, _ := other.(*IOUring)
	return slices.Equal(r.Access, o.Access) && r.Label == o.Label && r.Qualifier.Equals(o.Qualifier)
}

func (r *IOUring) String() string {
	return renderTemplate(r.Kind(), r)
}

func (r *IOUring) Constraint() constraint {
	return blockKind
}

func (r *IOUring) Kind() string {
	return tokIOURING
}
