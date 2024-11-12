#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load common

setup_file() {
    aa_setup
}

# bats test_tags=upower
@test "upower: Display power and battery information" {
    upower --dump
    aa_check
}

# bats test_tags=upower
@test "upower: List all power devices" {
    upower --enumerate
    aa_check
}

# bats test_tags=upower
@test "upower: Display version" {
    upower --version
    aa_check
}

