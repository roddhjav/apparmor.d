#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load ../common

@test "groupadd: Create a new group" {
    sudo groupadd user2
}

@test "groupadd: Create a new system group" {
    sudo groupadd --system system2
}

@test "groupadd: Create a new group with the specific groupid" {
    sudo groupadd --gid 3000 user3
}

@test "groupmod: Change the group name" {
    sudo groupmod --new-name user22 user2
}

@test "groupmod: Change the group ID" {
    sudo groupmod --gid 2222 user22
}

@test "groupdel: Delete newly created group" {
    sudo groupdel user22
    sudo groupdel system2
    sudo groupdel user3
}
