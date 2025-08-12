#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2025 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load ../common

@test "pgrep: Return PIDs of any running processes with a matching command string" {
    pgrep systemd
}

@test "pgrep: Search for processes including their command-line options" {
    pgrep --full 'systemd'
}

@test "pgrep: Search for processes run by a specific user" {
    pgrep --euid root systemd-udevd
}

