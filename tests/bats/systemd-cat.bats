#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

load common

setup_file() {
    aa_setup
}

# bats test_tags=systemd-cat
@test "systemd-cat: Write the output of the specified command to the journal (both output streams are captured)" {
    systemd-cat pwd
    aa_check
}

# bats test_tags=systemd-cat
@test "systemd-cat: Write the output of a pipeline to the journal (`stderr` stays connected to the terminal)" {
    echo apparmor.d-test-suite | systemd-cat
    aa_check
}
