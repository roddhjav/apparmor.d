// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"embed"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"text/template"
)

var (
	// Default indentation for apparmor profile (2 spaces)
	Indentation = "  "

	// The current indentation level
	IndentationLevel = 0

	//go:embed templates/*.j2
	//go:embed templates/rule/*.j2
	tmplFiles embed.FS

	// The functions available in the template
	tmplFunctionMap = template.FuncMap{
		"typeof":     typeOf,
		"join":       join,
		"cjoin":      cjoin,
		"indent":     indent,
		"overindent": indentDbus,
		"setindent":  setindent,
	}

	// The apparmor templates
	tmpl = generateTemplates([]string{
		"apparmor", tokPROFILE, "rules", // Global templates
		tokINCLUDE, tokRLIMIT, tokCAPABILITY, tokNETWORK,
		tokMOUNT, tokPIVOTROOT, tokCHANGEPROFILE, tokSIGNAL,
		tokPTRACE, tokUNIX, tokUSERNS, tokIOURING,
		tokDBUS, "file", "variable",
	})

	// convert apparmor requested mask to apparmor access mode
	maskToAccess = map[string]string{
		"a": "w",
		"c": "w",
		"d": "w",
	}

	// The order the apparmor rules should be sorted
	ruleAlphabet = []string{
		"include",
		"all",
		"rlimit",
		"userns",
		"capability",
		"network",
		"mount",
		"remount",
		"umount",
		"pivotroot",
		"changeprofile",
		"mqueue",
		"iouring",
		"signal",
		"ptrace",
		"unix",
		"dbus",
		"file",
		"profile",
		"include_if_exists",
	}
	ruleWeights = generateWeights(ruleAlphabet)

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
		"profile",             // 13. Subprofiles
	}
	fileWeights = generateWeights(fileAlphabet)

	// The order the rule values (access, type, domains, etc) should be sorted
	requirements        = map[string]requirement{}
	requirementsWeights map[string]map[string]map[string]int
)

func init() {
	requirementsWeights = generateRequirementsWeights(requirements)
}

func generateTemplates(names []string) map[string]*template.Template {
	res := make(map[string]*template.Template, len(names))
	base := template.New("").Funcs(tmplFunctionMap)
	base = template.Must(base.ParseFS(tmplFiles,
		"templates/*.j2", "templates/rule/*.j2",
	))
	for _, name := range names {
		t := template.Must(base.Clone())
		t = template.Must(t.Parse(
			fmt.Sprintf(`{{- template "%s" . -}}`, name),
		))
		res[name] = t
	}
	return res
}

func renderTemplate(name string, data any) string {
	var res strings.Builder
	template, ok := tmpl[name]
	if !ok {
		panic("template '" + name + "' not found")
	}
	err := template.Execute(&res, data)
	if err != nil {
		panic(err)
	}
	return res.String()
}

func generateWeights(alphabet []string) map[string]int {
	res := make(map[string]int, len(alphabet))
	for i, r := range alphabet {
		res[r] = i
	}
	return res
}

func generateRequirementsWeights(requirements map[string]requirement) map[string]map[string]map[string]int {
	res := make(map[string]map[string]map[string]int, len(requirements))
	for rule, req := range requirements {
		res[rule] = make(map[string]map[string]int, len(req))
		for key, values := range req {
			res[rule][key] = generateWeights(values)
		}
	}
	return res
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

func cjoin(i any) string {
	switch reflect.TypeOf(i).Kind() {
	case reflect.Slice:
		s := i.([]string)
		if len(s) == 1 {
			return s[0]
		}
		return "(" + strings.Join(s, " ") + ")"
	case reflect.Map:
		res := []string{}
		for k, v := range i.(map[string]string) {
			res = append(res, k+"="+v)
		}
		return "(" + strings.Join(res, " ") + ")"
	default:
		return i.(string)
	}
}

func typeOf(i any) string {
	if i == nil {
		return ""
	}
	return strings.TrimPrefix(reflect.TypeOf(i).String(), "*aa.")
}

func typeToValue(i reflect.Type) string {
	return strings.ToLower(strings.TrimPrefix(i.String(), "*aa."))
}

func setindent(i string) string {
	switch i {
	case "++":
		IndentationLevel++
	case "--":
		IndentationLevel--
	}
	return ""
}

func indent(s string) string {
	return strings.Repeat(Indentation, IndentationLevel) + s
}

func indentDbus(s string) string {
	return strings.Join([]string{Indentation, s}, "     ")
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
