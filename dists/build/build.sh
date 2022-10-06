#!/usr/bin/env bash
# Build the package in a clean Archlinux/Debian/Ubuntu container
# Copyright (C) 2022 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

# Usage: make <distribution>

set -eu

readonly BASEIMAGE="${BASEIMAGE:-}"
readonly IMAGEPREFIX="builder-"
readonly PKGNAME=apparmor.d
readonly VOLUME=/tmp/build
readonly BUILDIR=/home/build/tmp
readonly COMMAND="$1"
VERSION="0.$(git rev-list --count HEAD)-1"
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
    local img="$IMAGEPREFIX$dist"

    if _exist "$img"; then
        if ! _is_running "$img"; then
            _start "$img"
        fi
    else
        docker build -t "$BASEIMAGE$img" "dists/build/$dist"
        docker run -tid --name "$img" --volume "$PWD:$BUILDIR" \
            --env MAKEFLAGS="-j$(nproc)" --env PACKAGER="$PACKAGER" \
            --env PKGDEST="$BUILDIR" --env DIST="$dist" \
            "$BASEIMAGE$img"
    fi

    docker exec -i "$img" \
        makepkg -sfC --noconfirm --noprogressbar
    mv "$VOLUME/$PKGNAME"-*.pkg.* .
}

build_in_docker_dpkg() {
    local dist="$1"
    local img="$IMAGEPREFIX$dist"

    if _exist "$img"; then
        if ! _is_running "$img"; then
            _start "$img"
        fi
    else
        docker build -t "$BASEIMAGE$img" "dists/build/$dist"
        docker run -tid --name "$img" --volume "$VOLUME:$BUILDIR" \
            --env DEBIAN_FRONTEND=noninteractive --env DIST="$dist" \
            "$BASEIMAGE$img"
    fi

    docker exec --workdir="$BUILDIR/$PKGNAME" "$img" \
        dch --newversion="$VERSION" --urgency=medium --distribution=stable --controlmaint "Release $VERSION"
    docker exec --workdir="$BUILDIR/$PKGNAME" "$img" \
        dpkg-buildpackage -b -d --no-sign
    mv "$VOLUME/${PKGNAME}_${VERSION}"_*.* .
}

main() {
    case "$COMMAND" in
    archlinux)
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
