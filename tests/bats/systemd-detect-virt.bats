#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

# bats test_tags=systemd-detect-virt
@test "systemd-detect-virt: List detectable virtualization technologies" {
    systemd-detect-virt --list
}

# bats test_tags=systemd-detect-virt
@test "systemd-detect-virt: Detect virtualization, print the result and return a zero status code when running in a VM or a container, and a non-zero code otherwise" {
    systemd-detect-virt
}

# bats test_tags=systemd-detect-virt
@test "systemd-detect-virt: Silently check without printing anything" {
    systemd-detect-virt --quiet
}

# bats test_tags=systemd-detect-virt
@test "systemd-detect-virt: Only detect hardware virtualization" {
    systemd-detect-virt --vm
}

