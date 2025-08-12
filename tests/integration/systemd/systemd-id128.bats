#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load ../common

@test "systemd-id128: Generate a new random identifier" {
    systemd-id128 new
}

@test "systemd-id128: Print the identifier of the current machine" {
    systemd-id128 machine-id
}

@test "systemd-id128: Print the identifier of the current boot" {
    systemd-id128 boot-id
}

@test "systemd-id128: Generate a new random identifier and print it as a UUID (five groups of digits separated by hyphens)" {
    systemd-id128 new --uuid
}

