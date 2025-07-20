#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load ../common

@test "uptime: Print current time, uptime, number of logged-in users and other information" {
    uptime
}

@test "uptime: Show only the amount of time the system has been booted for" {
    uptime --pretty
}

@test "uptime: Print the date and time the system booted up at" {
    uptime --since
}
