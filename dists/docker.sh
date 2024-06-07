#!/usr/bin/env bash
# Build the package in a clean Archlinux/openSUSE/Debian/Ubuntu container
# Copyright (C) 2022-2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

# Usage: make package dist=<distribution>

set -eu -o pipefail

readonly BASEIMAGE="${BASEIMAGE:-registry.gitlab.com/roddhjav/builders}"
readonly PREFIX="builder-"
readonly PKGNAME=apparmor.d
readonly VOLUME=/tmp/build
readonly BUILDIR=/home/build/tmp
readonly COMMAND="$1"
VERSION="0.$(git rev-list --count HEAD)"
PACKAGER="$(git config user.name) <$(git config user.email)>"
readonly VERSION PACKAGER

_start() {
	local img="$1"
	docker start "$img"
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
			--env BUILDDIR=/tmp/build \
			"$BASEIMAGE/$dist"
	fi

	docker exec --workdir="$BUILDIR/$PKGNAME" "$img" bash dists/build.sh pkg
	mv "$VOLUME/$PKGNAME"-*.pkg.* .
}

build_in_docker_dpkg() {
	local dist="$1" target="$1"
	local img="$PREFIX$dist"

	[[ "$dist" == whonix ]] && dist=debian
	if _exist "$img"; then
		if ! _is_running "$img"; then
			_start "$img"
		fi
	else
		docker pull "$BASEIMAGE/$dist"
		docker run -tid --name "$img" --volume "$VOLUME:$BUILDIR" \
			--env DISTRIBUTION="$target" "$BASEIMAGE/$dist"
		docker exec "$img" sudo apt-get update -q
		docker exec "$img" sudo apt-get install -y config-package-dev rsync
		[[ "$dist" == debian ]] && aptopt=(-t bookworm-backports)
		docker exec "$img" sudo apt-get install -y "${aptopt[@]}" golang-go
	fi

	docker exec --workdir="$BUILDIR/$PKGNAME" "$img" bash dists/build.sh dpkg
	mv "$VOLUME/$PKGNAME/${PKGNAME}_${VERSION}-1"_*.* .
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
		docker exec "$img" sudo zypper install -y distribution-release golang-packaging rsync
	fi

	docker exec --workdir="$BUILDIR/$PKGNAME" "$img" bash dists/build.sh rpm
	mv "$VOLUME/$PKGNAME/$PKGNAME-$VERSION-"*.rpm .
}

main() {
	case "$COMMAND" in
	archlinux)
		# build_in_docker_makepkg "$COMMAND"
		PKGDEST=. makepkg -Cf
		;;

	debian | ubuntu | whonix)
		sync
		build_in_docker_dpkg "$COMMAND"
		;;

	opensuse)
		sync
		build_in_docker_rpm "$COMMAND"
		;;

	*) ;;
	esac
}

main "$@"
