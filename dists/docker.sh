#!/usr/bin/env bash
# Build the package in a clean Archlinux/openSUSE/Debian/Ubuntu container
# Copyright (C) 2022-2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

# Usage:
#  just package ubuntu 24.04
#  just package archlinux
#  just package opensuse
#  just package fedora 44

set -eu -o pipefail

readonly BASEIMAGE="${BASEIMAGE:-registry.gitlab.com/roddhjav/builders}"
readonly PREFIX="builder-"
readonly PKGNAME=apparmor.d
readonly VOLUME=/tmp/build
readonly BUILDIR=/home/build/tmp
readonly OUTPUT=".pkg"
readonly DISTRIBUTION="$1"
RELEASE="${2:-}"
FLAVOR="${3:-}"
PACKAGER="$(git config user.name) <$(git config user.email)>"
[[ "$RELEASE" == "-" ]] && RELEASE=""
readonly RELEASE FLAVOR PACKAGER

_start() {
	local img="$1"
	docker start "$img" || return 1
}

_is_running() {
	local img="$1"
	res="$(docker inspect -f '{{ .State.Running }}' "$img")" &>/dev/null
	exist=$?
	if [[ $exist -ne 0 ]]; then
		return $exist
	elif [[ "$res" == true ]]; then
		return 0
	else
		return 1
	fi
}

_exist() {
	local img="$1"
	docker inspect -f '{{ .State.Running }}' "$img" &>/dev/null
}

sync() {
	mkdir -p "$VOLUME"
	rsync -ra --delete . "$VOLUME/$PKGNAME"
}

build_in_docker_makepkg() {
	local dist="$1"
	local img="$PREFIX$dist"

	if _exist "$img"; then
		if ! _is_running "$img"; then
			_start "$img"
		fi
	else
		docker pull "$BASEIMAGE/$dist"
		docker run -tid --name "$img" --volume "$VOLUME:$BUILDIR" \
			--env PKGDEST="$BUILDIR" --env PACKAGER="$PACKAGER" \
			"$BASEIMAGE/$dist"
		docker exec "$img" sudo pacman -Sy --noconfirm --noprogressbar
	fi

	docker exec --workdir="$BUILDIR/$PKGNAME" "$img" just build-pkg
	mv "$VOLUME/$PKGNAME/$OUTPUT/$PKGNAME"*.pkg.* "$OUTPUT"
}

build_in_docker_dpkg() {
	local img dist="$1" target="$1" release="$2"

	[[ "$release" == 14 ]] && release="forky"
	img="$PREFIX$dist$release"

	# Adjustments for test flavor
	if [[ "$FLAVOR" == "test" ]]; then
		sed -i -e "s;just build=.build complain;just build=.build complain-test;" "$VOLUME/$PKGNAME/debian/rules"
	fi

	if _exist "$img"; then
		if ! _is_running "$img"; then
			_start "$img"
		fi
	else
		docker pull "$BASEIMAGE/$dist:$release"
		docker run -tid --name "$img" --volume "$VOLUME:$BUILDIR" \
			--env DISTRIBUTION="$target" "$BASEIMAGE/$dist:$release"
		docker exec "$img" sudo apt-get update -q
		docker exec "$img" sudo apt-get install -y lsb-release libdistro-info-perl
		local aptopt=()
		if [[ "$dist" == debian && "$release" == "13" ]]; then
			aptopt=(-t trixie-backports)
		fi
		docker exec "$img" sudo apt-get install -y "${aptopt[@]}" golang-go
	fi

	docker exec --workdir="$BUILDIR/$PKGNAME" "$img" just build-dpkg
	mv "$VOLUME/$PKGNAME/$OUTPUT/$PKGNAME"*.deb "$OUTPUT"
}

build_in_docker_rpm() {
	local dist="$1" release="${2:-}"
	local img="$PREFIX$dist$release"
	local tag="${release:-latest}"

	if _exist "$img"; then
		if ! _is_running "$img"; then
			_start "$img"
		fi
	else
		docker pull "$BASEIMAGE/$dist:$tag"
		docker run -tid --name "$img" --volume "$VOLUME:$BUILDIR" \
			"$BASEIMAGE/$dist:$tag"
		case "$dist" in
		opensuse)
			docker exec "$img" sudo zypper install -y distribution-release golang-packaging apparmor-profiles
			;;
		fedora)
			docker exec "$img" sudo dnf -y install golang go-rpm-macros
			;;
		esac
	fi

	docker exec --workdir="$BUILDIR/$PKGNAME" "$img" just build-rpm "$dist"
	mv "$VOLUME/$PKGNAME/$OUTPUT/$PKGNAME"*.rpm "$OUTPUT"
}

main() {
	case "$DISTRIBUTION" in
	archlinux)
		sync
		build_in_docker_makepkg "$DISTRIBUTION"
		;;

	debian | ubuntu)
		sync
		build_in_docker_dpkg "$DISTRIBUTION" "$RELEASE"
		;;

	opensuse | fedora)
		sync
		build_in_docker_rpm "$DISTRIBUTION" "$RELEASE"
		;;

	*) ;;
	esac
}

mkdir -p "$OUTPUT"
main "$@"
