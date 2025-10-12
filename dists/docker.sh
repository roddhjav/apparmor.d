#!/usr/bin/env bash
# Build the package in a clean Archlinux/openSUSE/Debian/Ubuntu container
# Copyright (C) 2022-2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

# Usage:
#  just package ubuntu 24.04
#  just package archlinux
#  just package opensuse

set -eu -o pipefail

readonly BASEIMAGE="${BASEIMAGE:-registry.gitlab.com/roddhjav/builders}"
readonly PREFIX="builder-"
readonly PKGNAME=apparmor.d
readonly VOLUME=/tmp/build
readonly BUILDIR=/home/build/tmp
readonly OUTDIR=".pkg"
readonly DISTRIBUTION="$1"
RELEASE="${2:-}"
FLAVOR="${3:-}"
VERSION="0.$(git rev-list --count HEAD)"
PACKAGER="$(git config user.name) <$(git config user.email)>"
[[ "$RELEASE" == "-" ]] && RELEASE=""
readonly OUTPUT="$PWD/$OUTDIR/$DISTRIBUTION/$RELEASE"
readonly RELEASE FLAVOR VERSION PACKAGER

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
	if [[ "$FLAVOR" == "test" ]]; then
		sed -i -e "s/just complain/just complain-test/" "$VOLUME/$PKGNAME/debian/rules"
	fi
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
			--env BUILDDIR=/tmp/build \
			"$BASEIMAGE/$dist"
		docker exec "$img" sudo pacman -Sy --noconfirm --noprogressbar
	fi

	docker exec --workdir="$BUILDIR/$PKGNAME" "$img" bash dists/build.sh pkg
	mv "$VOLUME/$PKGNAME/$OUTDIR/$PKGNAME"-*.pkg.* "$OUTPUT"
}

build_in_docker_dpkg() {
	local img dist="$1" target="$1" release="$2"

	[[ "$dist" == whonix ]] && dist=debian
	img="$PREFIX$dist$release"
	if _exist "$img"; then
		if ! _is_running "$img"; then
			_start "$img"
		fi
	else
		docker pull "$BASEIMAGE/$dist:$release"
		docker run -tid --name "$img" --volume "$VOLUME:$BUILDIR" \
			--env DISTRIBUTION="$target" "$BASEIMAGE/$dist:$release"
		docker exec "$img" sudo apt-get update -q
		docker exec "$img" sudo apt-get install -y config-package-dev lsb-release libdistro-info-perl
		if [[ "$dist" == debian && "$release" == "12"  ]]; then
			aptopt=(-t bookworm-backports)
		fi
		docker exec "$img" sudo apt-get install -y "${aptopt[@]}" golang-go
	fi

	docker exec --workdir="$BUILDIR/$PKGNAME" "$img" bash dists/build.sh dpkg
	mv "$VOLUME/$PKGNAME/$OUTDIR/${PKGNAME}_${VERSION}-1"_*.* "$OUTPUT"
}

build_in_docker_rpm() {
	local dist="$1"
	local img="$PREFIX$dist"

	if _exist "$img"; then
		if ! _is_running "$img"; then
			_start "$img"
		fi
	else
		docker pull "$BASEIMAGE/$dist"
		docker run -tid --name "$img" --volume "$VOLUME:$BUILDIR" \
			"$BASEIMAGE/$dist"
		docker exec "$img" sudo zypper install -y distribution-release golang-packaging apparmor-profiles
	fi

	docker exec --workdir="$BUILDIR/$PKGNAME" "$img" bash dists/build.sh rpm
	mv "$VOLUME/$PKGNAME/$OUTDIR/$PKGNAME-$VERSION-"*.rpm "$OUTPUT"
}

main() {
	case "$DISTRIBUTION" in
	archlinux)
		build_in_docker_makepkg "$DISTRIBUTION"
		;;

	debian | ubuntu | whonix)
		sync
		build_in_docker_dpkg "$DISTRIBUTION" "$RELEASE"
		;;

	opensuse)
		sync
		build_in_docker_rpm "$DISTRIBUTION"
		;;

	*) ;;
	esac
}

mkdir -p "$OUTPUT"
main "$@"
