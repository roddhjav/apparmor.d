#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load common

@test "hostname: Show current host name" {
    hostname
}

@test "hostname: Show the network address of the host name" {
    hostname -i
}

@test "hostname: Show all network addresses of the host" {
    hostname -I
}

@test "hostname: Show the FQDN (Fully Qualified Domain Name)" {
    hostname --fqdn
}

@test "hostname: Set current host name" {
    name=$(hostname)
    sudo hostname "new-$(hostname)"
    sudo hostname "$name"
}

