// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package directive

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/arduino/go-paths-helper"
)

// Define the directive keyword globally
const Keyword = "#aa:"

// Build the profiles with the following directive applied
var Directives = map[string]Directive{}

var regDirective = regexp.MustCompile(`(?m).*` + Keyword + `([a-z]*) (.*)`)

// Main directive interface
type Directive interface {
	Usage() string
	Message() string
	Apply(opt *Option, profile string) string
}

type DirectiveBase struct {
	message string
	usage   string
}

func (d *DirectiveBase) Usage() string {
	return d.usage
}

func (d *DirectiveBase) Message() string {
	return d.message
}

// Directive options
type Option struct {
	Name string
	Args map[string]string
	File *paths.Path
	Raw  string
}

func NewOption(file *paths.Path, match []string) *Option {
	if len(match) != 3 {
		panic(fmt.Sprintf("Invalid directive: %v", match))
	}
	args := map[string]string{}
	for _, t := range strings.Fields(match[2]) {
		tmp := strings.Split(t, "=")
		if len(tmp) < 2 {
			args[tmp[0]] = ""
		} else {
			args[tmp[0]] = tmp[1]
		}
	}
	return &Option{
		Name: match[1],
		Args: args,
		File: file,
		Raw:  match[0],
	}
}

func Run(file *paths.Path, profile string) string {
	for _, match := range regDirective.FindAllStringSubmatch(profile, -1) {
		opt := NewOption(file, match)
		drtv, ok := Directives[opt.Name]
		if !ok {
			panic(fmt.Sprintf("Unknown directive: %s", opt.Name))
		}
		profile = drtv.Apply(opt, profile)
	}
	return profile
}
