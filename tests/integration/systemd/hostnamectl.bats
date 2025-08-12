#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load ../common

@test "hostnamectl: Get the hostname of the computer" {
    hostnamectl
}

@test "hostnamectl: Get the location of the computer" {
    hostnamectl location
}

@test "hostnamectl: Set the hostname of the computer" {
    name=$(hostnamectl hostname)
    sudo hostnamectl set-hostname "new"
    sudo hostnamectl set-hostname "$name"
}
