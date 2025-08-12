#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load common

@test "snap: Search for a package" {
    snap find vim
}

@test "snap: Install a package" {
    sudo snap install vault
}

@test "snap: Update a package to another channel (track, risk, or branch)" {
    sudo snap refresh vault --channel=edge
}

@test "snap: Update all packages" {
    sudo snap refresh
}

@test "snap: Display basic information about installed snap software" {
    sudo snap list
}

@test "snap: lists information about the services" {
    sudo snap services
    sudo snap services vault
}

@test "snap: starts, and optionally enables, the given services" {
    sudo snap start --enable vault
}

@test "snap: logs of the given services" {
    sudo snap logs vault || true
}

@test "snap: restarts the given services" {
    sudo snap restart vault
}

@test "snap: stops, and optionally disables, the given services" {
    sudo snap stop --disable vault
}

@test "snap: Uninstall a package" {
    sudo snap remove vault
}

@test "snap: Check for recent snap changes in the system" {
    sudo snap changes
}
