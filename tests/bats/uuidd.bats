#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load common

setup_file() {
    aa_setup
}

# bats test_tags=uuidd
@test "uuidd: Generate a random UUID" {
    uuidd --random
    aa_check
}

# bats test_tags=uuidd
@test "uuidd: Generate a bulk number of random UUIDs" {
    uuidd --random --uuids 10
    aa_check
}

# bats test_tags=uuidd
@test "uuidd: Generate a time-based UUID, based on the current time and MAC address of the system" {
    uuidd --time
    aa_check
}

