// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2023 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

// Warning: this is purposely not using a Yacc parser. Its only aim is to
// extract variables and attachments for apparmor.d profile

package aa

import (
	"regexp"
	"strings"

	"github.com/arduino/go-paths-helper"
)

var (
	regVariablesDef = regexp.MustCompile(`@{(.*)}\s*[+=]+\s*(.*)`)
	regVariablesRef = regexp.MustCompile(`@{([^{}]+)}`)
)

// Default Apparmor magic directory: /etc/apparmor.d/.
var MagicRoot = paths.New("/etc/apparmor.d")

// DefaultTunables return a minimal working profile to build the profile
// It should not be used when loading file from /etc/apparmor.d
func DefaultTunables() *AppArmorProfile {
	return &AppArmorProfile{
		Preamble: Preamble{
			Variables: []Variable{
				{"bin", []string{"/{,usr/}{,s}bin"}},
				{"lib", []string{"/{,usr/}lib{,exec,32,64}"}},
				{"multiarch", []string{"*-linux-gnu*"}},
				{"user_share_dirs", []string{"/home/*/.local/share"}},
				{"etc_ro", []string{"/{,usr/}etc/"}},
				{"int", []string{"[0-9]{[0-9],}{[0-9],}{[0-9],}{[0-9],}{[0-9],}{[0-9],}{[0-9],}{[0-9],}{[0-9],}"}},
			},
		},
	}
}

// ParseVariables extract all variables from the profile
func (p *AppArmorProfile) ParseVariables(content string) {
	matches := regVariablesDef.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		if len(match) > 2 {
			key := match[1]
			values := strings.Split(match[2], " ")
			found := false
			for idx, variable := range p.Variables {
				if variable.Name == key {
					p.Variables[idx].Values = append(p.Variables[idx].Values, values...)
					found = true
					break
				}
			}
			if !found {
				variable := Variable{Name: key, Values: values}
				p.Variables = append(p.Variables, variable)
			}
		}
	}
}

// resolve recursively resolves all variables references
func (p *AppArmorProfile) resolve(str string) []string {
	if strings.Contains(str, "@{") {
		vars := []string{}
		match := regVariablesRef.FindStringSubmatch(str)
		if len(match) > 1 {
			variable := match[0]
			varname := match[1]
			for _, vrbl := range p.Variables {
				if vrbl.Name == varname {
					for _, value := range vrbl.Values {
						newVar := strings.ReplaceAll(str, variable, value)
						vars = append(vars, p.resolve(newVar)...)
					}
				}
			}
		} else {
			vars = append(vars, str)
		}
		return vars
	}
	return []string{str}
}

// ResolveAttachments resolve profile attachments defined in exec_path
func (p *AppArmorProfile) ResolveAttachments() {
	for _, variable := range p.Variables {
		if variable.Name == "exec_path" {
			for _, value := range variable.Values {
				p.Attachments = append(p.Attachments, p.resolve(value)...)
			}
		}
	}
}

// NestAttachments return a nested attachment string
func (p *AppArmorProfile) NestAttachments() string {
	if len(p.Attachments) == 0 {
		return ""
	} else if len(p.Attachments) == 1 {
		return p.Attachments[0]
	} else {
		res := []string{}
		for _, attachment := range p.Attachments {
			if strings.HasPrefix(attachment, "/") {
				res = append(res, attachment[1:])
			} else {
				res = append(res, attachment)
			}
		}
		return "/{" + strings.Join(res, ",") + "}"
	}
}

