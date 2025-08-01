#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load ../common

@test "dpkg-reconfigure: Reconfigure one or more packages" {
    sudo apt install -y pass
    sudo dpkg-reconfigure pass
}

