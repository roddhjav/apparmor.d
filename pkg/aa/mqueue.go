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

func newMqueue(q Qualifier, rule rule) (Rule, error) {
	access, name := "", ""
	r := rule.GetSlice()
	size := len(r)
	if size > 0 {
		access = strings.Join(r[:size-1], " ")
		name = r[size-1]
		if slices.Contains(requirements[MQUEUE]["access"], name) {
			access += " " + name
		}
	}
	accesses, err := toAccess(MQUEUE, access)
	if err != nil {
		return nil, err
	}
	return &Mqueue{
		RuleBase:  newBase(rule),
		Qualifier: q,
		Access:    accesses,
		Type:      rule.GetValuesAsString("type"),
		Label:     rule.GetValuesAsString("label"),
		Name:      name,
	}, nil
}

func newMqueueFromLog(log map[string]string) Rule {
	mqueueType := "posix"
	if strings.Contains(log["class"], "posix") {
		mqueueType = "posix"
	} else if strings.Contains(log["class"], "sysv") {
		mqueueType = "sysv"
	}
	return &Mqueue{
		RuleBase:  newBaseFromLog(log),
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

func (r *Mqueue) Merge(other Rule) bool {
	o, _ := other.(*Mqueue)

	if !r.Qualifier.Equal(o.Qualifier) {
		return false
	}
	if r.Type == o.Type && r.Label == o.Label && r.Name == o.Name {
		r.Access = merge(r.Kind(), "access", r.Access, o.Access)
		b := &r.RuleBase
		return b.merge(o.RuleBase)
	}
	return false
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
