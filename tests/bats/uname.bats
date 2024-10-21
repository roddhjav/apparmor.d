#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load common

setup_file() {
    aa_setup
}

# bats test_tags=uname
@test "uname: Print all information" {
    uname --all
    aa_check
}

# bats test_tags=uname
@test "uname: Print the current kernel name" {
    uname --kernel-name
    aa_check
}

# bats test_tags=uname
@test "uname: Print the current network node host name" {
    uname --nodename
    aa_check
}

# bats test_tags=uname
@test "uname: Print the current kernel release" {
    uname --kernel-release
    aa_check
}

# bats test_tags=uname
@test "uname: Print the current kernel version" {
    uname --kernel-version
    aa_check
}

# bats test_tags=uname
@test "uname: Print the current machine hardware name" {
    uname --machine
    aa_check
}

# bats test_tags=uname
@test "uname: Print the current processor type" {
    uname --processor
    aa_check
}

# bats test_tags=uname
@test "uname: Print the current operating system name" {
    uname --operating-system
    aa_check
}

