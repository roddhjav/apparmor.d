#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load common

setup_file() {
    aa_setup
}

# bats test_tags=blkid
@test "blkid: List all partitions" {
    sudo blkid
    aa_check
}

# bats test_tags=blkid
@test "blkid: List all partitions in a table, including current mountpoints" {
    sudo blkid -o list
    aa_check
}
