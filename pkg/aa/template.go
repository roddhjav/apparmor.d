// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2023 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	_ "embed"
	"text/template"
)

const indentation = "  "

//go:embed template.j2
var tmplFileAppArmorProfile string

var tmplFunctionMap = template.FuncMap{
	"indent":     indent,
	"overindent": indentDbus,
}

var tmplAppArmorProfile = template.Must(template.New("profile").
	Funcs(tmplFunctionMap).Parse(tmplFileAppArmorProfile))

func indent(s string) string {
	return indentation + s
}

func indentDbus(s string) string {
	return indentation + "     " + s
}

// TODO: Should be a map of slice, not exhausive yet
var maskToAccess = map[string]string{
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
	"wr":           "rw",
	"wrc":          "rw",
	"wrd":          "rw",
	"write":        "write",
	"x":            "rix",
}

