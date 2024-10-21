#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load common

setup_file() {
    aa_setup
}

# bats test_tags=aa-enforce
@test "aa-enforce: Disable profile" {
    sudo aa-disable pass
    aa_check
}

# bats test_tags=aa-enforce
@test "aa-enforce: Enforce a profile" {
    sudo aa-enforce pass
    aa_check
}

# bats test_tags=aa-enforce
@test "aa-enforce: Complain a profile" {
    sudo aa-complain pass
    aa_check
}

# bats test_tags=aa-enforce
@test "aa-enforce: Audit a profile" {
    sudo aa-audit pass
    aa_check
}
