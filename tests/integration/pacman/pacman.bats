#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load ../common

@test "pacman: Synchronize and update all packages" {
    sudo pacman -Syu --noconfirm
}

@test "pacman: Install a new package" {
    sudo pacman -S --noconfirm pass pass-otp
}

@test "pacman: Remove a package and its dependencies" {
    sudo pacman -Rs --noconfirm pass-otp
}

@test "pacman: List installed packages and versions" {
    pacman -Q
}

@test "pacman: List only the explicitly installed packages and versions" {
    pacman -Qe
}

@test "pacman: List orphan packages (installed as dependencies but not actually required by any package)" {
    pacman -Qtdq
}

@test "pacman: Empty the entire 'pacman' cache" {
    sudo pacman -Scc --noconfirm
}
