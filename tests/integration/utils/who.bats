#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load ../common

@test "who: Display the username, line, and time of all currently logged-in sessions" {
    who
}

@test "who: Display all available information" {
    who -a
}

@test "who: Display all available information with table headers" {
    who -a -H
}

