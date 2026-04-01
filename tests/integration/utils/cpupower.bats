#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2026 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load ../common

@test "cpupower: List CPUs" {
    sudo cpupower --cpu all info || true
}

@test "cpupower: Print information about all cores" {
    sudo cpupower --cpu all info || true
}

@test "cpupower: Set all CPUs to a power-saving frequency governor" {
    sudo cpupower --cpu all frequency-set --governor powersave || true
}

@test "cpupower: Print CPU 0's available frequency governors" {
    sudo cpupower --cpu 0 frequency-info --governors || true
}

@test "cpupower: Print CPU 4's frequency from the hardware, in a human-readable format" {
    sudo cpupower --cpu 0 frequency-info --hwfreq --human || true
}
