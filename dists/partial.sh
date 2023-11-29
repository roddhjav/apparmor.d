#!/usr/bin/env bash
# Partial install of apparmor profiles
# Copyright (C) 2023 monsieuremre <https://github.com/monsieuremre>
# Copyright (C) 2023 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

# Usage:
#     make
#     sudo make profile-names...

set -eu #-o pipefail

readonly BUILD=.build
readonly DESTDIR="$1"
shift

_install() {
    local profile="$1"
    if [[ ! -f "$BUILD/apparmor.d/$profile" ]]; then
        return
    fi
    if [[ -f "$DESTDIR/etc/apparmor.d/$profile" ]]; then
        return
    fi

    echo "Installing profile $profile"
    install -Dvm0644 "$BUILD/apparmor.d/$profile" "$DESTDIR/etc/apparmor.d/$profile"

    grep "rPx," "$BUILD/apparmor.d/$profile" | while read -r line; do
        [[ -z "$line" ]] && continue
        name="$(echo "$line" | awk '{print $1}')" # | awk -F"/" '{print $NF}')"
        _install "$name"
    done
    grep "rPx -> " "$BUILD/apparmor.d/$profile" | while read -r line; do
        [[ -z "$line" ]] && continue
        name=${line%%#*}
        name=$(echo "$name" | awk '{print $NF}')
        name=${name::-1}
        _install "$name"
    done
}

main() {
    for profile in "$@"; do
        _install "$profile"
    done
}

main "$@"
