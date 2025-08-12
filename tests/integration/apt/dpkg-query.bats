#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load ../common

@test "dpkg-query: List all installed packages" {
    dpkg-query --list
}

@test "dpkg-query: List installed packages matching a pattern" {
    dpkg-query --list 'libc6*'
}

@test "dpkg-query: List all files installed by a package" {
    dpkg-query --listfiles libc6
}

@test "dpkg-query: Show information about a package" {
    dpkg-query --status libc6
}

@test "dpkg-query: Search for packages that own files matching a pattern" {
    dpkg-query --search /etc/ld.so.conf.d
}

