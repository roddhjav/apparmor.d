#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2025 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load ../common

@test "importctl: Import an image as a machine" {
    sudo importctl pull-tar --force --class=machine -N https://cloud-images.ubuntu.com/noble/current/noble-server-cloudimg-amd64-root.tar.xz noble || true
}

@test "machinectl: Display a list of available images" {
    sudo machinectl list-images
}

@test "machinectl: Start a machine as a service using systemd-nspawn" {
    sudo machinectl start noble || true
}

@test "machinectl: Display a list of running machines" {
    sudo machinectl list
}

@test "machinectl: Stop a running machine" {
    sudo machinectl stop noble || true
}
