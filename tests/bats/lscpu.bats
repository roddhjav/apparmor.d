#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load common

setup_file() {
    aa_setup
}

# bats test_tags=lscpu
@test "lscpu: Display information about all CPUs" {
    lscpu
    aa_check
}

# bats test_tags=lscpu
@test "lscpu: Display information in a table" {
    lscpu --extended
    aa_check
}

# bats test_tags=lscpu
@test "lscpu: Display only information about offline CPUs in a table" {
    lscpu --extended --offline
    aa_check
}
