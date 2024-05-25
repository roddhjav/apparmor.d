// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

const (
	tokALL = "all"
)

type All struct {
	RuleBase
}

func (r *All) Validate() error {
	return nil
}

func (r *All) Less(other any) bool {
	return false
}

func (r *All) Equals(other any) bool {
	return false
}

func (r *All) String() string {
	return renderTemplate(r.Kind(), r)
}

func (r *All) Constraint() constraint {
	return blockKind
}

func (r *All) Kind() string {
	return tokALL
}
