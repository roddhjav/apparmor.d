#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load common

setup_file() {
    aa_setup
}

# bats test_tags=fc-cache
@test "fc-cache: Generate font cache files" {
    fc-cache
    aa_check
}

# bats test_tags=fc-cache
@test "fc-cache: Force a rebuild of all font cache files, without checking if cache is up-to-date" {
    fc-cache -f
    aa_check
}

# bats test_tags=fc-cache
@test "fc-cache: Erase font cache files, then generate new font cache files" {
    fc-cache -r
    aa_check
}

