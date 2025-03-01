#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load ../common

@test "lsusb: List all the USB devices available" {
    lsusb || true
}

@test "lsusb: List the USB hierarchy as a tree" {
    lsusb -t || true
}

@test "lsusb: List verbose information about USB devices" {
    lsusb --verbose || true
}
