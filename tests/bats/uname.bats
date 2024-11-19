#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load common

@test "uname: Print all information" {
    uname --all
}

@test "uname: Print the current kernel name" {
    uname --kernel-name
}

@test "uname: Print the current network node host name" {
    uname --nodename
}

@test "uname: Print the current kernel release" {
    uname --kernel-release
}

@test "uname: Print the current kernel version" {
    uname --kernel-version
}

@test "uname: Print the current machine hardware name" {
    uname --machine
}

@test "uname: Print the current processor type" {
    uname --processor
}

@test "uname: Print the current operating system name" {
    uname --operating-system
}

