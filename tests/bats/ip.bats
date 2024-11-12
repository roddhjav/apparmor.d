#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load common

setup_file() {
    aa_setup
}

# bats test_tags=ip
@test "ip: List interfaces with detailed info" {
    ip address
    aa_check
}

# bats test_tags=ip
@test "ip: List interfaces with brief link layer info" {
    ip link
    aa_check
}

# bats test_tags=ip
@test "ip: Display the routing table" {
    ip route
    aa_check
}

# bats test_tags=ip
@test "ip: Show neighbors (ARP table)" {
    ip neighbour
    aa_check
}

# bats test_tags=ip
@test "ip: Manage network namespace" {
    sudo ip netns add foo
    sudo ip netns list
    sudo ip netns exec foo bash -c "pwd"
    sudo ip netns delete foo
    aa_check
}


