// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package directive

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/arduino/go-paths-helper"
	"github.com/roddhjav/apparmor.d/pkg/prebuild/cfg"
)

var (
	// Define the directive keyword globally
	Keyword = "#aa:"

	// Build the profiles with the following directive applied
	Directives = map[string]Directive{}

	regDirective = regexp.MustCompile(`(?m).*` + Keyword + `([a-z]*) (.*)`)
)

// Main directive interface
type Directive interface {
	cfg.BaseInterface
	Apply(opt *Option, profile string) (string, error)
}

// Directive options
type Option struct {
	Name    string
	ArgMap  map[string]string
	ArgList []string
	File    *paths.Path
	Raw     string
}

func NewOption(file *paths.Path, match []string) *Option {
	if len(match) != 3 {
		panic(fmt.Sprintf("Invalid directive: %v", match))
	}
	argList := strings.Fields(match[2])
	argMap := map[string]string{}
	for _, t := range argList {
		tmp := strings.Split(t, "=")
		if len(tmp) < 2 {
			argMap[tmp[0]] = ""
		} else {
			argMap[tmp[0]] = tmp[1]
		}
	}
	return &Option{
		Name:    match[1],
		ArgMap:  argMap,
		ArgList: argList,
		File:    file,
		Raw:     match[0],
	}
}

func RegisterDirective(d Directive) {
	Directives[d.Name()] = d
}

func Run(file *paths.Path, profile string) (string, error) {
	var err error
	for _, match := range regDirective.FindAllStringSubmatch(profile, -1) {
		opt := NewOption(file, match)
		drtv, ok := Directives[opt.Name]
		if !ok {
			return "", fmt.Errorf("Unknown directive: %s", opt.Name)
		}
		profile, err = drtv.Apply(opt, profile)
		if err != nil {
			return "", err
		}
	}
	return profile, nil
}
