#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load common

setup_file() {
    aa_setup
}

# bats test_tags=uuidgen
@test "uuidgen: Create a random UUIDv4" {
    uuidgen --random
    aa_check
}

# bats test_tags=uuidgen
@test "uuidgen: Create a UUIDv1 based on the current time" {
    uuidgen --time
    aa_check
}

