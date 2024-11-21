#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load common

@test "fc-list: Return a list of installed fonts in your system" {
    fc-list
}

@test "fc-match: Return a sorted list of best matching fonts" {
    fc-match -s 'DejaVu Serif'
}

@test "fc-pattern: Display default information about a font" {
    fc-pattern --default 'DejaVu Serif'
}

@test "fc-pattern: Display configuration information about a font" {
    fc-pattern --config 'DejaVu Serif'
}
