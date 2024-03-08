#!/usr/bin/env bash
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2023-2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

set -eu

_lsb_release() {
	# shellcheck source=/dev/null
	. /etc/os-release
	echo "$ID"
}
DISTRIBUTION="$(_lsb_release)"
readonly SRC=/tmp/src
readonly DISTRIBUTION

main() {
	install -dm0750 -o "$SUDO_USER" -g "$SUDO_USER" "/home/$SUDO_USER/Projects/" "/home/$SUDO_USER/Projects/apparmor.d" "/home/$SUDO_USER/.config/"
	install -Dm0644 -o "$SUDO_USER" -g "$SUDO_USER" $SRC/.bash_aliases "/home/$SUDO_USER/.bash_aliases"
	install -Dm0644 -o "$SUDO_USER" -g "$SUDO_USER" $SRC/monitors.xml "/home/$SUDO_USER/.config/monitors.xml"
	install -Dm0644 -o "$SUDO_USER" -g "$SUDO_USER" $SRC/htoprc "/home/$SUDO_USER/.config/htop/htoprc"
	install -Dm0644 $SRC/site.local /etc/apparmor.d/tunables/multiarch.d/site.local
	install -Dm0755 $SRC/aa-update /usr/bin/aa-update
	install -Dm0755 $SRC/aa-log-clean /usr/bin/aa-log-clean
	cat $SRC/parser.conf >> /etc/apparmor/parser.conf
	chown -R "$SUDO_USER:$SUDO_USER" "/home/$SUDO_USER/.config/"
	case "$DISTRIBUTION" in
	arch) pacman --noconfirm -U $SRC/*.pkg.tar.zst ;;
	debian | ubuntu)
		apt-get update -y
		apt-get install -y apparmor-profiles build-essential config-package-dev \
			debhelper devscripts htop rsync vim
		dpkg -i $SRC/*.deb
		;;

	opensuse*)
		zypper install -y bash-completion git go htop make rsync vim
		sudo rpm -i $SRC/*.rpm
		;;

	esac
}

main "$@"
