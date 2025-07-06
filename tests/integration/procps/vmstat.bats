#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load ../common

@test "vmstat: Display virtual memory statistics" {
    vmstat
    vmstat --active
    vmstat --forks
}

@test "vmstat: Display disk statistics" {
    vmstat --disk
    vmstat --disk-sum 
}

@test "vmstat: Display slabinfo" {
    sudo vmstat --slabs
}

@test "vmstat: Display reports every second for 3 times" {
    vmstat 1 3
}
