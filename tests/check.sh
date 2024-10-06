#!/usr/bin/env bash
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

# Usage: make check
# shellcheck disable=SC2044

set -eu -o pipefail

readonly APPARMORD="apparmor.d"

check_profiles() {
    echo "⋅ Checking if all profiles contain:"
    echo "    - 'abi <abi/4.0>,'"
    echo "    - 'profile *profile_name* {'"
    echo "    - 'include if exists <local/*>'"
    echo "    - include if exists local for subprofiles"
    directories=("$APPARMORD/groups/*" "$APPARMORD/profiles-*-*")
    # shellcheck disable=SC2068
    for dir in ${directories[@]}; do
        for file in $(find "$dir" -maxdepth 1 -type f); do
            case "$file" in */README.md) continue ;; esac
            name="$(basename "$file")"
            name="${name/.apparmor.d/}"
            include="include if exists <local/$name>"
            if ! grep -q "^  *${include}$" "$file"; then
                echo "$name does not contain '$include'"
                exit 1
            fi
            if ! grep -q "^ *abi <abi/4.0>," "$file"; then
                echo "$name does not contain 'abi <abi/4.0>,'"
                exit 1
            fi
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
    echo "    - 'abi <abi/4.0>,'"
    echo "    - 'include if exists <abstractions/*.d>'"
    directories=(
        "$APPARMORD/abstractions/" "$APPARMORD/abstractions/app/"
        "$APPARMORD/abstractions/bus/" "$APPARMORD/abstractions/common/"
    )
    for dir in "${directories[@]}"; do
        for file in $(find "$dir" -maxdepth 1 -type f); do
            name="$(basename "$file")"
            root="${dir/${APPARMORD}\/abstractions\//}"
            include="include if exists <abstractions/${root}${name}.d>"
            if ! grep -q "^  *${include}$" "$file"; then
                echo "$file does not contain '$include'"
                exit 1
            fi
            # if ! grep -q "^ *abi <abi/4.0>," "$file"; then
            # 	echo "$file does not contain 'abi <abi/4.0>,'"
            # 	exit 1
            # fi
        done
    done

}

check_profiles
check_abstractions
