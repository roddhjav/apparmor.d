#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load common

@test "tlp: Apply settings (according to the actual power source)" {
    sudo tlp start
}

@test "tlp: Apply battery settings (ignoring the actual power source)" {
    sudo tlp bat
}

@test "tlp: Apply AC settings (ignoring the actual power source)" {
    sudo tlp ac
}

@test "tlp: Apply Disk settings" {
    sudo tlp diskid
}
