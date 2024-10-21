#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load common

setup_file() {
    aa_setup
}

# bats test_tags=groupadd
@test "groupadd: Create a new group" {
    sudo groupadd user2
    aa_check
}

# bats test_tags=groupadd
@test "groupadd: Create a new system group" {
    sudo groupadd --system system2
    aa_check
}

# bats test_tags=groupadd
@test "groupadd: Create a new group with the specific groupid" {
    sudo groupadd --gid 3000 user3
    aa_check
}

# bats test_tags=groupadd
@test "groupdel: Delete newly created group" {
    sudo groupdel user2
    sudo groupdel system2
    sudo groupdel user3
    aa_check
}
