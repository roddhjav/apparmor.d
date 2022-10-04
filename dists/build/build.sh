#!/usr/bin/env bash
# Build the package in a clean Archlinux/Debian/Ubuntu container
# Copyright (C) 2022 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

# Usage: make <distribution>

set -eu

readonly BASEIMAGE="${BASEIMAGE:-}"
readonly PKGNAME=apparmor.d
readonly VOLUME=/tmp/build
readonly BUILDIR=/home/build/tmp
readonly COMMAND="$1"
VERSION="0.$(git rev-list --count HEAD)-1"
PACKAGER="$(git config user.name) <$(git config user.email)>"
readonly VERSION PACKAGER

_start() {
    local name="$1"
    docker start "$name"
}

_is_running() {
    local name="$1"
    res="$(docker inspect -f '{{ .State.Running }}' "$name")" &>/dev/null
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
    local name="$1"
    docker inspect -f '{{ .State.Running }}' "$name" &>/dev/null
}

sync() {
    mkdir -p "$VOLUME"
    rsync -ra --delete . "$VOLUME/$PKGNAME"
}

build_in_docker_makepkg() {
    local name="$1"

    if _exist "$name"; then
        if ! _is_running "$name"; then
            _start "$name"
        fi
    else
        docker build -t "$BASEIMAGE$name" "dists/build/$name"
        docker run -tid --name "$name" --volume "$VOLUME:$BUILDIR" \
            --env MAKEFLAGS="-j$(nproc)" --env PACKAGER="$PACKAGER" \
            --env PKGDEST="$BUILDIR" --env DIST="$name" \
            "$BASEIMAGE$name"
    fi

    docker exec -i --workdir="$BUILDIR/$PKGNAME" "$name" \
        makepkg -sfC --noconfirm --noprogressbar
    mv "$VOLUME/$PKGNAME"-*.pkg.* .
}

build_in_docker_dpkg() {
    local name="$1"

    if _exist "$name"; then
        if ! _is_running "$name"; then
            _start "$name"
        fi
    else
        docker build -t "$BASEIMAGE$name" "dists/build/$name"
        docker run -tid --name "$name" --volume "$VOLUME:$BUILDIR" \
            --env DEBIAN_FRONTEND=noninteractive --env DIST="$name" \
            "$BASEIMAGE$name"
    fi

    docker exec --workdir="$BUILDIR/$PKGNAME" "$name" \
        dch --newversion="$VERSION" --urgency=medium --distribution=stable --controlmaint "Release $VERSION"
    docker exec --workdir="$BUILDIR/$PKGNAME" "$name" \
        dpkg-buildpackage -b -d --no-sign
    mv "$VOLUME/${PKGNAME}_${VERSION}"_*.* .
}

main() {
    case "$COMMAND" in
    archlinux)
        sync
        build_in_docker_makepkg "$COMMAND"
        ;;

    debian | ubuntu | whonix)
        sync
        build_in_docker_dpkg "$COMMAND"
        ;;

    *) ;;
    esac
}

main "$@"
