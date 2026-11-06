#!/usr/bin/env bash
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2026 Łukasz Kozak <lakis.kozak@gmail.com>
# SPDX-License-Identifier: GPL-2.0-only

# Preperation for DKMS tests

set -eu -o pipefail

# shellcheck source=/dev/null
_lsb_release() {
	. /etc/os-release || exit 1
	echo "$ID"
}
DISTRIBUTION="$(_lsb_release)"

# rules for access to the export/import DKMS modules' directory
echo '/var/tmp/dkms/{,**} rwlk,' | sudo tee /etc/apparmor.d/local/dkms > /dev/null
sudo mkdir -p /var/tmp/dkms/

case "$DISTRIBUTION" in
arch)
    # repo w/ ZFS packages, to skip AUR and build proccess
	if ! grep -q 'archzfs' /etc/pacman.conf; then
		echo '[archzfs]' | sudo tee -a /etc/pacman.conf > /dev/null
		echo 'SigLevel = Never' | sudo tee -a /etc/pacman.conf > /dev/null
		echo 'Server = https://github.com/archzfs/archzfs/releases/download/experimental' | sudo tee -a /etc/pacman.conf > /dev/null
		echo 'Server = https://mirrors.uni-plovdiv.net/archzfs/experimental' | sudo tee -a /etc/pacman.conf > /dev/null
    fi
	;;
debian | ubuntu)
    if [ $DISTRIBUTION == "debian" ]; then
	    sudo sed -i "s/Components: main/& contrib non-free/g" /etc/apt/sources.list.d/debian.sources
    fi
    # Skip accepting licence note
    echo "zfs-dkms zfs-dkms/note-incompatible-licenses note true" | sudo debconf-set-selections
	;;
opensuse*)
	;;
*) ;;
esac
