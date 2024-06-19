// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"fmt"
	"strings"
)

const (
	ABI      Kind = "abi"
	ALIAS    Kind = "alias"
	INCLUDE  Kind = "include"
	VARIABLE Kind = "variable"
	COMMENT  Kind = "comment"

	tokIFEXISTS = "if exists"
)

type Comment struct {
	RuleBase
}

func newComment(rule rule) (Rule, error) {
	base := newBase(rule)
	base.IsLineRule = true
	return &Comment{RuleBase: base}, nil
}

func (r *Comment) Validate() error {
	return nil
}

func (r *Comment) Compare(other Rule) int {
	return 0
}

func (r *Comment) String() string {
	return renderTemplate(r.Kind(), r)
}

func (r *Comment) Constraint() constraint {
	return anyKind
}

func (r *Comment) Kind() Kind {
	return COMMENT
}

type Abi struct {
	RuleBase
	Path    string
	IsMagic bool
}

func newAbi(q Qualifier, rule rule) (Rule, error) {
	var magic bool
	if len(rule) != 1 {
		return nil, fmt.Errorf("invalid abi format: %s", rule)
	}

	path := rule.Get(0)
	switch {
	case path[0] == '"':
		magic = false
	case path[0] == '<':
		magic = true
	default:
		return nil, fmt.Errorf("invalid path %s in rule: %s", path, rule)
	}
	return &Abi{
		RuleBase: newBase(rule),
		Path:     strings.Trim(path, "\"<>"),
		IsMagic:  magic,
	}, nil
}

func (r *Abi) Validate() error {
	return nil
}

func (r *Abi) Compare(other Rule) int {
	o, _ := other.(*Abi)
	if res := compare(r.Path, o.Path); res != 0 {
		return res
	}
	return compare(r.IsMagic, o.IsMagic)
}

func (r *Abi) String() string {
	return renderTemplate(r.Kind(), r)
}

func (r *Abi) Constraint() constraint {
	return preambleKind
}

func (r *Abi) Kind() Kind {
	return ABI
}

type Alias struct {
	RuleBase
	Path          string
	RewrittenPath string
}

func newAlias(q Qualifier, rule rule) (Rule, error) {
	if len(rule) != 3 {
		return nil, fmt.Errorf("invalid alias format: %s", rule)
	}
	if rule.Get(1) != tokARROW {
		return nil, fmt.Errorf("invalid alias format, missing %s in: %s", tokARROW, rule)
	}
	return &Alias{
		RuleBase:      newBase(rule),
		Path:          rule.Get(0),
		RewrittenPath: rule.Get(2),
	}, nil
}

func (r *Alias) Validate() error {
	return nil
}

func (r *Alias) Compare(other Rule) int {
	o, _ := other.(*Alias)
	if res := compare(r.Path, o.Path); res != 0 {
		return res
	}
	return compare(r.RewrittenPath, o.RewrittenPath)
}

func (r *Alias) String() string {
	return renderTemplate(r.Kind(), r)
}

func (r *Alias) Constraint() constraint {
	return preambleKind
}

func (r *Alias) Kind() Kind {
	return ALIAS
}

type Include struct {
	RuleBase
	IfExists bool
	Path     string
	IsMagic  bool
}

func newInclude(rule rule) (Rule, error) {
	var magic bool
	var ifexists bool

	size := len(rule)
	if size == 0 {
		return nil, fmt.Errorf("invalid include format: %v", rule)
	}

	r := rule.GetSlice()
	if size >= 3 && strings.Join(r[:2], " ") == tokIFEXISTS {
		ifexists = true
		r = r[2:]
	}

	path := r[0]
	switch {
	case path[0] == '"':
		magic = false
	case path[0] == '<':
		magic = true
	default:
		return nil, fmt.Errorf("invalid path format: %v", path)
	}
	return &Include{
		RuleBase: newBase(rule),
		IfExists: ifexists,
		Path:     strings.Trim(path, "\"<>"),
		IsMagic:  magic,
	}, nil
}

func (r *Include) Validate() error {
	return nil
}

func (r *Include) Compare(other Rule) int {
	const base = "abstractions/base"
	o, _ := other.(*Include)
	if res := compare(r.Path, o.Path); res != 0 {
		if r.Path == base {
			return -1
		}
		if o.Path == base {
			return 1
		}
		return res
	}
	if res := compare(r.IsMagic, o.IsMagic); res != 0 {
		return res
	}
	return compare(r.IfExists, o.IfExists)
}

func (r *Include) String() string {
	return renderTemplate(r.Kind(), r)
}

func (r *Include) Constraint() constraint {
	return anyKind
}

func (r *Include) Kind() Kind {
	return INCLUDE
}

type Variable struct {
	RuleBase
	Name   string
	Values []string
	Define bool
}

func newVariableFromRule(rule rule) (Rule, error) {
	var define bool
	var values []string
	if len(rule) < 3 {
		return nil, fmt.Errorf("invalid variable format: %v", rule)
	}

	r := rule.GetSlice()
	name := strings.Trim(rule.Get(0), VARIABLE.Tok()+"}")
	switch rule.Get(1) {
	case tokEQUAL:
		define = true
		values = r[2:]
	case tokPLUS + tokEQUAL:
		define = false
		values = r[2:]
	default:
		return nil, fmt.Errorf("invalid operator in variable: %v", rule)
	}
	return &Variable{
		RuleBase: newBase(rule),
		Name:     name,
		Values:   values,
		Define:   define,
	}, nil
}

func (r *Variable) Validate() error {
	return nil
}

func (r *Variable) Compare(other Rule) int {
	return 0
}

func (r *Variable) String() string {
	return renderTemplate(r.Kind(), r)
}

func (r *Variable) Constraint() constraint {
	return preambleKind
}

func (r *Variable) Kind() Kind {
	return VARIABLE
}
