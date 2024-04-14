// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

type Signal struct {
	Rule
	Qualifier
	Access string
	Set    string
	Peer   string
}

func newSignalFromLog(log map[string]string) *Signal {
	return &Signal{
		Rule:      newRuleFromLog(log),
		Qualifier: newQualifierFromLog(log),
		Access:    toAccess(log["requested_mask"]),
		Set:       log["signal"],
		Peer:      log["peer"],
	}
}

func (r *Signal) Less(other any) bool {
	o, _ := other.(*Signal)
	if r.Access != o.Access {
		return r.Access < o.Access
	}
	if r.Set != o.Set {
		return r.Set < o.Set
	}
	if r.Peer != o.Peer {
		return r.Peer < o.Peer
	}
	return r.Qualifier.Less(o.Qualifier)
}

func (r *Signal) Equals(other any) bool {
	o, _ := other.(*Signal)
	return r.Access == o.Access && r.Set == o.Set &&
		r.Peer == o.Peer && r.Qualifier.Equals(o.Qualifier)
}
