#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2025 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load ../common

@test "lsfd: List all open file descriptors" {
    lsfd
}

@test "lsfd: List all files kept open by a specific program" {
    sudo lsfd --filter 'PID == 1'
}

@test "lsfd: List open IPv4 or IPv6 sockets" {
    sudo lsfd -i4
    sudo lsfd -i6
}
