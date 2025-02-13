#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load common

@test "fwupdmgr: Display all devices detected by fwupd" {
    fwupdmgr get-devices
}

@test "fwupdmgr: Download the latest firmware metadata from LVFS" {
    fwupdmgr refresh || true
}

@test "fwupdmgr: List the updates available for devices on your system" {
    fwupdmgr get-updates || true
}

@test "fwupdmgr: Install firmware updates" {
    fwupdmgr update || true
}

