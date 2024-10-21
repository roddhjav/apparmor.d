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
    gpgconf --list-options scdaemon || true
    gpgconf --list-options dirmngr
    aa_check
}

# bats test_tags=gpgconf
@test "gpgconf: List programs and test whether they are runnable" {
    gpgconf --check-programs || true
    aa_check
}

# bats test_tags=gpgconf
@test "gpgconf: Reload a component" {
    gpgconf --reload gpg
    gpgconf --reload gpgsm
    gpgconf --reload gpg-agent
    gpgconf --reload scdaemon || true
    gpgconf --reload dirmngr
    aa_check
}
