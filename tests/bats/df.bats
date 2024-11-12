#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load common

setup_file() {
    aa_setup
}

# bats test_tags=df
@test "df: Display all filesystems and their disk usage" {
    df
    aa_check
}

# bats test_tags=df
@test "df: Display all filesystems and their disk usage in human-readable form" {
    df -h
    aa_check
}

# bats test_tags=df
@test "df: Display the filesystem and its disk usage containing the given file or directory" {
    df apparmor.d/
    aa_check
}

# bats test_tags=df
@test "df: Include statistics on the number of free inodes" {
    df --inodes
    aa_check
}

# bats test_tags=df
@test "df: Display filesystem types" {
    df --print-type
    aa_check
}
