#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load ../common

@test "w: Display information about all users who are currently logged in" {
    w
}

@test "w: Display information about a specific user" {
    w root
}

@test "w: Display information without including the header, the login, JCPU and PCPU columns" {
    w --no-header
    w --short
}
