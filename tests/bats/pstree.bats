#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load common

setup_file() {
    aa_setup
}

# bats test_tags=pstree
@test "pstree: Display a tree of processes" {
    pstree
    aa_check
}

# bats test_tags=pstree
@test "pstree: Display a tree of processes with PIDs" {
    pstree -p
    aa_check
}

# bats test_tags=pstree
@test "pstree: Display all process trees rooted at processes owned by specified user" {
    pstree root
    aa_check
}

