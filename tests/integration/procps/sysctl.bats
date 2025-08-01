#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load ../common

@test "sysctl: Show all available variables and their values" {
    sysctl -a
}

@test "sysctl: Set a changeable kernel state variable" {
    sudo sysctl -w vm.panic_on_oom=0
}

@test "sysctl: Get currently open file handlers" {
    sysctl fs.file-nr
}

@test "sysctl: Get limit for simultaneous open files" {
    sysctl fs.file-max
}

@test "sysctl: Apply changes from '/etc/sysctl.conf'" {
    sudo sysctl -p
}
