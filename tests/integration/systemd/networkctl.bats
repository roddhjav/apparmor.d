#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2025 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load ../common

@test "networkctl: List existing links with their status" {
    sudo networkctl list
}

@test "networkctl: Show an overall network status" {
    sudo networkctl status
}

@test "networkctl: Reload configuration files (.netdev and .network)" {
    sudo networkctl reload
}
