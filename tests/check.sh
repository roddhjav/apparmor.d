#!/usr/bin/env bash
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

# Usage: make check
# shellcheck disable=SC2044

set -eu -o pipefail

readonly APPARMORD="apparmor.d"
readonly HEADERS=(
    "# apparmor.d - Full set of apparmor profiles"
    "# Copyright (C) "
    "# SPDX-License-Identifier: GPL-2.0-only"
)

_die() {
    echo -e "\033[1;31m ✗ Error: \033[0m$*"
    exit 1
}

_ensure_header() {
    local file="$1"
    for header in "${HEADERS[@]}"; do
        if ! grep -q "^$header" "$file"; then
            _die "$file does not contain '$header'"
        fi
    done
}

_ensure_indentation() {
    local file="$1"
    local in_profile=false
    local first_line_after_profile=true
    local line_number=0

    while IFS= read -r line; do
        line_number=$((line_number + 1))

        if [[ "$line" =~ $'\t' ]]; then
            _die "$file:$line_number: tabs are not allowed."
        fi

        if [[ "$line" =~ ^profile ]]; then
            in_profile=true
            first_line_after_profile=true

        elif [[ "$line" =~ [[:space:]]+$ ]]; then
            _die "$file:$line_number: line has trailing whitespace."

        elif $in_profile; then
            if $first_line_after_profile; then
                local leading_spaces="${line%%[! ]*}"
                local num_spaces=${#leading_spaces}
                if ((num_spaces != 2)); then
                    _die "$file: profile must have a two-space indentation."
                fi
                first_line_after_profile=false

            else
                local leading_spaces="${line%%[! ]*}"
                local num_spaces=${#leading_spaces}

                if ((num_spaces % 2 != 0)); then
                    ok=false
                    for offset in 5 11; do
                        num_spaces=$((num_spaces - offset))
                        if ((num_spaces < 0)); then
                            break
                        fi
                        if ((num_spaces % 2 == 0)); then
                            ok=true
                            break
                        fi
                    done

                    if ! $ok; then
                        _die "$file:$line_number: invalid indentation."
                    fi
                fi
            fi
        fi
    done <"$file"
}

_ensure_include() {
    local file="$1"
    local include="$2"
    if ! grep -q "^  *${include}$" "$file"; then
        _die "$file does not contain '$include'"
    fi
}

_ensure_abi() {
    local file="$1"
    if ! grep -q "^ *abi <abi/4.0>," "$file"; then
        _die "$file does not contain 'abi <abi/4.0>,'"
    fi
}

_ensure_vim() {
    local file="$1"
    if ! grep -q "^# vim:syntax=apparmor" "$file"; then
        _die "$file does not contain '# vim:syntax=apparmor'"
    fi
}

check_sbin() {
    echo -e "\033[1m ⋅ \033[0mEnsuring '@{sbin}' is used in all profiles:"
    while IFS= read -r name; do
        mapfile -t files < <(grep -l -R "@{bin}/$name" apparmor.d)
        for file in "${files[@]}"; do
            _die "$file contains '@{bin}/$name' instead of '@{sbin}/$name'"
        done
    done <tests/sbin.list
}

check_profiles() {
    echo -e "\033[1m ⋅ \033[0mChecking if all profiles contain:"
    echo "    - apparmor.d header & license"
    echo "    - Check indentation: 2 spaces"
    echo "    - Check for trailing whitespaces"
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
            _ensure_indentation "$file"
            _ensure_include "$file" "$include"
            _ensure_abi "$file"
            _ensure_vim "$file"
            if ! grep -q "^profile $name" "$file"; then
                _die "$name does not contain 'profile $name'"
            fi
            mapfile -t subrofiles < <(grep "^  *profile*" "$file" | awk '{print $2}')
            for subprofile in "${subrofiles[@]}"; do
                include="include if exists <local/${name}_${subprofile}>"
                if ! grep -q "^  *${include}$" "$file"; then
                    _die "$name: $name//$subprofile does not contain '$include'"
                fi
            done
        done
    done
}

check_abstractions() {
    echo -e "\033[1m ⋅ \033[0mChecking if all abstractions contain:"
    echo "    - apparmor.d header & license"
    echo "    - Check indentation: 2 spaces"
    echo "    - Check for trailing whitespaces"
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
            _ensure_indentation "$file"
            _ensure_include "$file" "$include"
            _ensure_abi "$file"
            _ensure_vim "$file"
        done
    done
}

check_sbin
check_profiles
check_abstractions
