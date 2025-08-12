#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load ../common

@test "ps: List all running processes" {
    ps aux
}

@test "ps: List all running processes including the full command string" {
    ps auxww
}

@test "ps: List all processes of the current user in extra full format" {
    ps --user "$(id -u)" -F
}

@test "ps: List all processes of the current user as a tree" {
    ps --user "$(id -u)" -f
}

@test "ps: Get the parent PID of a process" {
    ps -o ppid= -p 1
}

@test "ps: Sort processes by memory consumption" {
    ps auxww --sort size
}
