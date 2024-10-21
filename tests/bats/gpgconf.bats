#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load common

setup_file() {
    aa_setup
}

# bats test_tags=gpgconf
@test "gpgconf: List all components" {
    gpgconf --list-components
    aa_check
}

# bats test_tags=gpgconf
@test "gpgconf: List the directories used by gpgconf" {
    gpgconf --list-dirs
    aa_check
}

# bats test_tags=gpgconf
@test "gpgconf: List all options of a component" {
    gpgconf --list-options gpg
    gpgconf --list-options gpgsm
    gpgconf --list-options gpg-agent
    gpgconf --list-options scdaemon
    gpgconf --list-options dirmngr
    aa_check
}

