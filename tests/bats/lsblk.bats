#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load common

setup_file() {
    aa_setup
}

# bats test_tags=lsblk
@test "lsblk: List all storage devices in a tree-like format" {
    lsblk
    aa_check
}

# bats test_tags=lsblk
@test "lsblk: Also list empty devices" {
    lsblk -a
    aa_check
}

# bats test_tags=lsblk
@test "lsblk: Print the SIZE column in bytes rather than in a human-readable format" {
    lsblk -b
    aa_check
}

# bats test_tags=lsblk
@test "lsblk: Output info about filesystems" {
    lsblk -f
    aa_check
}

# bats test_tags=lsblk
@test "lsblk: Use ASCII characters for tree formatting" {
    lsblk -i
    aa_check
}

# bats test_tags=lsblk
@test "lsblk: Output info about block-device topology" {
    lsblk -t
    aa_check
}

# bats test_tags=lsblk
@test "lsblk: Exclude the devices specified by the comma-separated list of major device numbers" {
    lsblk -e 1
    aa_check
}

# bats test_tags=lsblk
@test "lsblk: Display a customized summary using a comma-separated list of columns" {
    lsblk --output NAME,SERIAL,MODEL,TRAN,TYPE,SIZE,FSTYPE,MOUNTPOINT
    aa_check
}
