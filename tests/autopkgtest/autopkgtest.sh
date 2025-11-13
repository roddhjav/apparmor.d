#!/usr/bin/env bash
# Run autopkgtest in a VM
# Copyright (C) 2025 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

# Ubuntu:
#  just img ubuntu 25.10 test
#  just create ubuntu25.10 test
#  just halt ubuntu25.10 test
#  just autopkgtest ubuntu25.10
#
# Debian:
#  just img debian 13 test
#  just create debian13 test
#  just halt debian13 test
#  just autopkgtest debian13

set -eu -o pipefail

readonly COMMAND="$1"
readonly OSINFO="$2"
readonly FLAVOR="test"
readonly VERBOSE=${VERBOSE:-0}

# The maximum the host can handle
readonly CPUS=16
readonly RAM=16384
readonly TIMEOUT=1800

# As defined in Justfile
readonly PREFIX="$PREFIX"
readonly VM_DIR="$VM_DIR"
readonly USER="$USER"
readonly PASSWORD="$PASSWORD"
readonly SSH_OPT="$SSH_OPT"

readonly OUTPUT=".logs/autopkgtest/"
readonly VM_PATH="$VM_DIR/${PREFIX}${OSINFO}-${FLAVOR}.qcow2"
readonly PACKAGES_FILE="tests/autopkgtest/src-packages"
readonly reset='\e[0m' red='\e[0;31m' magenta='\e[0;35m'

_message() { printf '%b%s%b\n' "$magenta" "$*" "$reset" >&2; }
_verbose() { printf '%b>%b %s\n' "$magenta" "$reset" "$*" >&2; }
_log() { printf '%b%s%b\n' "$red" "$*" "$reset" >&2; }

_run() {
	local pkg=""
	coproc C { unbuffer -p ./tests/autopkgtest/autopkgtest.sh test "$OSINFO" 2>&1; }
	CMD_PID=$!
	while IFS= read -r line <&"${C[0]}"; do
		line="${line%$'\r'}"
		if [[ $VERBOSE -eq 0 ]]; then
			_verbose "$line"
		fi
		if [[ $line == *">>>> Testing package "* ]]; then
			pkg="${line#*">>>> Testing package "}"
			pkg="-${pkg%% <<<<*}"
		fi
		if [[ $line == "Press Enter to resume running tests." ]]; then
			# shellcheck disable=SC2086
			ssh -n $SSH_OPT -p 10022 "$USER@localhost" sudo aa-log --raw | while IFS= read -r log; do
				_log "$log"
				echo "$log" >>"$OUTPUT/aa-log-$(date +%Y%m%d-%H%M%S)$pkg"
			done
			printf '\n' >&"${C[1]}" # send Enter back over the PTY
		fi
	done
	wait $CMD_PID
}

_autopkgtest() {
	local start_from="openvpn"
	local end_at="xfsprogs"
	mapfile -t packages <"$PACKAGES_FILE"

	for pkg in "${packages[@]}"; do
		[[ "$pkg" < "$start_from" ]] && continue
		[[ "$pkg" > "$end_at" ]] && break
		_message ">>>> Testing package $pkg <<<<"
		autopkgtest "$pkg" --shell --timeout=$TIMEOUT --timeout-factor=4.0 \
			-- qemu --cpus=$CPUS --ram-size=$RAM \
			--user="$USER" --password="$PASSWORD" \
			"$VM_PATH" || true
	done
}

main() {
	case "$COMMAND" in
	run) _run ;;
	test) _autopkgtest ;;
	*) exit 1 ;;
	esac
}

mkdir -p "$OUTPUT"
main "$@"
