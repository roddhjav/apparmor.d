#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load common

@test "useradd: Create a new user with the specified shell" {
    sudo useradd --shell /bin/bash --create-home user2
}

@test "useradd: Create a new user with the specified user ID" {
    sudo useradd --uid 3000 user3
}

@test "useradd: Create a new user belonging to additional groups (mind the lack of whitespace)" {
    sudo useradd --groups adm user4
}

@test "useradd: Create a new system user without the home directory" {
    sudo useradd --system sys2
}

@test "userdel: Remove a user" {
    sudo userdel user3
    sudo userdel user4
    sudo userdel sys2
}

@test "userdel: Remove a user along with the home directory and mail spool" {
    sudo userdel --remove user2
}
