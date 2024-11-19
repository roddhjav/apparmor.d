#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load common

setup_file() {
    aa_setup
    skip
}

# bats test_tags=homectl
@test "homectl: Display help" {
    homectl --no-pager --help
    aa_check
}

# bats test_tags=homectl
@test "homectl: Create a user account and their associated home directory" {
    sudo homectl create user2
    aa_check
}

# bats test_tags=homectl
@test "homectl: List user accounts and their associated home directories" {
    homectl list
    aa_check
}

# bats test_tags=homectl
@test "homectl: Change the password for a specific user" {
    sudo homectl passwd user2
    aa_check
}

# bats test_tags=homectl
@test "homectl: Run a shell or a command with access to a specific home directory" {
    sudo homectl with user2 -- ls -al /home/user2
    aa_check
}

# bats test_tags=homectl
@test "homectl: Lock or unlock a specific home directory" {
    sudo homectl lock user2
    aa_check
}

# bats test_tags=homectl
@test "homectl: Change the disk space assigned to a specific home directory to 100 GiB" {
    sudo homectl resize user2 1G
    aa_check
}

# bats test_tags=homectl
@test "homectl: Remove a specific user and the associated home directory" {
    sudo homectl remove user2
    aa_check
}
