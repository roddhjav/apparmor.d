#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load common

setup() {
    aa_setup
}

# bats test_tags=hostnamectl
@test "hostnamectl: Get the hostname of the computer" {
    hostnamectl
}

# bats test_tags=hostnamectl
@test "hostnamectl: Get the location of the computer" {
    hostnamectl location
}

# bats test_tags=hostnamectl
@test "hostnamectl: Set the hostname of the computer" {
    name=$(hostnamectl hostname)
    sudo hostnamectl set-hostname "new"
    sudo hostnamectl set-hostname "$name"
}
