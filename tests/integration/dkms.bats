#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2026 Łukasz Kozak <lakis.kozak@gmail.com>
# SPDX-License-Identifier: GPL-2.0-only

load common

setup_file() {
    if [[ $(sudo dkms status -m zfs | grep -c 'installed') -eq 0 ]]; then
        echo "dkms: No kernel has a ZFS module installed. Skipping" >&3
        skip
    fi

    _lsb_release() {
	    . /etc/os-release || exit 1
	    echo "$ID"
    }
    DISTRIBUTION="$(_lsb_release)"

    case "$DISTRIBUTION" in
    opensuse*)
            # Tumbleweed distributes ZFS only as precompiled modules.
            # DKMS package would've needed to be built in the preparation phase
            skip
        ;;
    *) ;;
    esac

    DKMS_LINE=$(sudo dkms status -m zfs | grep installed | head -n 1)

    export KERNEL=$(echo $DKMS_LINE | sed "s/, $(uname -m): installed.*//g" | sed 's/zfs\/.*, //g')
    export ZFS_DKMS=$(echo $DKMS_LINE | sed "s/, $KERNEL,.*//g") 
}

@test "dkms: Check DKMS status" {
    sudo dkms status
}

@test "dkms: Pack a module to tarball" {
    sudo dkms mktarball $ZFS_DKMS -k $KERNEL \
    --archive /var/lib/dkms/zfs/zfs.tar
}

@test "dkms: Remove a module" {
    sudo dkms remove $ZFS_DKMS -k $KERNEL
}

@test "dkms: Install a module" {
    sudo dkms install $ZFS_DKMS -k $KERNEL
}

teardown_file() {
    sudo rm -f /var/lib/dkms/zfs/zfs.tar
}