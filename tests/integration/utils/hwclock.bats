#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load ../common

@test "hwclock: Display the current time as reported by the hardware clock" {
    sudo hwclock || true
}

@test "hwclock: Write the current software clock time to the hardware clock (sometimes used during system setup)" {
    sudo hwclock --systohc || true
}

@test "hwclock: Write the current hardware clock time to the software clock" {
    sudo hwclock --hctosys || true
}

