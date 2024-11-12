#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load common

setup_file() {
    aa_setup
}

# bats test_tags=sync
@test "sync: Flush all pending write operations on all disks" {
    sync
    aa_check
}

# bats test_tags=sync
@test "sync: Flush all pending write operations on a single file to disk" {
    sudo sync /
    aa_check
}
