#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load common

@test "ip-address: List network interfaces and their associated IP addresses" {
    ip address
}

@test "ip-address: Filter to show only active network interfaces" {
    ip address show up
}

@test "ip-route: Display the routing table" {
    ip route
}

@test "ip-route-get: Print route to a destination" {
    ip route get 1.1.1.1
}

@test "ip link: Show information about all network interfaces" {
    ip link
}

@test "ip neighbour: Display the neighbour/ARP table entries" {
    ip neighbour
}

@test "ip rule: Display the routing policy" {
    ip rule show
    ip rule list
}

@test "ip: Manage network namespace" {
    sudo ip netns add foo
    sudo ip netns list
    sudo ip netns exec foo bash -c "pwd"
    sudo ip netns delete foo
}
