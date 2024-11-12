#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load common

setup_file() {
    aa_setup
}

# bats test_tags=userdbctl
@test "userdbctl: List all known user records" {
    userdbctl --no-pager user
    aa_check
}

# bats test_tags=userdbctl
@test "userdbctl: Show details of a specific user" {
    userdbctl --no-pager user "$USER"
    aa_check
}

# bats test_tags=userdbctl
@test "userdbctl: List all known groups" {
    userdbctl --no-pager group
    aa_check
}

# bats test_tags=userdbctl
@test "userdbctl: Show details of a specific group" {
    sudo userdbctl --no-pager group "$USER"
    aa_check
}

# bats test_tags=userdbctl
@test "userdbctl: List all services currently providing user/group definitions to the system" {
    userdbctl --no-pager services
    aa_check
}

