#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load common

@test "cpuid: Display information for all CPUs" {
    cpuid
}

@test "cpuid: Display information only for the current CPU" {
    cpuid -1
}

@test "cpuid: Display raw hex information with no decoding" {
    cpuid -r
}
