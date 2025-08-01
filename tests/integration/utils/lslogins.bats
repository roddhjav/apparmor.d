#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2025 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load ../common

@test "lslogins: Display users in the system" {
    lslogins
    sudo lslogins
}

@test "lslogins: Display user accounts" {
    lslogins --user-accs
}

@test "lslogins: Display last logins" {
    lslogins --last
}

@test "lslogins: Display system accounts" {
    lslogins --system-accs
}

@test "lslogins: Display supplementary groups" {
    lslogins --supp-groups
}
