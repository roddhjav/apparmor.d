#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2025 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load ../common

@test "bootctl: Show information about the system firmware and the bootloaders" {
    sudo bootctl status
}

@test "bootctl: Show all available bootloader entries" {
    sudo bootctl list
}

@test "bootctl: Install 'systemd-boot' into the EFI system partition" {
    sudo bootctl install
}

@test "bootctl: Remove all installed versions of 'systemd-boot' from the EFI system partition" {
    sudo bootctl remove
}
