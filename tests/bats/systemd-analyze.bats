#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load common

setup_file() {
    aa_setup
}

# bats test_tags=systemd-analyze
@test "systemd-analyze: List all running units, ordered by the time they took to initialize" {
    systemd-analyze --no-pager blame
    aa_check
}

# bats test_tags=systemd-analyze
@test "systemd-analyze: Print a tree of the time-critical chain of units" {
    systemd-analyze --no-pager critical-chain
    aa_check
}

# bats test_tags=systemd-analyze
@test "systemd-analyze: Show security scores of running units" {
    systemd-analyze --no-pager security
    aa_check
}

