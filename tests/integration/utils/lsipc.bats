#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2025 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load ../common

@test "lsipc: Show information about all active IPC facilities" {
    lsipc
}

@test "lsipc: Show information about active shared memory segments, message queues or sempahore sets" {
    lsipc --shmems
    lsipc --queues
    lsipc --semaphores
}
