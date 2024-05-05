// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"slices"
)

const tokCAPABILITY = "capability"

type Capability struct {
	RuleBase
	Qualifier
	Names []string
}

}

func newCapabilityFromLog(log map[string]string) Rule {
	return &Capability{
		RuleBase:  newRuleFromLog(log),
		Qualifier: newQualifierFromLog(log),
		Names:     []string{log["capname"]},
	}
}

func (r *Capability) Less(other any) bool {
	o, _ := other.(*Capability)
	for i := 0; i < len(r.Names) && i < len(o.Names); i++ {
		if r.Names[i] != o.Names[i] {
			return r.Names[i] < o.Names[i]
		}
	}
	return r.Qualifier.Less(o.Qualifier)
}

func (r *Capability) Equals(other any) bool {
	o, _ := other.(*Capability)
	return slices.Equal(r.Names, o.Names) && r.Qualifier.Equals(o.Qualifier)
}

func (r *Capability) String() string {
	return renderTemplate(r.Kind(), r)
}

func (r *Capability) Constraint() constraint {
	return blockKind
}

func (r *Capability) Kind() string {
	return tokCAPABILITY
}
