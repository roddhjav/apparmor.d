#!/usr/bin/env bats
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2026 Łukasz Kozak <lakis.kozak@gmail.com>
# SPDX-License-Identifier: GPL-2.0-only

load common

setup_file() {
    if ! $(sudo dkms status | grep -q zfs); then
        echo "dkms: ZFS module wasn't prebuild. Skipping" >&3
        skip
    fi

    _lsb_release() {
	    . /etc/os-release || exit 1
	    echo "$ID"
    }
    DISTRIBUTION="$(_lsb_release)"

    case "$DISTRIBUTION" in
    arch)
            export ZFS_DKMS=$(pacman -Q zfs-dkms | sed 's/.*-dkms //g' | sed 's/-.*//g')
            export KERNEL=$(pacman -Q linux | sed 's/linux //g' | sed 's/\.arch/-arch/g')
        ;;
    ubuntu)
            export ZFS_DKMS=$(dpkg-query -W -f='${Version}\n' zfs-dkms | sed 's/-.*ubuntu.*//g' | sed 's/+.*//g')
            export KERNEL=$(uname -r)
        ;;
    debian)
            export ZFS_DKMS=$(dpkg-query -W -f='${Version}\n' zfs-dkms | sed 's/-.*deb.*//g')
            export KERNEL=$(uname -r)
        ;;
    opensuse*)
            # Tumbleweed distributes ZFS only as precompiled modules.
            # DKMS package would've needed to be built in the preparation phase
            skip
        ;;
    *) ;;
    esac

    ZFS_DKMS="zfs/$ZFS_DKMS"
}

@test "dkms: Check DKMS status" {
    sudo dkms status
}

@test "dkms: Pack the module to tarball" {
    sudo dkms mktarball $ZFS_DKMS -k $KERNEL \
    --archive /var/tmp/dkms/zfs.tar
}

@test "dkms: Remove a module" {
    sudo dkms remove $ZFS_DKMS -k $KERNEL
}

@test "dkms: Install a module" {
    sudo dkms install $ZFS_DKMS -k $KERNEL
}

@test "dkms: Load a module from tarball" {
    sudo dkms remove $ZFS_DKMS -k $KERNEL
    sudo dkms ldtarball /var/tmp/dkms/zfs.tar
}
