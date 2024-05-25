// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"slices"
	"strings"
)

const tokMQUEUE = "mqueue"

func init() {
	requirements[tokMQUEUE] = requirement{
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
		Access:    toAccess(tokMQUEUE, log["requested"]),
		Type:      mqueueType,
		Label:     log["label"],
		Name:      log["name"],
	}
}

func (r *Mqueue) Less(other any) bool {
	o, _ := other.(*Mqueue)
	if len(r.Access) != len(o.Access) {
		return len(r.Access) < len(o.Access)
	}
	if r.Type != o.Type {
		return r.Type < o.Type
	}
	if r.Label != o.Label {
		return r.Label < o.Label
	}
	return r.Qualifier.Less(o.Qualifier)
}

func (r *Mqueue) Equals(other any) bool {
	o, _ := other.(*Mqueue)
	return slices.Equal(r.Access, o.Access) && r.Type == o.Type && r.Label == o.Label &&
		r.Name == o.Name && r.Qualifier.Equals(o.Qualifier)
}

func (r *Mqueue) String() string {
	return renderTemplate(r.Kind(), r)
}

func (r *Mqueue) Constraint() constraint {
	return blockKind
}

func (r *Mqueue) Kind() string {
	return tokMQUEUE
}
