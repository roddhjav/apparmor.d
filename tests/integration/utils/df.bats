#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load ../common

@test "df: Display all filesystems and their disk usage" {
    df
}

@test "df: Display all filesystems and their disk usage in human-readable form" {
    df -h
}

@test "df: Display the filesystem and its disk usage containing the given file or directory" {
    df /etc/apparmor.d/
}

@test "df: Include statistics on the number of free inodes" {
    df --inodes
}

@test "df: Display filesystem types" {
    df --print-type
}
