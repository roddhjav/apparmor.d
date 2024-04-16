// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

// Warning: this is purposely not using a Yacc parser. Its only aim is to
// extract variables and attachments for apparmor.d profile

package aa

import (
	"regexp"
	"strings"
)

var (
	regVariablesDef = regexp.MustCompile(`@{(.*)}\s*[+=]+\s*(.*)`)
	regVariablesRef = regexp.MustCompile(`@{([^{}]+)}`)
)

// DefaultTunables return a minimal working profile to build the profile
// It should not be used when loading file from /etc/apparmor.d
func DefaultTunables() *AppArmorProfileFile {
	return &AppArmorProfileFile{
		Preamble: Preamble{
			Variables: []*Variable{
				{Name: "bin", Values: []string{"/{,usr/}{,s}bin"}},
				{Name: "lib", Values: []string{"/{,usr/}lib{,exec,32,64}"}},
				{Name: "multiarch", Values: []string{"*-linux-gnu*"}},
				{Name: "HOME", Values: []string{"/home/*"}},
				{Name: "user_share_dirs", Values: []string{"/home/*/.local/share"}},
				{Name: "etc_ro", Values: []string{"/{,usr/}etc/"}},
				{Name: "int", Values: []string{"[0-9]{[0-9],}{[0-9],}{[0-9],}{[0-9],}{[0-9],}{[0-9],}{[0-9],}{[0-9],}{[0-9],}"}},
			},
		},
	}
}

// ParseVariables extract all variables from the profile
func (f *AppArmorProfileFile) ParseVariables(content string) {
	matches := regVariablesDef.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		if len(match) > 2 {
			key := match[1]
			values := strings.Split(match[2], " ")
			found := false
			for idx, variable := range f.Variables {
				if variable.Name == key {
					f.Variables[idx].Values = append(f.Variables[idx].Values, values...)
					found = true
					break
				}
			}
			if !found {
				variable := &Variable{Name: key, Values: values}
				f.Variables = append(f.Variables, variable)
			}
		}
	}
}

// resolve recursively resolves all variables references
func (f *AppArmorProfileFile) resolve(str string) []string {
	if strings.Contains(str, "@{") {
		vars := []string{}
		match := regVariablesRef.FindStringSubmatch(str)
		if len(match) > 1 {
			variable := match[0]
			varname := match[1]
			for _, vrbl := range f.Variables {
				if vrbl.Name == varname {
					for _, value := range vrbl.Values {
						newVar := strings.ReplaceAll(str, variable, value)
						vars = append(vars, f.resolve(newVar)...)
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
func (f *AppArmorProfileFile) ResolveAttachments() {
	p := f.GetDefaultProfile()

	for _, variable := range profile.Variables {
		if variable.Name == "exec_path" {
			for _, value := range variable.Values {
				attachments := profile.resolve(value)
				if len(attachments) == 0 {
					panic("Variable not defined in: " + value)
				}
				p.Attachments = append(p.Attachments, attachments...)
			}
		}
	}
}

// NestAttachments return a nested attachment string
func (f *AppArmorProfileFile) NestAttachments() string {
	p := f.GetDefaultProfile()
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
