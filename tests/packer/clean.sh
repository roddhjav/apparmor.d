#!/usr/bin/env bash
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2023-2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

set -eu -o pipefail

# shellcheck source=/dev/null
_lsb_release() {
	. /etc/os-release || exit 1
	echo "$ID"
}
DISTRIBUTION="$(_lsb_release)"
readonly SELF="$0"
readonly DISTRIBUTION

export HISTSIZE=0
export DEBIAN_FRONTEND=noninteractive

readonly green='\e[0;32m' Bold='\e[1m'
_msg() { printf '%b- %s%b\n' "$Bold" "$*" "$green" >&2; }
_disksize() { df --total --block-size=1 --output=size | tail -1; }
_diskused() { df --total --block-size=1 --output=used | tail -1; }

_sshdgenkeys() {
	cat <<-_EOF >/usr/lib/systemd/system/sshdgenkeys.service
		[Unit]
		Description=SSH Key Generation
		ConditionPathExists=|!/etc/ssh/ssh_host_ecdsa_key
		ConditionPathExists=|!/etc/ssh/ssh_host_ecdsa_key.pub
		ConditionPathExists=|!/etc/ssh/ssh_host_ed25519_key
		ConditionPathExists=|!/etc/ssh/ssh_host_ed25519_key.pub
		ConditionPathExists=|!/etc/ssh/ssh_host_rsa_key
		ConditionPathExists=|!/etc/ssh/ssh_host_rsa_key.pub

		[Service]
		ExecStart=/usr/bin/ssh-keygen -A
		Type=oneshot
		RemainAfterExit=yes
	_EOF
	mkdir -p /usr/lib/systemd/system/ssh.service.d
	cat <<-_EOF >/usr/lib/systemd/system/ssh.service.d/sshdgenkeys.conf
		[Unit]
		Wants=sshdgenkeys.service
		After=sshdgenkeys.service
	_EOF
}

clean_apt() {
	_msg "Cleaning the apt cache"
	apt-get -y autoremove --purge
	apt-get -y autoclean
	apt-get -y clean
}

clean_pacman() {
	_msg "Cleaning pacman cache"
	pacman -Scc --noconfirm
}

clean_zypper() {
	_msg "Cleaning zypper cache"
	zypper update -y
	zypper clean -y
}

# Make the image as impersonal as possible.
impersonalize() {
	_msg "Making the image as impersonal as possible."

	# Remove remaining pkg file, docs and caches
	dirs=(
		/usr/share/doc
		/usr/share/man
		/var/cache/
		/var/lib/apt
		/var/lib/dhcp
		/var/tmp
	)
	for dir in "${dirs[@]}"; do
		if [[ -d "$dir" ]]; then
			find "$dir" -mindepth 1 -delete
		fi
	done

	# Truncate any logs that have built up during the install
	find /var/log -type f -exec truncate --size=0 {} \;

	# Truncate the machine-id
	truncate --size=0 /etc/machine-id
	truncate --size=0 /var/lib/dbus/machine-id

	remove=(
		# Remove remaining pkg file, docs and caches
		/usr/share/info/
		/usr/share/lintian/
		/usr/share/linda/

		# Remove history & unique ids
		/etc/adjtime
		/etc/ansible/
		/etc/ssh/ssh_host_*_key*
		/home/*/.ansible/
		/home/*/.bash_history
		/home/*/.sudo_as_admin_successful
		/home/*/.zsh_history
		/root/.ansible/
		/root/.bash_history
		/root/.wget-hsts
		/var/cache/ldconfig/aux-cache
		/var/lib/systemd/random-seed
		/var/log/private

		# Remove itself
		"$(readlink -f "$SELF")"
	)
	rm -rf "${remove[@]}"
}

# Free all unused storage block.
trim() {
	local swapuuid swappart

	if swapuuid=$(blkid -o value -l -s UUID -t TYPE=swap); then
		swappart=$(readlink -f /dev/disk/by-uuid/"$swapuuid")
		_msg "Fill with 0 the swap partition to reduce box size"
		swapoff "$swappart"
		dd if=/dev/zero of="$swappart" bs=1M || true
		mkswap -U "$swapuuid" "$swappart"
	elif [[ -f /swap/swapfile ]]; then
		swapoff /swap/swapfile
		truncate --size=0 /swap/swapfile
	fi

	# Block until the empty file has been removed, otherwise, Packer will
	# try to kill the box while the disk is still full and that is bad.
	sync
}

main() {
	local begin end
	begin=$(_diskused)
	case "$DISTRIBUTION" in
	debian | ubuntu)
		clean_apt
		_sshdgenkeys
		;;

	opensuse*)
		clean_zypper
		;;

	arch)
		clean_pacman
		;;
	esac
	impersonalize
	trim

	end="$(_diskused)"
	((res = begin - end))
	echo "Inital used space: $(numfmt --to=iec "$begin")"
	echo "Final space: $(numfmt --to=iec "$end")"
	echo "Saved space: $(numfmt --to=iec "$res")"
}

main "$@"
