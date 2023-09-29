// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2023 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"embed"
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

	// The apparmor profile template
	tmplAppArmorProfile = generateTemplate()

	// convert apparmor requested mask to apparmor access mode
	// TODO: Should be a map of slice, not exhausive yet
	maskToAccess = map[string]string{
		"a":            "w",
		"c":            "w",
		"d":            "w",
		"k":            "rk",
		"l":            "l",
		"m":            "rm",
		"r":            "r",
		"ra":           "rw",
		"read write":   "read write",
		"read":         "read",
		"readby":       "readby",
		"receive":      "receive",
		"rm":           "rm",
		"rw":           "rw",
		"send receive": "send receive",
		"send":         "send",
		"w":            "w",
		"wc":           "w",
		"wd":           "w",
		"wk":           "wk",
		"wr":           "rw",
		"wrc":          "rw",
		"wrd":          "rw",
		"write":        "write",
		"x":            "rix",
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
		"include_if_exists",
	}
	ruleWeights = map[string]int{}

	// The order the apparmor file rules should be sorted
	fileAlphabet = []string{
		"@{exec_path}",        // 1. entry point
		"@{bin}",              // 2.1 binaries
		"@{lib}",              // 2.2 libraries
		"/opt",                // 2.3 opt binaries & libraries
		"/usr/share",          // 3. shared data
		"/etc",                // 4. system configuration
		"/",                   // 5.1 system data
		"/var",                // 5.2 system data read/write data
		"/boot",               // 5.3 boot files
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

func generateTemplate() *template.Template {
	res := template.New("profile.j2").Funcs(tmplFunctionMap)
	res = template.Must(res.ParseFS(tmplFiles, "templates/*.j2"))
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

func indent(s string) string {
	return indentation + s
}

func indentDbus(s string) string {
	return indentation + "     " + s
}
