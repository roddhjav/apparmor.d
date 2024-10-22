#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load common

setup_file() {
    aa_setup
}

# bats test_tags=lsusb
@test "lsusb: List all the USB devices available" {
    lsusb
    aa_check
}

# bats test_tags=lsusb
@test "lsusb: List the USB hierarchy as a tree" {
    lsusb -t
    aa_check
}

# bats test_tags=lsusb
@test "lsusb: List verbose information about USB devices" {
    lsusb --verbose
    aa_check
}
