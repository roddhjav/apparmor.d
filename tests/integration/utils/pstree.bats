#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load ../common

@test "pstree: Display a tree of processes" {
    pstree
}

@test "pstree: Display a tree of processes with PIDs" {
    pstree -p
}

@test "pstree: Display all process trees rooted at processes owned by specified user" {
    pstree root
}

