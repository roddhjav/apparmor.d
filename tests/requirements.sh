#!/usr/bin/env bash
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

# Dependencies for the bats integration tests

set -eu -o pipefail

# shellcheck source=/dev/null
_lsb_release() {
	. /etc/os-release || exit 1
	echo "$ID"
}
DISTRIBUTION="$(_lsb_release)"

case "$DISTRIBUTION" in
arch)
    # repo w/ ZFS packages, to skip AUR and build proccess
	if ! grep -q 'archzfs' /etc/pacman.conf; then
		echo '[archzfs]' | sudo tee -a /etc/pacman.conf > /dev/null
		echo 'SigLevel = Never' | sudo tee -a /etc/pacman.conf > /dev/null
		echo 'Server = https://github.com/archzfs/archzfs/releases/download/experimental' | sudo tee -a /etc/pacman.conf > /dev/null
		echo 'Server = https://mirrors.uni-plovdiv.net/archzfs/experimental' | sudo tee -a /etc/pacman.conf > /dev/null
	fi

	sudo pacman -Syu --noconfirm \
		bats bats-support \
		pacman-contrib tlp flatpak networkmanager \
		linux-headers zfs-dkms
	;;
debian | ubuntu)
    # Debian requires manual adding of contrib repo, Ubuntu has it in the standard installation
    if [ $DISTRIBUTION == "debian" ]; then
	    sudo sed -i "s/Components: main/& contrib non-free/g" /etc/apt/sources.list.d/debian.sources
    fi
    # Skip accepting licence note
    echo "zfs-dkms zfs-dkms/note-incompatible-licenses note true" | sudo debconf-set-selections

	sudo apt-get update -y
	sudo apt-get install -y \
		bats bats-support \
		cpuid dfc systemd-boot systemd-userdbd systemd-homed systemd-container tlp \
		network-manager systemd-container flatpak util-linux-extra \
		dkms zfs-dkms
	;;
opensuse*)
	;;
*) ;;
esac
