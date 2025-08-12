#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load ../common

@test "lslocks: List all local system locks" {
    sudo lslocks
}

@test "lslocks: List locks producing a raw output (no columns), and without column headers" {
    sudo lslocks --raw --noheadings
}

@test "lslocks: List locks by PID input" {
    sudo lslocks --pid "$(sudo lslocks --raw --noheadings  --output PID | head -1)"
}

@test "lslocks: List locks with JSON output to stdout" {
    lslocks --json
}
