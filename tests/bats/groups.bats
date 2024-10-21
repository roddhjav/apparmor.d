#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load common

setup_file() {
    aa_setup
}

# bats test_tags=groups
@test "groups: Print group memberships for the current user" {
    groups
    aa_check
}

# bats test_tags=groups
@test "groups: Print group memberships for a list of users" {
    groups root
    aa_check
}

