#!/usr/bin/env bash
# Generate markdown documentation from abstractions, tunables, and profiles
# Copyright (C) 2025 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

# Usage:
#  just docstring

set -eu -o pipefail

readonly REPO_URL="${REPO_URL:-https://github.com/roddhjav/apparmor.d}"
readonly REPO_BRANCH="${REPO_BRANCH:-main}"
readonly UPSTREAM_URL="${UPSTREAM_URL:-https://gitlab.com/apparmor/apparmor}"
readonly UPSTREAM_BRANCH="${UPSTREAM_BRANCH:-master}"

# Line-anchored links must pin to the exact commit whose content was parsed;
# a moving branch ref would drift the #L line number. Falls back to the branch.
readonly REPO_COMMIT="${REPO_COMMIT:-$(git rev-parse HEAD 2>/dev/null || echo "$REPO_BRANCH")}"

readonly ABSTRACTIONS_DIR="apparmor.d/abstractions"
readonly ABSTRACTIONS_DOCS_DIR="docs/abstractions"

readonly TUNABLES_DEFAULT_DIR="${TUNABLES_DEFAULT_DIR:-../apparmor/profiles/apparmor.d/tunables}"
readonly TUNABLES_DIR="apparmor.d/tunables"
readonly TUNABLES_DOCS_DIR="docs/tunables"

readonly PROFILES_DIR="apparmor.d"
readonly PROFILES_DOCS_DIR="docs/profiles"

