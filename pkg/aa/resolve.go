// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/roddhjav/apparmor.d/pkg/paths"
)

var (
	includeCache map[*Include]*AppArmorProfileFile = make(map[*Include]*AppArmorProfileFile)

	regVariableReference = regexp.MustCompile(`@{([^{}]+)}`)
)

// Resolve resolves variables and includes definied in the profile preamble
func (f *AppArmorProfileFile) Resolve() error {
	// Resolve preamble includes
	// for _, include := range f.Preamble.GetIncludes() {
	// 	err := f.resolveInclude(include)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// Append value to variable
	seen := map[string]*Variable{}
	for idx, variable := range f.Preamble.GetVariables() {
		if _, ok := seen[variable.Name]; ok {
			if variable.Define {
				return fmt.Errorf("variable %s already defined", variable.Name)
			}
			seen[variable.Name].Values = append(seen[variable.Name].Values, variable.Values...)
			f.Preamble = f.Preamble.Delete(idx)
		}
		if variable.Define {
			seen[variable.Name] = variable
		}
	}

	// Resolve variables
	for _, variable := range f.Preamble.GetVariables() {
		newValues := []string{}
		for _, value := range variable.Values {
			vars, err := f.resolveValues(value)
			if err != nil {
				return err
			}
			newValues = append(newValues, vars...)
		}
		variable.Values = newValues
	}

	// Resolve variables in attachements
	for _, profile := range f.Profiles {
		attachments := []string{}
		for _, att := range profile.Attachments {
			vars, err := f.resolveValues(att)
			if err != nil {
				return err
			}
			attachments = append(attachments, vars...)
		}
		profile.Attachments = attachments
	}

	return nil
}

func (f *AppArmorProfileFile) resolveValues(input string) ([]string, error) {
	if !strings.Contains(input, VARIABLE.Tok()) {
		return []string{input}, nil
	}

	values := []string{}
	match := regVariableReference.FindStringSubmatch(input)
	if len(match) == 0 {
		return nil, fmt.Errorf("invalid variable reference: %s", input)
	}

	variable := match[0]
	varname := match[1]
	found := false
	for _, vrbl := range f.Preamble.GetVariables() {
		if vrbl.Name == varname {
			found = true
			for _, v := range vrbl.Values {
				if strings.Contains(v, VARIABLE.Tok()+varname+"}") {
					return nil, fmt.Errorf("recursive variable found in: %s", varname)
				}
				newValues := strings.ReplaceAll(input, variable, v)
				newValues = strings.ReplaceAll(newValues, "//", "/")
				res, err := f.resolveValues(newValues)
				if err != nil {
					return nil, err
				}
				values = append(values, res...)
			}
		}
	}

	if !found {
		return nil, fmt.Errorf("variable %s not defined", varname)
	}
	return values, nil
}

// resolveInclude resolves all includes defined in the profile preamble
func (f *AppArmorProfileFile) resolveInclude(include *Include) error {
	if include == nil || include.Path == "" {
		return fmt.Errorf("invalid include: %v", include)
	}

	_, isCached := includeCache[include]
	if !isCached {
		var files paths.PathList
		var err error

		path := MagicRoot.Join(include.Path)
		if !include.IsMagic {
			path = paths.New(include.Path)
		}

		if path.IsDir() {
			files, err = path.ReadDir(paths.FilterOutDirectories())
			if err != nil {
				if include.IfExists {
					return nil
				}
				return fmt.Errorf("File %s not found: %v", path, err)
			}

		} else if path.Exist() {
			files = append(files, path)

		} else {
			if include.IfExists {
				return nil
			}
			return fmt.Errorf("File %s not found", path)

		}

		iFile := &AppArmorProfileFile{}
		for _, file := range files {
			raw, err := file.ReadFileAsString()
			if err != nil {
				return err
			}
			if _, err := iFile.Parse(raw); err != nil {
				return err
			}
		}
		if err := iFile.Validate(); err != nil {
			return err
		}
		for _, inc := range iFile.Preamble.GetIncludes() {
			if err := iFile.resolveInclude(inc); err != nil {
				return err
			}
		}

		// Remove all includes in iFile
		iFile.Preamble = iFile.Preamble.DeleteKind(INCLUDE)

		// Cache the included file
		includeCache[include] = iFile
	}

	// Insert iFile in the place of include in the current file
	index := f.Preamble.Index(include)
	f.Preamble = f.Preamble.Replace(index, includeCache[include].Preamble...)
	return nil
}
