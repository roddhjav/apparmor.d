#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load ../common

@test "fstrim: Trim unused blocks on all mounted partitions that support it" {
    sudo fstrim --all
}

@test "fstrim: Trim unused blocks on a specified partition" {
    sudo fstrim --verbose /
}
