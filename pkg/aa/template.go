// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2023 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	_ "embed"
	"strings"
)

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

