#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load ../common

@test "dmesg: Show kernel messages" {
    sudo dmesg
}

@test "dmesg: Show kernel error messages" {
    sudo dmesg --level err
}

@test "dmesg: Show how much physical memory is available on this system" {
    sudo dmesg | grep -i memory
}

@test "dmesg: Show kernel messages with a timestamp (available in kernels 3.5.0 and newer)" {
    sudo dmesg -T
}

@test "dmesg: Show kernel messages in human-readable form (available in kernels 3.5.0 and newer)" {
    sudo dmesg -H
}

@test "dmesg: Colorize output (available in kernels 3.5.0 and newer)" {
    sudo dmesg -L
}
