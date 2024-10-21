#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load common

setup_file() {
    aa_setup
}

# bats test_tags=flatpak
@test "flatpak: List installed applications, ignoring runtimes" {
    flatpak list --app
    aa_check
}

# bats test_tags=flatpak
@test "flatpak: Install an application from a remote source" {
    flatpak install --noninteractive org.vim.Vim
    aa_check
}

# bats test_tags=flatpak
@test "flatpak: Show information about an installed application" {
    flatpak info org.vim.Vim
    aa_check
}

# bats test_tags=flatpak
@test "flatpak: Run an installed application" {
    flatpak run org.vim.Vim
    aa_check
}

# bats test_tags=flatpak
@test "flatpak: Update all installed applications and runtimes" {
    flatpak update --noninteractive
    aa_check
}

# bats test_tags=flatpak
@test "flatpak: Remove an installed application" {
    flatpak remove --noninteractive org.vim.Vim
    aa_check
}

# bats test_tags=flatpak
@test "flatpak: Remove all unused applications" {
    flatpak remove --unused
    aa_check
}
