#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load ../common

setup_file() {
    skip
}

@test "aa-enforce: Disable profile" {
    sudo aa-disable pass
}

@test "aa-enforce: Enforce a profile" {
    sudo aa-enforce pass
}

@test "aa-enforce: Complain a profile" {
    sudo aa-complain pass
}

@test "aa-enforce: Audit a profile" {
    sudo aa-audit pass
}
