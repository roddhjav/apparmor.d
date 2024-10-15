#!/usr/bin/env bash
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

# Usage: make check
# shellcheck disable=SC2044

set -eu -o pipefail

readonly APPARMORD="apparmor.d"

_ensure_header() {
    local file="$1"
    headers=(
        "# apparmor.d - Full set of apparmor profiles"
        "# Copyright (C) "
        "# SPDX-License-Identifier: GPL-2.0-only"
    )
    for header in "${headers[@]}"; do
        if ! grep -q "^$header" "$file"; then
            echo "$file does not contain '$header'"
            exit 1
        fi
    done
}

_ensure_include() {
    local file="$1"
    local include="$2"
    if ! grep -q "^  *${include}$" "$file"; then
        echo "$file does not contain '$include'"
        exit 1
    fi
}

_ensure_abi() {
    local file="$1"
    if ! grep -q "^ *abi <abi/4.0>," "$file"; then
        echo "$file does not contain 'abi <abi/4.0>,'"
        exit 1
    fi
}

_ensure_vim() {
    local file="$1"
    if ! grep -q "^# vim:syntax=apparmor" "$file"; then
        echo "$file does not contain '# vim:syntax=apparmor'"
        exit 1
    fi
}

check_profiles() {
    echo "⋅ Checking if all profiles contain:"
    echo "    - apparmor.d header & license"
    echo "    - 'abi <abi/4.0>,'"
    echo "    - 'profile <profile_name>'"
    echo "    - 'include if exists <local/*>'"
    echo "    - include if exists local for subprofiles"
    echo "    - vim:syntax=apparmor"
    directories=("$APPARMORD/groups/*" "$APPARMORD/profiles-*-*")
    # shellcheck disable=SC2068
    for dir in ${directories[@]}; do
        for file in $(find "$dir" -maxdepth 1 -type f); do
            case "$file" in */README.md) continue ;; esac
            name="$(basename "$file")"
            name="${name/.apparmor.d/}"
            include="include if exists <local/$name>"
            _ensure_header "$file"
            _ensure_include "$file" "$include"
            _ensure_abi "$file"
            _ensure_vim "$file"
            if ! grep -q "^profile $name" "$file"; then
                echo "$name does not contain 'profile $name'"
                exit 1
            fi
            mapfile -t subrofiles < <(grep "^  *profile*" "$file" | awk '{print $2}')
            for subprofile in "${subrofiles[@]}"; do
                include="include if exists <local/${name}_${subprofile}>"
                if ! grep -q "^  *${include}$" "$file"; then
                    echo "$name: $name//$subprofile does not contain '$include'"
                    exit 1
                fi
            done
        done
    done
}

check_abstractions() {
    echo "⋅ Checking if all abstractions contain:"
    echo "    - apparmor.d header & license"
    echo "    - 'abi <abi/4.0>,'"
    echo "    - 'include if exists <abstractions/*.d>'"
    echo "    - vim:syntax=apparmor"
    directories=(
        "$APPARMORD/abstractions/" "$APPARMORD/abstractions/app/"
        "$APPARMORD/abstractions/attached/"
        "$APPARMORD/abstractions/bus/" "$APPARMORD/abstractions/common/"
    )
    for dir in "${directories[@]}"; do
        for file in $(find "$dir" -maxdepth 1 -type f); do
            name="$(basename "$file")"
            root="${dir/${APPARMORD}\/abstractions\//}"
            include="include if exists <abstractions/${root}${name}.d>"
            _ensure_header "$file"
            _ensure_include "$file" "$include"
            _ensure_abi "$file"
            _ensure_vim "$file"
        done
    done

}

check_profiles
check_abstractions
