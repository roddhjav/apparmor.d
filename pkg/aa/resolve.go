// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	regVariableReference = regexp.MustCompile(`@{([^{}]+)}`)
)

// Resolve resolves all variables and includes in the profile and merge the rules in the profile
func (f *AppArmorProfileFile) Resolve() error {
	// Resolve variables
	for _, variable := range f.Preamble.GetVariables() {
		newValues := []string{}
		for _, value := range variable.Values {
			vars := f.resolveVariable(value)
			if len(vars) == 0 {
				return fmt.Errorf("Variable not defined in: %s", value)
			}
			newValues = append(newValues, vars...)
		}
		variable.Values = newValues
	}

	// Resolve variables in attachements
	for _, profile := range f.Profiles {
		attachments := []string{}
		for _, att := range profile.Attachments {
			vars := f.resolveVariable(att)
			if len(vars) == 0 {
				return fmt.Errorf("Variable not defined in: %s", att)
			}
			attachments = append(attachments, vars...)
		}
		profile.Attachments = attachments
	}

	return nil
}

func (f *AppArmorProfileFile) resolveVariable(input string) []string {
	if !strings.Contains(input, tokVARIABLE) {
		return []string{input}
	}

	vars := []string{}
	match := regVariableReference.FindStringSubmatch(input)
	if len(match) > 1 {
		variable := match[0]
		varname := match[1]
		for _, vrbl := range f.Preamble.GetVariables() {
			if vrbl.Name == varname {
				for _, v := range vrbl.Values {
					newVar := strings.ReplaceAll(input, variable, v)
					res := f.resolveVariable(newVar)
					vars = append(vars, res...)
				}
			}
		}
	} else {
		vars = append(vars, input)
	}
	return vars
}
