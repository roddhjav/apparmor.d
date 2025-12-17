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
	Base
	Key   string
	Op    string
	Value string
}

func newRlimit(q Qualifier, rule rule) (Rule, error) {
	if len(rule) != 4 {
		return nil, fmt.Errorf("invalid set format: %s", rule)
	}
	if rule.Get(0) != RLIMIT.Tok() {
		return nil, fmt.Errorf("invalid rlimit format: %s", rule)
	}
	return &Rlimit{
		Base:  newBase(rule),
		Key:   rule.Get(1),
		Op:    rule.Get(2),
		Value: rule.Get(3),
	}, rule.ValidateMapKeys([]string{})
}

func newRlimitFromLog(log map[string]string) Rule {
	return &Rlimit{
		Base:  newBaseFromLog(log),
		Key:   log["rlimit"],
		Op:    "<=",
		Value: log["value"],
	}
}

func (r *Rlimit) Kind() Kind {
	return RLIMIT
}

func (r *Rlimit) Constraint() Constraint {
	return BlockRule
}

func (r *Rlimit) String() string {
	return renderTemplate(r.Kind(), r)
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

func (r *Rlimit) Merge(other Rule) bool {
	return false // Never merge rlimit
}

func (r *Rlimit) Lengths() []int {
	return []int{
		length("", r.Key),
		length("", r.Op),
		length("", r.Value),
	}
}

func (r *Rlimit) setPaddings(max []int) {
	r.Paddings = setPaddings(
		max, []string{"", "", ""},
		[]any{r.Key, r.Op, r.Value},
	)
}
