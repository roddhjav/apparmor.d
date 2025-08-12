#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load ../common

@test "systemd-analyze: List all running units, ordered by the time they took to initialize" {
    systemd-analyze --no-pager blame
}

@test "systemd-analyze: Print a tree of the time-critical chain of units" {
    systemd-analyze --no-pager critical-chain
}

@test "systemd-analyze: Show security scores of running units" {
    systemd-analyze --no-pager security
}
