// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"fmt"
	"slices"
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

func newComment(rule []string) (Rule, error) {
	base := newRule(rule)
	base.IsLineRule = true
	return &Comment{RuleBase: base}, nil
}

func (r *Comment) Validate() error {
	return nil
}

func (r *Comment) Less(other any) bool {
	return false
}

func (r *Comment) Equals(other any) bool {
	return false
}

func (r *Comment) String() string {
	return renderTemplate(r.Kind(), r)
}

func (r *Comment) IsPreamble() bool {
	return true
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

func newAbi(rule []string) (Rule, error) {
	var magic bool
	if len(rule) > 0 && rule[0] == ABI.Tok() {
		rule = rule[1:]
	}
	if len(rule) != 1 {
		return nil, fmt.Errorf("invalid abi format: %s", rule)
	}

	path := rule[0]
	switch {
	case path[0] == '"':
		magic = false
	case path[0] == '<':
		magic = true
	default:
		return nil, fmt.Errorf("invalid path %s in rule: %s", path, rule)
	}
	return &Abi{
		RuleBase: newRule(rule),
		Path:     strings.Trim(path, "\"<>"),
		IsMagic:  magic,
	}, nil
}

func (r *Abi) Validate() error {
	return nil
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

func newAlias(rule []string) (Rule, error) {
	if len(rule) > 0 && rule[0] == ALIAS.Tok() {
		rule = rule[1:]
	}
	if len(rule) != 3 {
		return nil, fmt.Errorf("invalid alias format: %s", rule)
	}
	if rule[1] != tokARROW {
		return nil, fmt.Errorf("invalid alias format, missing %s in: %s", tokARROW, rule)
	}
	return &Alias{
		RuleBase:      newRule(rule),
		Path:          rule[0],
		RewrittenPath: rule[2],
	}, nil
}

func (r *Alias) Validate() error {
	return nil
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

func newInclude(rule []string) (Rule, error) {
	var magic bool
	var ifexists bool

	if len(rule) > 0 && rule[0] == INCLUDE.Tok() {
		rule = rule[1:]
	}

	size := len(rule)
	if size == 0 {
		return nil, fmt.Errorf("invalid include format: %v", rule)
	}

	if size >= 3 && strings.Join(rule[:2], " ") == tokIFEXISTS {
		ifexists = true
		rule = rule[2:]
	}

	path := rule[0]
	switch {
	case path[0] == '"':
		magic = false
	case path[0] == '<':
		magic = true
	default:
		return nil, fmt.Errorf("invalid path format: %v", path)
	}
	return &Include{
		RuleBase: newRule(rule),
		IfExists: ifexists,
		Path:     strings.Trim(path, "\"<>"),
		IsMagic:  magic,
	}, nil
}

func (r *Include) Validate() error {
	return nil
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

func newVariable(rule []string) (Rule, error) {
	var define bool
	var values []string
	if len(rule) < 3 {
		return nil, fmt.Errorf("invalid variable format: %v", rule)
	}

	name := strings.Trim(rule[0], VARIABLE.Tok()+"}")
	switch rule[1] {
	case tokEQUAL:
		define = true
		values = tokensStripComment(rule[2:])
	case tokPLUS:
		if rule[2] != tokEQUAL {
			return nil, fmt.Errorf("invalid operator in variable: %v", rule)
		}
		define = false
		values = tokensStripComment(rule[3:])
	default:
		return nil, fmt.Errorf("invalid operator in variable: %v", rule)
	}
	return &Variable{
		RuleBase: newRule(rule),
		Name:     name,
		Values:   values,
		Define:   define,
	}, nil
}

func (r *Variable) Validate() error {
	return nil
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
	return renderTemplate(r.Kind(), r)
}

func (r *Variable) Constraint() constraint {
	return preambleKind
}

func (r *Variable) Kind() Kind {
	return VARIABLE
}
