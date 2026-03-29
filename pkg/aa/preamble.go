// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
)

const (
	ABI      Kind = "abi"
	ALIAS    Kind = "alias"
	INCLUDE  Kind = "include"
	VARIABLE Kind = "variable"
	BOOLEAN  Kind = "boolean"
	COMMENT  Kind = "comment"

	tokIFEXISTS = "if exists"
)

var (
	// reValidVarName matches valid variable names: must start with a letter,
	// followed by letters, digits, or underscores.
	reValidVarName = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*$`)
)

type Comment struct {
	Base
}

func newComment(rule rule) (Rule, error) {
	base := newBase(rule)
	base.IsLineRule = true
	return &Comment{Base: base}, nil
}

func (r *Comment) Kind() Kind {
	return COMMENT
}

func (r *Comment) Constraint() Constraint {
	return AnyRule
}

func (r *Comment) String() string {
	return renderTemplate(r.Kind(), r)
}

func (r *Comment) Validate() error {
	return nil
}

func (r *Comment) Compare(other Rule) int {
	return 0 // Comments are always equal to each other as they are not compared
}

func (r *Comment) Merge(other Rule) bool {
	return false // Never merge comments
}

func (r *Comment) Lengths() []int {
	return []int{} // No len for comments
}

func (r *Comment) setPaddings(max []int) {} // No paddings for comments

type Abi struct {
	Base
	Path    string
	IsMagic bool
}

func newAbi(q Qualifier, rule rule) (Rule, error) {
	var magic bool

	// Rejoin tokens that were split by spaces (e.g., "< includes/path >")
	path := strings.Join(rule.GetSlice(), " ")
	if path == "" {
		return nil, fmt.Errorf("invalid abi format: %s", rule)
	}
	switch path[0] {
	case '"':
		magic = false
		if !strings.HasSuffix(path, "\"") || len(path) < 3 {
			return nil, fmt.Errorf("invalid path %s in rule: %s", path, rule)
		}
	case '<':
		magic = true
		if !strings.HasSuffix(path, ">") || len(path) < 3 {
			return nil, fmt.Errorf("invalid path %s in rule: %s", path, rule)
		}
	default:
		return nil, fmt.Errorf("invalid path %s in rule: %s", path, rule)
	}
	path = strings.Trim(path, "\"<>")
	path = strings.Trim(path, " \t")
	path = strings.Trim(path, "\"")
	return &Abi{
		Base:    newBase(rule),
		Path:    path,
		IsMagic: magic,
	}, rule.ValidateMapKeys([]string{})
}

func (r *Abi) Kind() Kind {
	return ABI
}

func (r *Abi) Constraint() Constraint {
	return PreambleRule
}

func (r *Abi) String() string {
	return renderTemplate(r.Kind(), r)
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

func (r *Abi) Merge(other Rule) bool {
	return false // Never merge abi
}

func (r *Abi) Lengths() []int {
	return []int{} // No len for abi
}

func (r *Abi) setPaddings(max []int) {} // No paddings for abi

type Alias struct {
	Base
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
		Base:          newBase(rule),
		Path:          rule.Get(0),
		RewrittenPath: rule.Get(2),
	}, rule.ValidateMapKeys([]string{})
}

func (r *Alias) Kind() Kind {
	return ALIAS
}

func (r *Alias) Constraint() Constraint {
	return PreambleRule
}

func (r *Alias) String() string {
	return renderTemplate(r.Kind(), r)
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

func (r *Alias) Merge(other Rule) bool {
	return false // Never merge alias
}

func (r *Alias) Lengths() []int {
	return []int{} // No len for alias
}

func (r *Alias) setPaddings(max []int) {} // No paddings for alias

type Include struct {
	Base
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

	path := strings.Join(r, " ") // Rejoin in case of spaces in path
	switch path[0] {
	case '"':
		magic = false
		if !strings.HasSuffix(path, "\"") || len(path) < 3 {
			return nil, fmt.Errorf("invalid path %s in rule: %s", path, rule)
		}
		path = strings.Trim(path, "\"")
	case '<':
		magic = true
		if !strings.HasSuffix(path, ">") || len(path) < 3 {
			return nil, fmt.Errorf("invalid path %s in rule: %s", path, rule)
		}
		path = strings.Trim(path, "<>")
		path = strings.Trim(path, " \t")
	default:
		// Allow bare/relative paths (e.g., simple_tests/includes/...)
		magic = false
	}
	return &Include{
		Base:     newBase(rule),
		IfExists: ifexists,
		Path:     path,
		IsMagic:  magic,
	}, rule.ValidateMapKeys([]string{})
}

func (r *Include) Kind() Kind {
	return INCLUDE
}

func (r *Include) Constraint() Constraint {
	return AnyRule
}

func (r *Include) String() string {
	return renderTemplate(r.Kind(), r)
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

func (r *Include) Merge(other Rule) bool {
	return false // Never merge include
}

func (r *Include) Lengths() []int {
	return []int{} // No len for include
}

func (r *Include) setPaddings(max []int) {} // No paddings for include

type Variable struct {
	Base
	Name   string
	Values []string
	Define bool
}

func newVariable(rule rule) (Rule, error) {
	var define bool
	var values []string
	if len(rule) < 3 {
		return nil, fmt.Errorf("invalid variable format: %v", rule)
	}

	r := rule.GetSlice()
	name := strings.Trim(rule.Get(0), VARIABLE.Tok()+"}")
	if !reValidVarName.MatchString(name) {
		return nil, fmt.Errorf("invalid variable name '%s': must start with a letter", name)
	}

	switch rule.Get(1) {
	case tokEQUAL, "?=", ":=":
		define = true
		values = r[2:]
	case tokPLUS + tokEQUAL:
		define = false
		values = r[2:]
	default:
		return nil, fmt.Errorf("invalid operator in variable: %v", rule)
	}

	// Validate variable values
	for _, v := range values {
		if strings.HasSuffix(v, ",") {
			return nil, fmt.Errorf("trailing comma in variable value: %s", v)
		}
		quoteCount := strings.Count(v, "\"")
		if quoteCount%2 != 0 {
			return nil, fmt.Errorf("unbalanced quotes in variable value: %s", v)
		}
		if strings.Contains(v, "!") {
			return nil, fmt.Errorf("invalid character '!' in variable value: %s", v)
		}
	}

	return &Variable{
		Base:   newBase(rule),
		Name:   name,
		Values: values,
		Define: define,
	}, nil
}

func (r *Variable) Kind() Kind {
	return VARIABLE
}

func (r *Variable) Constraint() Constraint {
	return PreambleRule
}

func (r *Variable) String() string {
	return renderTemplate(r.Kind(), r)
}

func (r *Variable) Validate() error {
	return nil
}

func (r *Variable) Compare(other Rule) int {
	o, _ := other.(*Variable)
	if res := compare(r.Name, o.Name); res != 0 {
		return res
	}
	if res := compare(r.Define, o.Define); res != 0 {
		return res
	}
	return compare(r.Values, o.Values)
}

func (r *Variable) Merge(other Rule) bool {
	o, _ := other.(*Variable)

	if r.Name == o.Name && r.Define == o.Define {
		r.Values = merge(r.Kind(), "access", r.Values, o.Values)
		b := &r.Base
		return b.merge(o.Base)
	}
	return false
}

func (r *Variable) Lengths() []int {
	return []int{} // No len for variable
}

func (r *Variable) setPaddings(max []int) {} // No paddings for variable

type Boolean struct {
	Base
	Name  string
	Value bool
}

func newBoolean(rule rule) (Rule, error) {
	name, value := "", ""

	switch len(rule) {
	case 1:
		name = strings.Trim(rule.Get(0), BOOLEAN.Tok()+"{}")
		value = rule.GetValuesAsString(rule.Get(0))

	case 3:
		name = strings.Trim(rule.Get(0), BOOLEAN.Tok()+"{}")
		op := rule.Get(1)
		if op != tokEQUAL && op != "?=" && op != ":=" {
			return nil, fmt.Errorf("invalid boolean format, missing %s in: %s", tokEQUAL, rule)
		}
		value = rule.Get(2)

	case 4:
		// Handle ?= and := operators: "$VAR ? = value" or "$VAR : = value"
		name = strings.Trim(rule.Get(0), BOOLEAN.Tok()+"{}")
		op := rule.Get(1)
		if (op == "?" || op == ":") && rule.Get(2) == tokEQUAL {
			value = rule.Get(3)
		} else {
			return nil, fmt.Errorf("invalid boolean format: %v", rule)
		}

	default:
		return nil, fmt.Errorf("invalid boolean format: %v", rule)
	}

	valueLower := strings.ToLower(value)
	if !slices.Contains([]string{"true", "false"}, valueLower) {
		return nil, fmt.Errorf("invalid boolean value %s in rule: %s", value, rule)
	}
	return &Boolean{
		Base:  newBase(rule),
		Name:  name,
		Value: valueLower == "true",
	}, nil
}

func (r *Boolean) Kind() Kind {
	return BOOLEAN
}

func (r *Boolean) Constraint() Constraint {
	return PreambleRule
}

func (r *Boolean) String() string {
	return renderTemplate(r.Kind(), r)
}

func (r *Boolean) Validate() error {
	return nil
}

func (r *Boolean) Compare(other Rule) int {
	o, _ := other.(*Boolean)
	if res := compare(r.Name, o.Name); res != 0 {
		return res
	}
	return compare(r.Value, o.Value)
}

func (r *Boolean) Merge(other Rule) bool {
	o, _ := other.(*Boolean)

	if r.Name == o.Name && r.Value == o.Value {
		b := &r.Base
		return b.merge(o.Base)
	}
	return false
}

func (r *Boolean) Lengths() []int {
	return []int{} // No len for boolean
}

func (r *Boolean) setPaddings(max []int) {} // No paddings for boolean
