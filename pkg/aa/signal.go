// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2023 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

type Signal struct {
	Qualifier
	Access string
	Set    string
	Peer   string
}

func SignalFromLog(log map[string]string) ApparmorRule {
	return &Signal{
		Qualifier: NewQualifierFromLog(log),
		Access:    maskToAccess[log["requested_mask"]],
		Set:       log["signal"],
		Peer:      log["peer"],
	}
}

func (r *Signal) Less(other any) bool {
	o, _ := other.(*Signal)
	if r.Qualifier.Equals(o.Qualifier) {
		if r.Access == o.Access {
			if r.Set == o.Set {
				return r.Peer < o.Peer
			}
			return r.Set < o.Set
		}
		return r.Access < o.Access
	}
	return r.Qualifier.Less(o.Qualifier)
}

func (r *Signal) Equals(other any) bool {
	o, _ := other.(*Signal)
	return r.Access == o.Access && r.Set == o.Set &&
		r.Peer == o.Peer && r.Qualifier.Equals(o.Qualifier)
}
