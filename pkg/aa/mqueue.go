// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import "strings"

type Mqueue struct {
	Qualifier
	Access string
	Type   string
	Label  string
	Name   string
}

func MqueueFromLog(log map[string]string) ApparmorRule {
	mqueueType := "posix"
	if strings.Contains(log["class"], "posix") {
		mqueueType = "posix"
	} else if strings.Contains(log["class"], "sysv") {
		mqueueType = "sysv"
	}
	return &Mqueue{
		Qualifier: NewQualifierFromLog(log),
		Access:    toAccess(log["requested"]),
		Type:      mqueueType,
		Label:     log["label"],
		Name:      log["name"],
	}
}

func (r *Mqueue) Less(other any) bool {
	o, _ := other.(*Mqueue)
	if r.Qualifier.Equals(o.Qualifier) {
		if r.Access == o.Access {
			if r.Type == o.Type {
				return r.Label < o.Label
			}
			return r.Type < o.Type
		}
		return r.Access < o.Access
	}
	return r.Qualifier.Less(o.Qualifier)
}

func (r *Mqueue) Equals(other any) bool {
	o, _ := other.(*Mqueue)
	return r.Access == o.Access && r.Type == o.Type && r.Label == o.Label && r.Qualifier.Equals(o.Qualifier)
}
