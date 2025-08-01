#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load ../common

@test "lsns: List all namespaces" {
    lsns
    sudo lsns
}

@test "lsns: List namespaces in JSON format" {
    sudo lsns --json
}

@test "lsns: List namespaces associated with the specified process" {
    sudo lsns --task 1
}

@test "lsns: List the specified type of namespaces only" {
    sudo lsns --type mnt
    sudo lsns --type net
    sudo lsns --type ipc
    sudo lsns --type user
    sudo lsns --type pid
    sudo lsns --type uts
    sudo lsns --type cgroup
    sudo lsns --type time
}

