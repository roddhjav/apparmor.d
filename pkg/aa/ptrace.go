// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

type Ptrace struct {
	Rule
	Qualifier
	Access string
	Peer   string
}

func newPtraceFromLog(log map[string]string) *Ptrace {
	return &Ptrace{
		Rule:      newRuleFromLog(log),
		Qualifier: newQualifierFromLog(log),
		Access:    toAccess(log["requested_mask"]),
		Peer:      log["peer"],
	}
}

func (r *Ptrace) Less(other any) bool {
	o, _ := other.(*Ptrace)
	if r.Access != o.Access {
		return r.Access < o.Access
	}
	if r.Peer != o.Peer {
		return r.Peer == o.Peer
	}
	return r.Qualifier.Less(o.Qualifier)
}

func (r *Ptrace) Equals(other any) bool {
	o, _ := other.(*Ptrace)
	return r.Access == o.Access && r.Peer == o.Peer &&
		r.Qualifier.Equals(o.Qualifier)
}
