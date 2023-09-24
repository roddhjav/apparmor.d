// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2023 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	_ "embed"
	"reflect"
	"strings"
	"text/template"
)

// Default indentation for apparmor profile (2 spaces)
const indentation = "  "

var (
	//go:embed template.j2
	tmplFileAppArmorProfile string

	// tmplFunctionMap is the list of function available in the template
	tmplFunctionMap = template.FuncMap{
		"typeof":     typeOf,
		"join":       join,
		"indent":     indent,
		"overindent": indentDbus,
	}

	// The apparmor profile template
	tmplAppArmorProfile = template.Must(template.New("profile").
				Funcs(tmplFunctionMap).Parse(tmplFileAppArmorProfile))

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

)


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
