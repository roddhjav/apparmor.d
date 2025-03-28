#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load ../common

@test "eject: Display the default device" {
    eject -d || true
}

@test "eject: Eject the default device" {
    eject || true
}
