// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"fmt"
)

const IOURING Kind = "io_uring"

func init() {
	requirements[IOURING] = requirement{
		"access": []string{"sqpoll", "override_creds"},
	}
}

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
		Access:    Must(toAccess(IOURING, log["requested"])),
		Label:     log["label"],
	}
}

func (r *IOUring) Validate() error {
	if err := validateValues(r.Kind(), "access", r.Access); err != nil {
		return fmt.Errorf("%s: %w", r, err)
	}
	return nil
}

func (r *IOUring) Compare(other Rule) int {
	o, _ := other.(*IOUring)
	if res := compare(r.Access, o.Access); res != 0 {
		return res
	}
	if res := compare(r.Label, o.Label); res != 0 {
		return res
	}
	return r.Qualifier.Compare(o.Qualifier)
}

func (r *IOUring) String() string {
	return renderTemplate(r.Kind(), r)
}

func (r *IOUring) Constraint() constraint {
	return blockKind
}

func (r *IOUring) Kind() Kind {
	return IOURING
}
