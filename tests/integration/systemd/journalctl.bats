#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2025 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load ../common

@test "journalctl: Show all messages with priority level 3 (errors) from this boot" {
    sudo journalctl -b --priority=3
}

@test "journalctl: Show only the last N lines of the journal" {
    sudo journalctl --lines 100
}

@test "journalctl: Show all messages by a specific [u]nit" {
    sudo journalctl --unit apparmor.service
}

@test "journalctl: Show all messages by a specific process" {
    sudo journalctl _PID=1
}

@test "journalctl: Show all messages by a specific executable" {
    sudo journalctl /usr/bin/bootctl
}

@test "journalctl: Delete journal logs which are older than 10 seconds" {
    sudo journalctl --vacuum-time=10s
}
