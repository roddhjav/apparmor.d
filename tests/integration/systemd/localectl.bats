#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2025 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load ../common

@test "localectl: Show the current settings of the system locale and keyboard mapping" {
    localectl
}

@test "localectl: List available locales" {
    localectl list-locales
}

@test "localectl: Set a system locale variable" {
    sudo localectl set-locale LANG=en_US.UTF-8
}

@test "localectl: List available keymaps" {
    localectl list-keymaps || true
}

@test "localectl: Set the system keyboard mapping for the console and X11" {
    sudo localectl set-keymap uk || true
}

