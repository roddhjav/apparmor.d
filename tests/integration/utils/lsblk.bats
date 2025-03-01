#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load ../common

@test "lsblk: List all storage devices in a tree-like format" {
    lsblk
}

@test "lsblk: Also list empty devices" {
    lsblk -a
}

@test "lsblk: Print the SIZE column in bytes rather than in a human-readable format" {
    lsblk -b
}

@test "lsblk: Output info about filesystems" {
    lsblk -f
}

@test "lsblk: Use ASCII characters for tree formatting" {
    lsblk -i
}

@test "lsblk: Output info about block-device topology" {
    lsblk -t
}

@test "lsblk: Exclude the devices specified by the comma-separated list of major device numbers" {
    lsblk -e 1
}

@test "lsblk: Display a customized summary using a comma-separated list of columns" {
    lsblk --output NAME,SERIAL,MODEL,TRAN,TYPE,SIZE,FSTYPE,MOUNTPOINT
}
