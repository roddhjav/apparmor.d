#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load common

setup_file() {
    aa_setup
}

# bats test_tags=dfc
@test "dfc: Display filesystems and their disk usage in human-readable form with colors and graphs" {
    dfc
    aa_check
}

# bats test_tags=dfc
@test "dfc: Display all filesystems including pseudo, duplicate and inaccessible filesystems" {
    dfc -a
    aa_check
}

# bats test_tags=dfc
@test "dfc: Display filesystems without color" {
    dfc -c never
    aa_check
}

# bats test_tags=dfc
@test "dfc: Display filesystems containing &#34;ext&#34; in the filesystem type" {
    dfc -t ext
    aa_check
}
