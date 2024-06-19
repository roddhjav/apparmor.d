// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import "fmt"

const (
	RLIMIT Kind = "rlimit"
)

func init() {
	requirements[RLIMIT] = requirement{
		"keys": {
			"cpu", "fsize", "data", "stack", "core", "rss", "nofile", "ofile",
			"as", "nproc", "memlock", "locks", "sigpending", "msgqueue", "nice",
			"rtprio", "rttime",
		},
	}
}

type Rlimit struct {
	RuleBase
	Key   string
	Op    string
	Value string
}

func newRlimitFromLog(log map[string]string) Rule {
	return &Rlimit{
		RuleBase: newRuleFromLog(log),
		Key:      log["key"],
		Op:       log["op"],
		Value:    log["value"],
	}
}

func (r *Rlimit) Validate() error {
	if err := validateValues(r.Kind(), "keys", []string{r.Key}); err != nil {
		return fmt.Errorf("%s: %w", r, err)
	}
	return nil
}

func (r *Rlimit) Compare(other Rule) int {
	o, _ := other.(*Rlimit)
	if res := compare(r.Key, o.Key); res != 0 {
		return res
	}
	if res := compare(r.Op, o.Op); res != 0 {
		return res
	}
	return compare(r.Value, o.Value)
}

func (r *Rlimit) String() string {
	return renderTemplate(r.Kind(), r)
}

func (r *Rlimit) Constraint() constraint {
	return blockKind
}

func (r *Rlimit) Kind() Kind {
	return RLIMIT
}
