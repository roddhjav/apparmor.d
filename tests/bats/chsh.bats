#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load common

setup_file() {
    aa_setup
}

# bats test_tags=chsh
@test "chsh: [l]ist available shells" {
    chsh --list-shells
    aa_check
}

# bats test_tags=chsh
@test "chsh: Set a specific login [s]hell for the current user" {
    chsh --shell /usr/bin/bash
    aa_check
}

# bats test_tags=chsh
@test "chsh: Set a login [s]hell for a specific user" {
    sudo chsh --shell /usr/bin/sh root
    aa_check
}
