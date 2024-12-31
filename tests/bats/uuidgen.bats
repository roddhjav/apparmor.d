#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load common

@test "uuidgen: Create a random UUIDv4" {
    uuidgen --random
}

@test "uuidgen: Create a UUIDv1 based on the current time" {
    uuidgen --time
}
