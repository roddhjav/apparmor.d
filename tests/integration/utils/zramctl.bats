#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load ../common

@test "zramctl: Check if zram is enabled; enable it if needed" {
    lsmod | grep -i zram || sudo modprobe zram || true
}

@test "zramctl: Find and initialize the next free zram device to a 1 GB virtual drive using LZ4 compression" {
    sudo zramctl --find --size 1GB --algorithm lz4 || true
}

@test "zramctl: List currently initialized devices" {
    sudo zramctl || true
}
