#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load common

setup_file() {
    skip "mqueue raised despite the rule being present. See https://gitlab.com/apparmor/apparmor/-/issues/362"
}

@test "needrestart: List outdated processes" {
    needrestart
}

@test "needrestart: Interactively restart services" {
    sudo needrestart
}

@test "needrestart: List outdated processes in verbose mode" {
    needrestart -v
}

@test "needrestart: Check if the kernel is outdated" {
    needrestart -k
}

@test "needrestart: Check if the CPU microcode is outdated" {
    needrestart -w
}

@test "needrestart: List outdated processes in batch mode" {
    needrestart -b
}

@test "needrestart: Display help" {
    needrestart --help
}
