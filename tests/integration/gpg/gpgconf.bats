#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load ../common

@test "gpgconf: List all components" {
    gpgconf --list-components
}

@test "gpgconf: List the directories used by gpgconf" {
    gpgconf --list-dirs
}

@test "gpgconf: List all options of a component" {
    gpgconf --list-options gpg
    gpgconf --list-options gpgsm
    gpgconf --list-options gpg-agent
    gpgconf --list-options scdaemon || true
    gpgconf --list-options dirmngr
}

@test "gpgconf: List programs and test whether they are runnable" {
    gpgconf --check-programs || true
}

@test "gpgconf: Reload a component" {
    gpgconf --reload gpg
    gpgconf --reload gpgsm
    gpgconf --reload gpg-agent
    gpgconf --reload scdaemon || true
    gpgconf --reload dirmngr
}
