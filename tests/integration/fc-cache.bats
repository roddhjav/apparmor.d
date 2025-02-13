#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load common

@test "fc-cache: Generate font cache files" {
    fc-cache
}

@test "fc-cache: Force a rebuild of all font cache files, without checking if cache is up-to-date" {
    fc-cache -f
}

@test "fc-cache: Erase font cache files, then generate new font cache files" {
    fc-cache -r
}
