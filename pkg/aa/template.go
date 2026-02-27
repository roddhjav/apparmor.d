// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"embed"
	"fmt"
	"strings"
	"text/template"
)

var (
	// Default indentation for apparmor profile (2 spaces)
	Indentation = "  "

	// The current indentation level
	IndentationLevel = 0

	//go:embed templates/*.gotmpl
	//go:embed templates/rule/*.gotmpl
	tmplFiles embed.FS

	// The functions available in the template
	tmplFunctionMap = template.FuncMap{
		"kindof":     kindOf,
		"join":       join,
		"cjoin":      cjoin,
		"indent":     indent,
		"overindent": indentDbus,
		"setindent":  setindent,
	}

	// The apparmor templates
	tmpl = generateTemplates([]Kind{
		// Global templates
		"apparmor", PROFILE, HAT, "rules",

		// Preamble templates
		ABI, ALIAS, INCLUDE, VARIABLE, COMMENT,

		// Rules templates
		ALL, RLIMIT, USERNS, CAPABILITY, NETWORK,
		MOUNT, REMOUNT, UMOUNT, PIVOTROOT, CHANGEPROFILE,
		MQUEUE, IOURING, UNIX, PTRACE, SIGNAL, DBUS,
		FILE, LINK,
	})

	// convert apparmor requested mask to apparmor access mode
	maskToAccess = map[string]string{
		"a":  "w",
		"c":  "w",
		"d":  "w",
		"wc": "w",
		"x":  "ix",
	}

	// The order the apparmor rules should be sorted
	ruleAlphabet = []Kind{
		INCLUDE,
		ALL,
		RLIMIT,
		USERNS,
		CAPABILITY,
		NETWORK,
		MOUNT,
		REMOUNT,
		UMOUNT,
		PIVOTROOT,
		CHANGEPROFILE,
		MQUEUE,
		IOURING,
		SIGNAL,
		PTRACE,
		UNIX,
		DBUS,
		FILE,
		LINK,
		PROFILE,
		HAT,
		"include_if_exists",
	}
	ruleWeights = generateWeights(ruleAlphabet)

	// The order the apparmor file rules should be sorted
	fileAlphabet = []string{
		"@{exec_path}",        // 1. entry point
		"@{sh_path}",          // 2.1 shells
		"@{coreutils_path}",   // 2.2 coreutils
		"@{open_path}",        // 2.3 binaries paths
		"@{bin}",              // 2.3 binaries
		"@{lib}",              // 2.4 libraries
		"/opt",                // 2.5 opt binaries & libraries
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
		"@{tmp}",              // 8.1. User temporary data
		"/dev/shm",            // 8.2 Shared memory
		"@{run}",              // 8.3 Runtime data
		"@{sys}",              // 9. Sys files
		"@{PROC}",             // 10. Proc files
		"/dev",                // 11. Dev files
		"deny",                // 12. Deny rules
		"profile",             // 13. Subprofiles
	}
	fileWeights = generateWeights(fileAlphabet)

	// Some file rule should be sorted in the same group
	fileAlphabetGroups = map[string]string{
		"@{exec_path}":      "exec",
		"@{sh_path}":        "exec",
		"@{coreutils_path}": "exec",
		"@{open_path}":      "exec",
		"@{bin}":            "exec",
		"@{lib}":            "exec",
		"/opt":              "exec",
		"/home":             "home",
		"@{HOME}":           "home",
		"/tmp":              "tmp",
		"@{tmp}":            "tmp",
		"/dev/shm":          "tmp",
	}

	// The order AARE should be sorted
	stringAlphabet = []byte(
		"!\"#$%&'*(){}[]@+,-./:;<=>?\\^_`|~0123456789abcdefghijklmnopqrstuvwxyz",
	)
	stringWeights = generateWeights(stringAlphabet)

	// The order the rule values (access, type, domains, etc) should be sorted
	requirements        = map[Kind]requirement{}
	requirementsWeights map[Kind]map[string]map[string]int

	// Pairs of mutually exclusive values that cannot coexist
	conflicts = map[Kind]map[string][][]string{}
)

func init() {
	requirementsWeights = generateRequirementsWeights(requirements)
}

func generateTemplates(names []Kind) map[Kind]*template.Template {
	res := make(map[Kind]*template.Template, len(names))
	base := template.New("").Funcs(tmplFunctionMap)
	base = template.Must(base.ParseFS(tmplFiles,
		"templates/*.gotmpl", "templates/rule/*.gotmpl",
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

func renderTemplate(name Kind, data any) string {
	var res strings.Builder
	template, ok := tmpl[name]
	if !ok {
		panic("template '" + name.String() + "' not found")
	}
	err := template.Execute(&res, data)
	if err != nil {
		panic(err)
	}
	return res.String()
}

func generateWeights[T comparable](alphabet []T) map[T]int {
	res := make(map[T]int, len(alphabet))
	for i, r := range alphabet {
		res[r] = i
	}
	return res
}

func generateRequirementsWeights(requirements map[Kind]requirement) map[Kind]map[string]map[string]int {
	res := make(map[Kind]map[string]map[string]int, len(requirements))
	for rule, req := range requirements {
		res[rule] = make(map[string]map[string]int, len(req))
		for key, values := range req {
			res[rule][key] = generateWeights(values)
		}
	}
	return res
}

func join(i any) string {
	switch i := i.(type) {
	case []string:
		return strings.Join(i, " ")
	case map[string]string:
		res := []string{}
		for k, v := range i {
			res = append(res, k+"="+v)
		}
		return strings.Join(res, " ")
	default:
		return i.(string)
	}
}

func cjoin(i any) string {
	switch i := i.(type) {
	case []string:
		if len(i) == 1 {
			return i[0]
		}
		return "(" + strings.Join(i, " ") + ")"
	case map[string]string:
		res := []string{}
		for k, v := range i {
			res = append(res, k+"="+v)
		}
		return "(" + strings.Join(res, " ") + ")"
	default:
		return i.(string)
	}
}

func kindOf(i Rule) string {
	if i == nil {
		return ""
	}
	return i.Kind().String()
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
