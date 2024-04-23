// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"embed"
	"fmt"
	"reflect"
	"strings"
	"text/template"
)

// Default indentation for apparmor profile (2 spaces)
const indentation = "  "

var (
	//go:embed templates/*.j2
	tmplFiles embed.FS

	// The functions available in the template
	tmplFunctionMap = template.FuncMap{
		"typeof":     typeOf,
		"join":       join,
		"indent":     indent,
		"overindent": indentDbus,
	}

	// The apparmor templates
	tmpl = map[string]*template.Template{
		"apparmor": generateTemplate("apparmor.j2"),
	}

	// convert apparmor requested mask to apparmor access mode
	maskToAccess = map[string]string{
		"a": "w",
		"c": "w",
		"d": "w",
	}

	// The order the apparmor rules should be sorted
	ruleAlphabet = []string{
		"include",
		"rlimit",
		"capability",
		"network",
		"mount",
		"remount",
		"umount",
		"pivotroot",
		"changeprofile",
		"mqueue",
		"signal",
		"ptrace",
		"unix",
		"userns",
		"iouring",
		"dbus",
		"file",
		"profile",
		"include_if_exists",
	}
	ruleWeights = map[string]int{}

	// The order the apparmor file rules should be sorted
	fileAlphabet = []string{
		"@{exec_path}",        // 1. entry point
		"@{sh_path}",          // 2.1 shells
		"@{bin}",              // 2.1 binaries
		"@{lib}",              // 2.2 libraries
		"/opt",                // 2.3 opt binaries & libraries
		"/usr/share",          // 3. shared data
		"/etc",                // 4. system configuration
		"/var",                // 5.1 system read/write data
		"/boot",               // 5.2 boot files
		"/home",               // 6.1 user data
		"@{HOME}",             // 6.2 home files
		"@{user_cache_dirs}",  // 7.1 user caches
		"@{user_config_dirs}", // 7.2 user config
		"@{user_share_dirs}",  // 7.3 user shared
		"/tmp",                // 8.1 Temporary data
		"@{run}",              // 8.2 Runtime data
		"/dev/shm",            // 8.3 Shared memory
		"@{sys}",              // 9. Sys files
		"@{PROC}",             // 10. Proc files
		"/dev",                // 11. Dev files
		"deny",                // 12. Deny rules
	}
	fileWeights = map[string]int{}
)

func generateTemplate(name string) *template.Template {
	res := template.New(name).Funcs(tmplFunctionMap)
	switch name {
	case "apparmor.j2":
		res = template.Must(res.ParseFS(tmplFiles,
			"templates/*.j2", "templates/rule/*.j2",
		))
	case "profile.j2":
		res = template.Must(res.Parse("{{ template \"profile\" . }}"))
		res = template.Must(res.ParseFS(tmplFiles,
			"templates/profile.j2", "templates/rule/*.j2",
		))
	default:
		res = template.Must(res.Parse(
			fmt.Sprintf("{{ template \"%s\" . }}", name),
		))
		res = template.Must(res.ParseFS(tmplFiles,
			fmt.Sprintf("templates/rule/%s.j2", name),
			"templates/rule/qualifier.j2", "templates/rule/comment.j2",
		))
	}
	return res
}

func init() {
	for i, r := range fileAlphabet {
		fileWeights[r] = i
	}
	for i, r := range ruleAlphabet {
		ruleWeights[r] = i
	}
}

func join(i any) string {
	switch reflect.TypeOf(i).Kind() {
	case reflect.Slice:
		return strings.Join(i.([]string), " ")
	case reflect.Map:
		res := []string{}
		for k, v := range i.(map[string]string) {
			res = append(res, k+"="+v)
		}
		return strings.Join(res, " ")
	default:
		return i.(string)
	}
}

func typeOf(i any) string {
	return strings.TrimPrefix(reflect.TypeOf(i).String(), "*aa.")
}

func typeToValue(i reflect.Type) string {
	return strings.ToLower(strings.TrimPrefix(i.String(), "*aa."))
}

func indent(s string) string {
	return indentation + s
}

func indentDbus(s string) string {
	return indentation + "     " + s
}

func getLetterIn(alphabet []string, in string) string {
	for _, letter := range alphabet {
		if strings.HasPrefix(in, letter) {
			return letter
		}
	}
	return ""
}

// Helper function to convert a access string to slice of access
func toAccess(constraint string, input string) []string {
	var res []string

	switch constraint {
	case "file", "file-log":
		raw := strings.Split(input, "")
		trans := []string{}
		for _, access := range raw {
			if slices.Contains(fileAccess, access) {
				res = append(res, access)
			} else if maskToAccess[access] != "" {
				res = append(res, maskToAccess[access])
				trans = append(trans, access)
			}
		}

		if constraint != "file-log" {
			transition := strings.Join(trans, "")
			if len(transition) > 0 {
				if slices.Contains(fileExecTransition, transition) {
					res = append(res, transition)
				} else {
					panic("unrecognized pattern: " + transition)
				}
			}
		}
		return res

	default:
		res = strings.Fields(input)
		slices.Sort(res)
		return slices.Compact(res)
	}
}
