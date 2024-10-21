#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load common

setup_file() {
    aa_setup
}

# bats test_tags=aa-status
@test "aa-status: Check status" {
    sudo aa-status
    aa_check
}

# bats test_tags=aa-status
@test "aa-status: Display the number of loaded policies" {
    sudo aa-status --profiled
    aa_check
}

# bats test_tags=aa-status
@test "aa-status: Display the number of loaded enforicing policies" {
    sudo aa-status --enforced
    aa_check
}

# bats test_tags=aa-status
@test "aa-status: Display the number of loaded non-enforcing policies" {
    sudo aa-status --complaining
    aa_check
}

# bats test_tags=aa-status
@test "aa-status: Display the number of loaded enforcing policies that kill tasks" {
    sudo aa-status --kill
    aa_check
}
