#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load common

@test "id: Display current user&#39;s ID (UID), group ID (GID) and groups to which they belong" {
    id
}

@test "id: Display the current user identity" {
    id -un
}

@test "id: Display the current user identity as a number" {
    id -u
}

@test "id: Display the current primary group identity" {
    id -gn
}

@test "id: Display the current primary group identity as a number" {
    id -g
}

@test "id: Display an arbitrary user ID (UID), group ID (GID) and groups to which they belong" {
    id root
}
