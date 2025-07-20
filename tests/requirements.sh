#!/usr/bin/env bash
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

# Dependencies for the bats integration tests

set -eu

# shellcheck source=/dev/null
_lsb_release() {
	. /etc/os-release || exit 1
	echo "$ID"
}
DISTRIBUTION="$(_lsb_release)"

case "$DISTRIBUTION" in
arch)
	sudo pacman -Syu --noconfirm \
		bats bats-support \
		pacman-contrib tlp flatpak networkmanager
	;;
debian | ubuntu | whonix)
	sudo apt update -y
	sudo apt install -y \
		bats bats-support \
		cpuid dfc systemd-boot systemd-userdbd systemd-homed systemd-container tlp \
		network-manager systemd-container flatpak util-linux-extra
	;;
opensuse*)
	;;
*) ;;
esac
