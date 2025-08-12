#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load ../common

@test "sync: Flush all pending write operations on all disks" {
    sync
}

@test "sync: Flush all pending write operations on a single file to disk" {
    sudo sync /
}
