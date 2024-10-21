#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load common

setup_file() {
    aa_setup
}

# bats test_tags=id
@test "id: Display current user&#39;s ID (UID), group ID (GID) and groups to which they belong" {
    id
    aa_check
}

# bats test_tags=id
@test "id: Display the current user identity" {
    id -un
    aa_check
}

# bats test_tags=id
@test "id: Display the current user identity as a number" {
    id -u
    aa_check
}

# bats test_tags=id
@test "id: Display the current primary group identity" {
    id -gn
    aa_check
}

# bats test_tags=id
@test "id: Display the current primary group identity as a number" {
    id -g
    aa_check
}

# bats test_tags=id
@test "id: Display an arbitrary user ID (UID), group ID (GID) and groups to which they belong" {
    id root
}
