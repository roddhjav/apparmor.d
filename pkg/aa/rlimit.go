// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

type Rlimit struct {
	Rule
	Key   string
	Op    string
	Value string
}

func newRlimitFromLog(log map[string]string) *Rlimit {
	return &Rlimit{
		Rule:  newRuleFromLog(log),
		Key:   log["key"],
		Op:    log["op"],
		Value: log["value"],
	}
}

func (r *Rlimit) Less(other any) bool {
	o, _ := other.(*Rlimit)
	if r.Key != o.Key {
		return r.Key < o.Key
	}
	if r.Op != o.Op {
		return r.Op < o.Op
	}
	return r.Value < o.Value
}

func (r *Rlimit) Equals(other any) bool {
	o, _ := other.(*Rlimit)
	return r.Key == o.Key && r.Op == o.Op && r.Value == o.Value
}
