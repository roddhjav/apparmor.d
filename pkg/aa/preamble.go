// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"slices"
)

type Abi struct {
	RuleBase
	Path    string
	IsMagic bool
}

func (r *Abi) Less(other any) bool {
	o, _ := other.(*Abi)
	if r.Path != o.Path {
		return r.Path < o.Path
	}
	return r.IsMagic == o.IsMagic
}

func (r *Abi) Equals(other any) bool {
	o, _ := other.(*Abi)
	return r.Path == o.Path && r.IsMagic == o.IsMagic
}

type Alias struct {
	RuleBase
	Path          string
	RewrittenPath string
}

func (r Alias) Less(other any) bool {
	o, _ := other.(*Alias)
	if r.Path != o.Path {
		return r.Path < o.Path
	}
	return r.RewrittenPath < o.RewrittenPath
}

func (r Alias) Equals(other any) bool {
	o, _ := other.(*Alias)
	return r.Path == o.Path && r.RewrittenPath == o.RewrittenPath
}

type Include struct {
	RuleBase
	IfExists bool
	Path     string
	IsMagic  bool
}

func (r *Include) Less(other any) bool {
	o, _ := other.(*Include)
	if r.Path == o.Path {
		return r.Path < o.Path
	}
	if r.IsMagic != o.IsMagic {
		return r.IsMagic
	}
	return r.IfExists
}

func (r *Include) Equals(other any) bool {
	o, _ := other.(*Include)
	return r.Path == o.Path && r.IsMagic == o.IsMagic && r.IfExists == o.IfExists
}

type Variable struct {
	RuleBase
	Name   string
	Values []string
}

func (r *Variable) Less(other any) bool {
	o, _ := other.(*Variable)
	if r.Name != o.Name {
		return r.Name < o.Name
	}
	return len(r.Values) < len(o.Values)
}

func (r *Variable) Equals(other any) bool {
	o, _ := other.(*Variable)
	return r.Name == o.Name && slices.Equal(r.Values, o.Values)
}
