// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package directive

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/roddhjav/apparmor.d/pkg/paths"
	"github.com/roddhjav/apparmor.d/pkg/prebuild"
)

var (
	// Define the directive keyword globally
	Keyword = "#aa:"

	// Build the profiles with the following directive applied
	Directives = map[string]Directive{}

	regDirective = regexp.MustCompile(`(?m).*` + Keyword + `([a-z]*)( .*)?`)
)

// Directive main interface
type Directive interface {
	prebuild.BaseInterface
	Apply(opt *Option, profile string) (string, error)
}

func Usage() string {
	res := "Directive:\n"
	for _, d := range Directives {
		for _, h := range d.Usage() {
			res += fmt.Sprintf("    %s%s %s\n", Keyword, d.Name(), h)
		}
	}
	return res
}

// Option for the directive
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

// Clean removes selected directive line from input string.
// Useful to remove directive text applied on some condition only
func (o *Option) Clean(input string) string {
	return strings.Replace(input, o.Raw, o.cleanKeyword(o.Raw), 1)
}

// cleanKeyword removes the dirextive keywork (#aa:...) from the input string
func (o *Option) cleanKeyword(input string) string {
	reg := regexp.MustCompile(`\s*` + Keyword + o.Name + `( .*)?$`)
	return reg.ReplaceAllString(input, "")
}

// IsInline checks if either the directive is in one line or if it is a paragraph
func (o *Option) IsInline() bool {
	inline := true
	tmp := strings.Split(o.Raw, Keyword)
	if len(tmp) >= 1 {
		left := strings.TrimSpace(tmp[0])
		if len(left) == 0 {
			inline = false
		}
	}
	return inline
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
			if opt.Name == "lint" {
				continue
			}
			return "", fmt.Errorf("unknown directive '%s' in %s", opt.Name, opt.File)
		}
		profile, err = drtv.Apply(opt, profile)
		if err != nil {
			return "", fmt.Errorf("%s %s: %w", drtv.Name(), opt.File, err)
		}
	}
	return profile, nil
}
