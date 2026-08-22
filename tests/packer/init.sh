#!/usr/bin/env bash
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2023-2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

set -eux -o pipefail

# shellcheck source=/dev/null
source /etc/os-release || exit 1
readonly SRC=/tmp/

main() {
	install -dm0750 -o "$SUDO_USER" -g "$SUDO_USER" "/home/$SUDO_USER/Projects/" "/home/$SUDO_USER/Projects/apparmor.d" "/home/$SUDO_USER/.config/"
	install -Dm0644 -o "$SUDO_USER" -g "$SUDO_USER" $SRC/.bash_aliases "/home/$SUDO_USER/.bash_aliases"
	install -Dm0644 -o "$SUDO_USER" -g "$SUDO_USER" $SRC/htoprc "/home/$SUDO_USER/.config/htop/htoprc"
	install -Dm0644 $SRC/parser.conf /etc/apparmor/parser.conf
	install -Dm0644 $SRC/site.local /etc/apparmor.d/tunables/multiarch.d/site.local
	install -Dm0755 $SRC/aa-update /usr/bin/aa-update
	install -Dm0755 $SRC/aa-clean /usr/bin/aa-clean
	chown -R "$SUDO_USER:$SUDO_USER" "/home/$SUDO_USER/.config/"

	case "$ID" in
	arch)
		rm -f $SRC/*.sig # Ignore signature files
		pacman -Rdd --noconfirm dbus-broker dbus-broker-units dbus-units || true
		pacman -U --noconfirm --overwrite '*' $SRC/*.pkg.tar.zst || true
		systemctl enable aa-flatpak || true
		;;

	debian | ubuntu)
		if [ "$ID" = ubuntu ]; then
			suffix="ubuntu1~"
		else
			suffix="1+deb"
		fi
		apt-get install -y --no-install-recommends $SRC/*-"${suffix}"*.deb
		;;

	opensuse*)
		mv "/home/$SUDO_USER/.bash_aliases" "/home/$SUDO_USER/.alias"
		rpm -i $SRC/*.rpm || true
		;;

	fedora)
		mkdir -p "/home/$SUDO_USER/.bashrc.d"
		chown -R "$SUDO_USER:$SUDO_USER" "/home/$SUDO_USER/.bashrc.d"
		mv "/home/$SUDO_USER/.bash_aliases" "/home/$SUDO_USER/.bashrc.d/aliases"
		rpm -i $SRC/*.rpm || true
		;;

	esac
}

main "$@"
