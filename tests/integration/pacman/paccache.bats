#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load ../common

@test "paccache: Perform a dry-run and show the number of candidate packages for deletion" {
    sudo paccache -d
}

@test "paccache: Move candidate packages to a directory instead of deleting them" {
    sudo paccache -m "$USER_BUILD_DIRS"
}

@test "paccache: Remove all but the 3 most recent package versions from the `pacman` cache" {
    sudo paccache -r
}

@test "paccache: Set the number of package versions to keep" {
    sudo paccache -rk 3
}
