#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load common

setup_file() {
    aa_setup
}

# bats test_tags=who
@test "who: Display the username, line, and time of all currently logged-in sessions" {
    who
    aa_check
}

# bats test_tags=who
@test "who: Display all available information" {
    who -a
    aa_check
}

# bats test_tags=who
@test "who: Display all available information with table headers" {
    who -a -H
    aa_check
}

