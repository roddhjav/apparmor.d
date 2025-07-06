#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load ../common

@test "pidof: List all process IDs with given name" {
    pidof systemd
    pidof bash
}

@test "pidof: List a single process ID with given name" {
    pidof -s bash
}

@test "pidof: List process IDs including scripts with given name" {
    pidof -x bash
}
