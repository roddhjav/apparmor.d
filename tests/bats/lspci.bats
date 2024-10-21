#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load common

setup_file() {
    aa_setup
}

# bats test_tags=lspci
@test "lspci: Show a brief list of devices" {
    lspci
    aa_check
}

# bats test_tags=lspci
@test "lspci: Display additional info" {
    lspci -v
    aa_check
}

# bats test_tags=lspci
@test "lspci: Display drivers and modules handling each device" {
    lspci -k
    aa_check
}

# bats test_tags=lspci
@test "lspci: Show a specific device" {
    lspci -s 00:00.0
    aa_check
}

# bats test_tags=lspci
@test "lspci: Dump info in a readable form" {
    lspci -vm
    aa_check
}
