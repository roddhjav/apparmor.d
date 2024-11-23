#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load common

@test "upower: Display power and battery information" {
    upower --dump
}

@test "upower: List all power devices" {
    upower --enumerate
}

@test "upower: Display version" {
    upower --version
}

