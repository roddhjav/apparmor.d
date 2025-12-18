#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load common

@test "flatpak: Add a new remote repository (by URL)" {
    sudo flatpak remote-add --if-not-exists flathub https://flathub.org/repo/flathub.flatpakrepo
}

@test "flatpak: List all remote repositories" {
    flatpak remotes
}

@test "flatpak: Search for an application in a remote repository" {
    sudo flatpak search vim
    sudo flatpak search org.freedesktop.Platform
}

@test "flatpak: Install an application from a remote source" {
    sudo flatpak install --noninteractive org.vim.Vim
}

@test "flatpak: List installed applications, ignoring runtimes" {
    flatpak list --app
}

@test "flatpak: Show information about an installed application" {
    flatpak info org.vim.Vim
}

@test "flatpak: List exported files" {
    flatpak documents
}

@test "flatpak: List dynamic permissions" {
    flatpak permissions
}

# @test "flatpak: Run an installed application" {
#     _timeout flatpak run org.vim.Vim
# }

@test "flatpak: Update all installed applications and runtimes" {
    sudo flatpak update --noninteractive
}

@test "flatpak: Remove an installed application" {
    sudo flatpak remove --noninteractive org.vim.Vim
}

@test "flatpak: Remove all unused applications" {
    sudo flatpak remove --noninteractive --unused
}
