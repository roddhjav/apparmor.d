// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

// Package state detects system state values and edits the state tunable files
package state

import (
	"strings"

	"github.com/roddhjav/apparmor.d/pkg/aa"
	"github.com/roddhjav/apparmor.d/pkg/paths"
)

// Key normalizes a variable name to its `@{NAME}` form.
func Key(name string) string {
	if strings.HasPrefix(name, "@{") {
		return name
	}
	return "@{" + name + "}"
}

// varName returns the aa.Variable name form of a key, without `@{}`.
func varName(key string) string {
	return strings.TrimSuffix(strings.TrimPrefix(key, "@{"), "}")
}

// File is a state tunable file
type File struct {
	path *paths.Path
	file *aa.AppArmorProfileFile
}

// Load reads and parses a tunable file.
func Load(path *paths.Path) (*File, error) {
	content, err := path.ReadFileAsString()
	if err != nil {
		return nil, err
	}
	file := &aa.AppArmorProfileFile{}
	if _, err := file.Parse(content); err != nil {
		return nil, err
	}
	return &File{path: path, file: file}, nil
}

// NewFile returns an empty in-memory tunable file, created on Save.
func NewFile(path *paths.Path) *File {
	return &File{path: path, file: &aa.AppArmorProfileFile{}}
}

// variable returns the definition of a variable, nil when not defined.
func (f *File) variable(key string) *aa.Variable {
	name := varName(key)
	for _, rule := range f.file.Preamble {
		if v, ok := rule.(*aa.Variable); ok && v.Name == name {
			return v
		}
	}
	return nil
}

// Get returns the value of a variable definition.
func (f *File) Get(name string) (string, bool) {
	if v := f.variable(name); v != nil {
		return strings.Join(v.Values, " "), true
	}
	return "", false
}

// Names returns the defined variable names, in file order.
func (f *File) Names() []string {
	var names []string
	for _, rule := range f.file.Preamble {
		if v, ok := rule.(*aa.Variable); ok {
			names = append(names, Key(v.Name))
		}
	}
	return names
}

// Set replaces the value of an existing definition in place. It reports
// whether the variable was found.
func (f *File) Set(name, value string) bool {
	if v := f.variable(name); v != nil {
		v.Values = strings.Fields(value)
		return true
	}
	return false
}

// Add appends a definition at the end of the file.
func (f *File) Add(name, value string) {
	f.file.Preamble = append(f.file.Preamble, &aa.Variable{
		Name:   varName(name),
		Values: strings.Fields(value),
		Define: true,
	})
}

// Save writes the file back, creating parent directories as needed.
func (f *File) Save() error {
	if err := f.path.Parent().MkdirAll(); err != nil {
		return err
	}
	return f.path.WriteFile([]byte(f.file.String()))
}
