#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load ../common

@test "systemd-cgls: Display the whole control group hierarchy on your system" {
    systemd-cgls --no-pager
}

@test "systemd-cgls: Display a control group tree of a specific resource controller" {
    systemd-cgls --no-pager io
}

@test "systemd-cgls: Display the control group hierarchy of one or more systemd units" {
    systemd-cgls --no-pager --unit systemd-logind
}
