#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load common

@test "dfc: Display filesystems and their disk usage in human-readable form with colors and graphs" {
    dfc
}

@test "dfc: Display all filesystems including pseudo, duplicate and inaccessible filesystems" {
    dfc -a
}

@test "dfc: Display filesystems without color" {
    dfc -c never
}

@test "dfc: Display filesystems containing &#34;ext&#34; in the filesystem type" {
    dfc -t ext
}
