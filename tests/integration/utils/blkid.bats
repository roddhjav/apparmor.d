#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load ../common

@test "blkid: List all partitions" {
    sudo blkid
}

@test "blkid: List all partitions in a table, including current mountpoints" {
    sudo blkid -o list
}
