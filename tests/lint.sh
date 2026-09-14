#!/usr/bin/env bash
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2026 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

# Usage: just check

# Lint a profile tree with apparmor-lint

set -eu -o pipefail

readonly WORKDIR="$(pwd)"
readonly BASE="$WORKDIR/apparmor.d"
readonly EXTRA_ARG=(
	-I "$(realpath "$WORKDIR/../apparmor/profiles/apparmor.d")"
	--disable prefer-attach-disconnected-path
	--disable suggest-variable
)
readonly OUTPUT=lint-output.txt

main() {
	local rc
	set +e
	find "$BASE" -type f -print0 |
		xargs -0 -P "$(nproc)" -n 10 apparmor-lint --no-parser -b "$BASE" "${EXTRA_ARG[@]}" |
		tee "$OUTPUT"
	rc=${PIPESTATUS[1]}
	set -e

	if [ -s "$OUTPUT" ]; then
		echo ""
		echo "Diagnostic code counts:"
		grep -oE '\[[a-z][a-z-]+\]' "$OUTPUT" |
			sort | uniq -c | sort -rn |
			awk '{printf "  %6d %s\n", $1, $2}'
	fi

	[ "$rc" = 123 ] && exit 1
	exit "$rc"
}

main "$@"
