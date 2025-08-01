#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load ../common

@test "apt: Update the list of available packages and versions" {
    sudo apt update
}

@test "apt: Search for a given package" {
    apt search apparmor
}

@test "apt: Show information for a package" {
    apt show apparmor
}

@test "apt: Install a package, or update it to the latest available version" {
    sudo apt install -y pass
}

@test "apt: Remove a package and its configuration files" {
    sudo apt purge -y pass
}

@test "apt: Upgrade all installed packages to their newest available versions" {
    sudo apt upgrade -y
}

@test "apt: Upgrade installed packages, but remove obsolete packages and install additional packages to meet new dependencies" {
    sudo apt dist-upgrade -y
}

@test "apt: Clean the local repository - removing package files (.deb) from interrupted downloads that can no longer be downloaded" {
    sudo apt autoclean -y
}

@test "apt: Remove all packages that are no longer needed" {
    sudo apt autoremove -y
}

@test "apt: List all packages" {
    apt list
}

@test "apt: List installed packages" {
    apt list --installed
}

@test "apt: Print a cow easter egg" {
    apt moo
}
