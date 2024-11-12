#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load common

setup_file() {
    aa_setup
}

# bats test_tags=systemd-ac-power
@test "systemd-ac-power: Report whether we are connected to an external power source." {
    systemd-ac-power || true
    aa_check
}

# bats test_tags=systemd-ac-power
@test "systemd-ac-power: Check if battery is discharging and low" {
    systemd-ac-power --low || true
    aa_check
}

