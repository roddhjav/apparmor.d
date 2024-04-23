// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"slices"

const (
	tokABI      = "abi"
	tokALIAS    = "alias"
	tokINCLUDE  = "include"
	tokIFEXISTS = "if exists"
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

func (r *Abi) String() string {
	return renderTemplate(tokABI, r)
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

func (r *Alias) String() string {
	return renderTemplate(tokALIAS, r)
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

func (r *Include) String() string {
	return renderTemplate(tokINCLUDE, r)
}

type Variable struct {
	RuleBase
	Name   string
	Values []string
	Define bool
}

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

func (r *Variable) String() string {
	return renderTemplate("variable", r)
}
