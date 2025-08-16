#!/usr/bin/env bash
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2023-2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

set -eux

_lsb_release() {
	# shellcheck source=/dev/null
	. /etc/os-release
	echo "$ID"
}
DISTRIBUTION="$(_lsb_release)"
readonly SRC=/tmp/
readonly DISTRIBUTION

main() {
	install -dm0750 -o "$SUDO_USER" -g "$SUDO_USER" "/home/$SUDO_USER/Projects/" "/home/$SUDO_USER/Projects/apparmor.d" "/home/$SUDO_USER/.config/"
	install -Dm0644 -o "$SUDO_USER" -g "$SUDO_USER" $SRC/.bash_aliases "/home/$SUDO_USER/.bash_aliases"
	install -Dm0644 -o "$SUDO_USER" -g "$SUDO_USER" $SRC/htoprc "/home/$SUDO_USER/.config/htop/htoprc"
	install -Dm0644 $SRC/parser.conf /etc/apparmor/parser.conf
	install -Dm0644 $SRC/site.local /etc/apparmor.d/tunables/multiarch.d/site.local
	install -Dm0755 $SRC/aa-update /usr/bin/aa-update
	install -Dm0755 $SRC/aa-clean /usr/bin/aa-clean
	chown -R "$SUDO_USER:$SUDO_USER" "/home/$SUDO_USER/.config/"

	case "$DISTRIBUTION" in
	arch)
		rm -f $SRC/*.sig # Ignore signature files
		pacman --noconfirm -U $SRC/*.pkg.tar.zst
		;;

	debian | ubuntu)
		apt install -y apparmor-profiles
		dpkg -i $SRC/*.deb || true
		;;

	opensuse*)
		mv "/home/$SUDO_USER/.bash_aliases" "/home/$SUDO_USER/.alias"
		rpm -i $SRC/*.rpm
		;;

	esac

	verb="start"
	rm -rf /var/cache/apparmor/* || true
	if systemctl is-active -q apparmor; then
		verb="reload"
	fi
	systemctl "$verb" apparmor.service || journalctl -xeu apparmor.service
}

main "$@"
