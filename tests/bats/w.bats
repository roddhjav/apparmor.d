#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load common

setup_file() {
    aa_setup
}

# bats test_tags=w
@test "w: Display information about all users who are currently logged in" {
    w
    aa_check
}

# bats test_tags=w
@test "w: Display information about a specific user" {
    w root
    aa_check
}
