#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load common

setup_file() {
    aa_setup
}

# bats test_tags=snap
@test "snap: Search for a package" {
    snap find vim
    aa_check
}

# bats test_tags=snap
@test "snap: Install a package" {
    sudo snap install nano-strict
    aa_check
}

# bats test_tags=snap
@test "snap: Update a package to another channel (track, risk, or branch)" {
    sudo snap refresh nano-strict --channel=edge
    aa_check
}

# bats test_tags=snap
@test "snap: Update all packages" {
    sudo snap refresh
    aa_check
}

# bats test_tags=snap
@test "snap: Display basic information about installed snap software" {
    sudo snap list
    aa_check
}

# bats test_tags=snap
@test "snap: Check for recent snap changes in the system" {
    sudo snap changes
    aa_check
}

# bats test_tags=snap
@test "snap: Uninstall a package" {
    sudo snap remove nano-strict
    aa_check
}
