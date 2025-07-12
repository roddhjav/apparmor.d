#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2025 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load common

@test "whois: Get information about a domain name" {
    whois google.fr
}

@test "whois: Get information about an IP address" {
    whois 8.8.8.8
}

@test "whois: Get abuse contact for an IP address" {
    whois -b 8.8.8.8
}

