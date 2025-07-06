#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load ../common

@test "free: Display system memory" {
    free
}

@test "free: Display memory in GB" {
    free -g
}

@test "free: Display memory in human-readable units" {
    free -h
}
