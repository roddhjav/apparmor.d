// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

type Signal struct {
	RuleBase
	Qualifier
	Access []string
	Set    []string
	Peer   string
}

func newSignalFromLog(log map[string]string) Rule {
	return &Signal{
		RuleBase:  newRuleFromLog(log),
		Qualifier: newQualifierFromLog(log),
		Access:    toAccess(tokSIGNAL, log["requested_mask"]),
		Set:       toAccess(tokSIGNAL, log["signal"]),
		Peer:      log["peer"],
	}
}

func (r *Signal) Less(other any) bool {
	o, _ := other.(*Signal)
	if len(r.Access) != len(o.Access) {
		return len(r.Access) < len(o.Access)
	}
	if len(r.Set) != len(o.Set) {
		return len(r.Set) < len(o.Set)
	}
	if r.Peer != o.Peer {
		return r.Peer < o.Peer
	}
	return r.Qualifier.Less(o.Qualifier)
}

func (r *Signal) Equals(other any) bool {
	o, _ := other.(*Signal)
	return slices.Equal(r.Access, o.Access) && slices.Equal(r.Set, o.Set) &&
		r.Peer == o.Peer && r.Qualifier.Equals(o.Qualifier)
}
