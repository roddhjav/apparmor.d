#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load ../common

@test "pacman-key: Initialize the 'pacman' keyring" {
    sudo pacman-key --init
}

@test "pacman-key: Add the default Arch Linux keys" {
    sudo pacman-key --populate
}

@test "pacman-key: List keys from the public keyring" {
    pacman-key --list-keys
}

@test "pacman-key: Receive a key from a key server" {
    sudo pacman-key --recv-keys 06A26D531D56C42D66805049C5469996F0DF68EC
}

@test "pacman-key: Print the fingerprint of a specific key" {
    pacman-key --finger 06A26D531D56C42D66805049C5469996F0DF68EC
}

@test "pacman-key: Sign an imported key locally" {
    sudo pacman-key --lsign-key 06A26D531D56C42D66805049C5469996F0DF68EC
}

@test "pacman-key: Remove a specific key" {
    sudo pacman-key --delete 06A26D531D56C42D66805049C5469996F0DF68EC
}
