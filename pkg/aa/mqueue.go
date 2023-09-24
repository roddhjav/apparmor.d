// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2023 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

type Mqueue struct {
	Qualifier
	Access string
	Type   string
	Label  string
}

func MqueueFromLog(log map[string]string, noNewPrivs, fileInherit bool) ApparmorRule {
	return &Mqueue{
		Qualifier: NewQualifier(false, noNewPrivs, fileInherit),
		Access:    maskToAccess[log["requested_mask"]],
		Type:      log["type"],
		Label:     log["label"],
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