# Extract description from abstraction file header
# Skips copyright/license and extracts meaningful comments
_get_docs() {
    local file="$1"
    while IFS= read -r line; do
        [[ -z "$line" ]] && continue
        [[ ! "$line" =~ ^# ]] && break
        [[ "$line" =~ (Copyright|SPDX-License-Identifier|^#\ apparmor\.d\ ) ]] && continue
        [[ "$line" =~ (LOGPROF-SUGGEST|NEEDS-VARIABLE) ]] && continue
        [[ "$line" =~ ^#aa: ]] && continue

        local comment="${line#\# }"
        if [[ $comment == "#" ]]; then
            echo ""
        elif [[ -n "$comment" ]]; then
            echo "$comment"
        fi
    done < "$file"
}

# Extract required variables from abstraction file
_get_variables() {
    local file="$1"
    grep "^# NEEDS-VARIABLE:" "$file" 2>/dev/null | sed 's/^# NEEDS-VARIABLE: //' | tr '\n' ', ' | sed 's/, $//' || true
}

# Check if an abstraction is core (has no external includes)
# Returns 0 if core, 1 if not core
_is_core_abstraction() {
    local file="$1"
    local name="$2"

    # Get base name without path components
    local base_name
    base_name="$(basename "$name")"

    # Check for include directives, excluding self-includes to .d directories
    while IFS= read -r line; do
        # Skip comments and empty lines
        [[ "$line" =~ ^[[:space:]]*# ]] && continue
        [[ -z "$line" ]] && continue

        # Check for include directives
        if [[ "$line" =~ ^[[:space:]]*include ]]; then
            # Extract the included path
            local included_path="${line#*<abstractions/}"
            included_path="${included_path%%>*}"

            # Skip if it's an include to own .d directory (e.g., amdgpu.d, amdgpu-strict.d)
            [[ "$included_path" =~ ^${base_name}(\.d|-strict\.d)$ ]] && continue

            # Found an external include, not core
            return 1
        fi
    done < "$file"

    # No external includes found, it's core
    return 0
}

# Abstraction category by top-level directory. Anything without a matching
# subdirectory falls back to core (no external includes) or generic.
declare -A ABSTRACTION_TYPES=(
    [app]=app
    [attached]=attached
    [bus]=dbus
    [common]=common
    [flatpak]=flatpak
    [mapping]=mapping
    [sys]=sys
    [udev]=udev
)

# Categorize abstraction based on path and name
_get_abstraction_type() {
    local name="$1"
    local file="$2"

    if [[ "$name" == */* ]]; then
        local type="${ABSTRACTION_TYPES[${name%%/*}]:-}"
        [[ -n "$type" ]] && { echo "$type"; return; }
    fi

    # Check if it's a core abstraction (no external includes)
    if _is_core_abstraction "$file" "$name"; then
        echo "core"
        return
    fi

    # Generic abstractions (layer 1) - generic access patterns
    # Everything else that's not in a subdirectory
    echo "generic"
}

# Project tunable doc category by path. Only multiarch.d files with a matching
# docs template are documented; anything else returns "" and is skipped.
declare -A TUNABLE_TYPES=(
    [multiarch.d/base]=base
    [multiarch.d/extensions]=extensions
    [multiarch.d/paths]=paths
    [multiarch.d/profiles]=profiles
    [multiarch.d/programs]=programs
    [multiarch.d/system]=system
    [multiarch.d/system-users]=system-users
    [home.d/apparmor.d]=home
    [xdg-user-dirs.d/apparmor.d]=home
)

# Categorize tunable based on filename
_get_tunable_type() {
    echo "${TUNABLE_TYPES[$1]:-}"
}

# Extract variables and their documentation from a tunable file
# Groups variables by paragraph (separated by empty lines)
# For each paragraph: description (comments) + all variable definitions
_get_tunable_variables() {
    local file="$1"
    local source_url="${2:-}"
    local comments=""
    local preamble=""
    local variables=""
    local sections=""
    local first_var_seen=false
    local lineno=0
    local section_line=0

    # Function to format a section (paragraph of variables)
    _format_section() {
        local section_comments="$1"
        local section_variables="$2"
        local section_url="${3:-}"

        [[ -z "$section_variables" ]] && return

        # Use first variable name as section title
        local first_var
        first_var=$(echo "$section_variables" | head -n1 | sed 's/^@{\([^}]*\)}.*/\1/')

        printf "\n### @{%s}\n\n" "$first_var"
        [[ -n "$section_url" ]] && \
            printf '[:material-file-eye:](%s "View source"){ .abs-source }\n\n' "$section_url"
        if [[ -n "$section_comments" ]]; then
            printf "%s\n\n" "$section_comments"
        fi
        # shellcheck disable=SC2016
        printf '```\n%s\n```\n' "$section_variables"
    }

    while IFS= read -r line; do
        lineno=$((lineno + 1))
        # Empty line marks end of a paragraph
        if [[ -z "$line" ]]; then
            # Save pending section if any
            if [[ -n "$variables" ]]; then
                local section_url=""
                [[ -n "$source_url" ]] && section_url="$source_url#L$section_line"
                sections="$sections$(_format_section "$comments" "$variables" "$section_url")"
                variables=""
                comments=""
                first_var_seen=true
            elif ! $first_var_seen && [[ -n "$comments" ]]; then
                # Save pre-first-variable comments as preamble
                if [[ -n "$preamble" ]]; then
                    preamble="$preamble"$'\n\n'"$comments"
                else
                    preamble="$comments"
                fi
                comments=""
            fi
            continue
        fi

        # Handle comments
        if [[ "$line" =~ ^# ]]; then
            # Skip copyright/license boilerplate.
            [[ "$line" =~ (Copyright|SPDX-License-Identifier|free software|redistribute|modify it under|License published|Free Software Foundation|GNU General Public|WITHOUT ANY WARRANTY|Author:|Last Modified:) ]] && continue
            [[ "$line" =~ ^#[[:space:]]*-+[[:space:]]*$ ]] && continue
            [[ "$line" =~ ^#\ apparmor\.d\  ]] && continue

            # Skip vim modeline and similar
            [[ "$line" =~ ^#\ vim: ]] && continue

            # Skip commented-out variable definitions (e.g. "# @{FOO}=bar" or
            # "#@{FOO}+=bar"), but keep prose that merely starts with a variable name.
            [[ "$line" =~ ^#[[:space:]]*@\{[^}]+\}[[:space:]]*\+?= ]] && continue

            # Skip directive lines
            [[ "$line" =~ ^#aa: ]] && continue

            # Extract comment text. Strip the leading marker (and one space) so a
            # comment like "#same pattern ..." never renders as a markdown heading.
            local comment="${line#\#}"
            comment="${comment# }"
            if [[ -n "$comments" ]]; then
                comments="$comments"$'\n'"$comment"
            else
                comments="$comment"
            fi
            continue
        fi

        # Handle variable definitions
        if [[ "$line" =~ ^@\{[^}]+\}[[:space:]]*\+?= ]]; then
            # Add to variables list
            if [[ -n "$variables" ]]; then
                variables="$variables"$'\n'"$line"
            else
                variables="$line"
                section_line=$lineno
            fi
        fi
    done < "$file"

    # Save last section if any
    if [[ -n "$variables" ]]; then
        local section_url=""
        [[ -n "$source_url" ]] && section_url="$source_url#L$section_line"
        sections="$sections$(_format_section "$comments" "$variables" "$section_url")"
    fi

    # Output preamble first, then sections
    if [[ -n "$preamble" ]]; then
        printf "\n%s\n" "$preamble"
    fi
    printf "%s" "$sections"
}

# Generate documentation for tunables
_generate_tunables_docs() {
    echo "Generating tunable documentation..."

    while IFS= read -r -d '' file; do
        name="${file#apparmor.d/tunables/}"
        category="$(_get_tunable_type "$name")"
        [[ -z "$category" ]] && continue

        _get_tunable_variables "$file" "$REPO_URL/blob/$REPO_COMMIT/apparmor.d/tunables/$name" \
            >> "$TUNABLES_DOCS_DIR/$category.md"
    done < <(find "$TUNABLES_DIR" -type f -print0 | sort -z)
}

# Group an upstream tunable file into a preset page. Anything unmapped falls
# back to 'system' so a newly added upstream file is never dropped.
declare -A DEFAULT_TUNABLE_TYPES=(
    [system]=base
    [home]=user
    [rygel]=user
    [share]=user
    [xdg-user-dirs]=user
    [dovecot]=programs
    [gs]=programs
    [ntpd]=programs
    [print-devices]=programs
)

_default_tunable_category() {
    echo "${DEFAULT_TUNABLE_TYPES[$1]:-system}"
}

# Generate documentation for the upstream/default tunable variables.
# Runs only when the upstream apparmor tunables dir is present (not automatic).
# Upstream files are bucketed into three preset pages under tunables/default/
# (system/user/programs) to avoid clashing with the project's own category pages.
# The nav entries are added/removed accordingly so they never dangle.
_generate_default_tunables_docs() {
    local outdir="$TUNABLES_DOCS_DIR/default"
    local -a entries=()

    # Always drop the stray default.md copied from the template; it is never a
    # real page (the real pages live under default/).
    rm -f "$TUNABLES_DOCS_DIR/default.md" || true

    # Without the upstream checkout we cannot regenerate the default/ pages. They
    # are committed, so leave them (and the nav) untouched rather than wiping them
    # and dangling the index.md "Default" link.
    if [[ ! -d "$TUNABLES_DEFAULT_DIR" ]]; then
        return
    fi

    # Drop previous output so pages for removed upstream tunables don't linger.
    rm -rf "$outdir" || true

    echo "Generating default (upstream) tunable documentation..."
    mkdir -p "$outdir"
    # Pin to the parsed upstream commit so #L anchors match (see REPO_COMMIT).
    local upstream_commit
    upstream_commit="$(git -C "$TUNABLES_DEFAULT_DIR" rev-parse HEAD 2>/dev/null || echo "$UPSTREAM_BRANCH")"
    local -A seen=()
    while IFS= read -r -d '' file; do
        local name="${file#"$TUNABLES_DEFAULT_DIR"/}"
        # Skip site-local stubs in .d directories
        [[ "$name" == *.d/* ]] && continue
        # Only files that actually define variables
        grep -q '^@{' "$file" || continue
        local vars
        vars="$(_get_tunable_variables "$file" "$UPSTREAM_URL/-/blob/$upstream_commit/profiles/apparmor.d/tunables/$name")"
        [[ -z "$vars" ]] && continue

        local category
        category="$(_default_tunable_category "$name")"
        local out="$outdir/$category.md"
        if [[ -z "${seen[$category]:-}" ]]; then
            printf -- '---\ntitle: %s\ntags:\n  - tunables\n  - default\n---\n' "${category^}" > "$out"
            seen[$category]=1
        fi
        printf '\n## %s\n%s\n' "$name" "$vars" >> "$out"
    done < <(find "$TUNABLES_DEFAULT_DIR" -type f -print0 | sort -z)

    # Fixed nav order; only categories that got content.
    local category
    for category in base system user programs; do
        [[ -f "$outdir/$category.md" ]] && entries+=("tunables/default/$category.md")
    done

    _update_zensical_tunables_nav "${entries[@]}"
}

# Update zensical.toml default-tunables nav (same marker scheme as profiles).
# Emits a dedicated "Default" section, but only when it has entries, so the whole
# group disappears (rather than rendering empty) when upstream is absent.
_update_zensical_tunables_nav() {
    local nav_entries="" entry
    for entry in "$@"; do
        nav_entries+="            \"$entry\","$'\n'
    done

    local block=""
    [[ -n "$nav_entries" ]] && block="        { \"Default\" = [
${nav_entries}        ] },
"

    awk -v block="$block" '
        /# BEGIN DEFAULT TUNABLES/ {
            print
            printf "%s", block
            skip = 1
            next
        }
        /# END DEFAULT TUNABLES/ {
            skip = 0
        }
        !skip { print }
    ' zensical.toml > zensical.toml.tmp && mv zensical.toml.tmp zensical.toml
}

# Generate documentation for abstractions
_generate_abstractions_docs() {
    echo "Generating abstraction documentation..."

    while IFS= read -r -d '' file; do
        name="${file#apparmor.d/abstractions/}"
        # Skip files inside .d directories (e.g., bash.d/complete)
        [[ "$name" == *".d/"* ]] && continue
        category="$(_get_abstraction_type "$name" "$file")"
        docs="$(_get_docs "$file")"
        variables="$(_get_variables "$file")"

        if [[ -z "$docs" && -z "$variables" ]]; then
            continue
        fi
        printf "\n## %s\n\n" "$name" >> "$ABSTRACTIONS_DOCS_DIR/$category.md"
        printf '[:material-file-eye:](%s/blob/%s/apparmor.d/abstractions/%s "View source"){ .abs-source }\n\n' \
            "$REPO_URL" "$REPO_BRANCH" "$name" >> "$ABSTRACTIONS_DOCS_DIR/$category.md"

        if [[ -n "$docs" ]]; then
            printf "%s\n" "$docs" >> "$ABSTRACTIONS_DOCS_DIR/$category.md"
        fi
        if [[ -n "$variables" ]]; then
            printf "\n**Required variables:** %s\n" "$variables" >> "$ABSTRACTIONS_DOCS_DIR/$category.md"
        fi
    done < <(find "$ABSTRACTIONS_DIR" -type f -print0 | sort -z)
}

# Append one profile entry to the per-category docs file.
# Reads category_has_content from the caller's scope (bash dynamic scope)
# to write a category header exactly once.
_emit_profile_doc() {
    local file="$1" category="$2" profile_name="$3"
    local docs
    docs="$(_get_docs "$file")"
    [[ -z "$docs" ]] && return

    local outfile="$PROFILES_DOCS_DIR/$category.md"
    if [[ -z "${category_has_content[$category]:-}" ]]; then
        printf "# %s\n" "$category" > "$outfile"
        category_has_content[$category]=1
    fi
    {
        printf "\n## %s\n\n" "$profile_name"
        printf '[:material-file-eye:](%s/blob/%s/%s "View source"){ .abs-source }\n\n' \
            "$REPO_URL" "$REPO_BRANCH" "$file"
        printf "%s\n" "$docs"
    } >> "$outfile"
}

# Generate documentation for profiles
# Processes profiles in groups/ and profiles-*-*/ directories
# Only includes profiles with actual docstrings
_generate_profiles_docs() {
    echo "Generating profile documentation..."

    local -A category_has_content
    local -a generated_groups=()

    while IFS= read -r -d '' file; do
        local relpath="${file#apparmor.d/groups/}"
        local group="${relpath%%/*}"
        local was_seen="${category_has_content[$group]:-}"
        _emit_profile_doc "$file" "$group" "${relpath##*/}"
        # Capture each group exactly once, only after a real (non-skipped) emit.
        if [[ -z "$was_seen" && -n "${category_has_content[$group]:-}" ]]; then
            generated_groups+=("$group")
        fi
    done < <(find "$PROFILES_DIR/groups" -type f -print0 | sort -z)

    for profiles_subdir in "$PROFILES_DIR"/profiles-*; do
        [[ -d "$profiles_subdir" ]] || continue
        local category
        category="$(basename "$profiles_subdir")"
        while IFS= read -r -d '' file; do
            _emit_profile_doc "$file" "$category" "$(basename "$file")"
        done < <(find "$profiles_subdir" -type f -print0 | sort -z)
    done

    _update_zensical_nav "${generated_groups[@]}"
}

# Update zensical.toml with generated profile groups nav entries.
_update_zensical_nav() {
    local zensical_file="zensical.toml"
    local nav_entries=""
    local group
    for group in "$@"; do
        nav_entries+="            \"profiles/$group.md\","$'\n'
    done

    awk -v entries="$nav_entries" '
        /# BEGIN PROFILE GROUPS/ {
            print
            printf "%s", entries
            skip = 1
            next
        }
        /# END PROFILE GROUPS/ {
            skip = 0
        }
        !skip { print }
    ' "$zensical_file" > "${zensical_file}.tmp" && mv "${zensical_file}.tmp" "$zensical_file"
}

main() {
	# Initialize documentation files with headers
	mkdir -p "$ABSTRACTIONS_DOCS_DIR" "$TUNABLES_DOCS_DIR" "$PROFILES_DOCS_DIR"
    for template in "$ABSTRACTIONS_DOCS_DIR"/.template/*.md; do
        [[ -f "$template" ]] || continue
        cp "$template" "$ABSTRACTIONS_DOCS_DIR/$(basename "$template")"
    done
    for template in "$TUNABLES_DOCS_DIR"/.template/*.md; do
        [[ -f "$template" ]] || continue
        cp "$template" "$TUNABLES_DOCS_DIR/$(basename "$template")"
    done

	# Generate the docs
	_generate_abstractions_docs
	_generate_tunables_docs
	_generate_default_tunables_docs
	_generate_profiles_docs
}

main "$@"
