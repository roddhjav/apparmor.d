#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load common

@test "snap: Search for a package" {
    snap find vim
}

@test "snap: Install a package" {
    sudo snap install nano-strict
}

@test "snap: Update a package to another channel (track, risk, or branch)" {
    sudo snap refresh nano-strict --channel=edge
}

@test "snap: Update all packages" {
    sudo snap refresh
}

@test "snap: Display basic information about installed snap software" {
    sudo snap list
}

@test "snap: Check for recent snap changes in the system" {
    sudo snap changes
}

@test "snap: Uninstall a package" {
    sudo snap remove nano-strict
}
