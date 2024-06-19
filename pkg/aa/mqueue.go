// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"fmt"
	"slices"
	"strings"
)

const MQUEUE Kind = "mqueue"

func init() {
	requirements[MQUEUE] = requirement{
		"access": []string{
			"r", "w", "rw", "read", "write", "create", "open",
			"delete", "getattr", "setattr",
		},
		"type": []string{"posix", "sysv"},
	}
}

type Mqueue struct {
	RuleBase
	Qualifier
	Access []string
	Type   string
	Label  string
	Name   string
}

func newMqueueFromLog(log map[string]string) Rule {
	mqueueType := "posix"
	if strings.Contains(log["class"], "posix") {
		mqueueType = "posix"
	} else if strings.Contains(log["class"], "sysv") {
		mqueueType = "sysv"
	}
	return &Mqueue{
		RuleBase:  newRuleFromLog(log),
		Qualifier: newQualifierFromLog(log),
		Access:    Must(toAccess(MQUEUE, log["requested"])),
		Type:      mqueueType,
		Label:     log["label"],
		Name:      log["name"],
	}
}

func (r *Mqueue) Validate() error {
	if err := validateValues(r.Kind(), "access", r.Access); err != nil {
		return fmt.Errorf("%s: %w", r, err)
	}
	if err := validateValues(r.Kind(), "type", []string{r.Type}); err != nil {
		return fmt.Errorf("%s: %w", r, err)
	}
	return nil
}

func (r *Mqueue) Compare(other Rule) int {
	o, _ := other.(*Mqueue)
	if res := compare(r.Access, o.Access); res != 0 {
		return res
	}
	if res := compare(r.Type, o.Type); res != 0 {
		return res
	}
	if res := compare(r.Label, o.Label); res != 0 {
		return res
	}
	return r.Qualifier.Compare(o.Qualifier)
}

func (r *Mqueue) String() string {
	return renderTemplate(r.Kind(), r)
}

func (r *Mqueue) Constraint() constraint {
	return blockKind
}

func (r *Mqueue) Kind() Kind {
	return MQUEUE
}
