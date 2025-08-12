#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2025 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load ../common

@test "dpkg: Remove a package" {
    sudo apt install -y pass
    sudo dpkg -r pass
}

@test "dpkg: List installed packages" {
    dpkg -l apparmor
}

@test "dpkg: List a package's contents" {
    dpkg -L apparmor.d
}

@test "dpkg: Find out which package owns a file" {
    dpkg -S /etc/apparmor/parser.conf
}

@test "dpkg: Purge an installed or already removed package, including configuration" {
    sudo dpkg -P pass
}
