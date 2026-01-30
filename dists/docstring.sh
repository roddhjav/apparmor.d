#!/usr/bin/env bash
# Generate markdown documentation from abstractions and tunables
# Copyright (C) 2025 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

# Usage:
#  just docstring

set -eu -o pipefail

readonly ABSTRACTIONS_DIR="apparmor.d/abstractions"
readonly TUNABLES_DIR="apparmor.d/tunables"
readonly ABSTRACTIONS_DOCS_DIR="docs/abstractions/"
readonly TUNABLES_DOCS_DIR="docs/tunables/"

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

# Categorize abstraction based on path and name
_get_abstraction_type() {
    local name="$1"
    local file="$2"

    if [[ "$name" =~ ^bus/ ]]; then
        echo "dbus"
        return

    elif [[ "$name" =~ ^sys/ ]]; then
        echo "sys"
        return

    elif [[ "$name" =~ ^udev/ ]]; then
        echo "udev"
        return

    elif [[ "$name" =~ ^attached/ ]]; then
        echo "attached"
        return

    elif [[ "$name" =~ ^app/ ]]; then
        echo "app"
        return

    elif [[ "$name" =~ ^flatpak/ ]]; then
        echo "flatpak"
        return


    elif [[ "$name" =~ ^mapping/ ]]; then
        echo "mapping"
        return

    elif [[ "$name" =~ ^common/ ]]; then
        echo "common"
        return

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

# Initialize documentation files with headers
_init_files() {
    for template in "$ABSTRACTIONS_DOCS_DIR"/.template/*.md; do
        [[ -f "$template" ]] || continue
        cp "$template" "$ABSTRACTIONS_DOCS_DIR/$(basename "$template")"
    done
    for template in "$TUNABLES_DOCS_DIR"/.template/*.md; do
        [[ -f "$template" ]] || continue
        cp "$template" "$TUNABLES_DOCS_DIR/$(basename "$template")"
    done
}

# Categorize tunable based on filename
_get_tunable_type() {
    local name="$1"

    case "$name" in
        multiarch.d/base)
            echo "base"
            ;;
        multiarch.d/extensions)
            echo "extensions"
            ;;
        multiarch.d/paths)
            echo "paths"
            ;;
        multiarch.d/profiles)
            echo "profiles"
            ;;
        multiarch.d/programs)
            echo "programs"
            ;;
        multiarch.d/system)
            echo "system"
            ;;
        multiarch.d/system-users)
            echo "system-users"
            ;;
        *)
            echo ""
            ;;
    esac
}

# Extract variables and their documentation from a tunable file
# Groups variables by paragraph (separated by empty lines)
# For each paragraph: description (comments) + all variable definitions
_get_tunable_variables() {
    local file="$1"
    local comments=""
    local preamble=""
    local variables=""
    local sections=""
    local in_header=true
    local first_var_seen=false

    # Function to format a section (paragraph of variables)
    _format_section() {
        local section_comments="$1"
        local section_variables="$2"

        [[ -z "$section_variables" ]] && return

        # Use first variable name as section title
        local first_var
        first_var=$(echo "$section_variables" | head -n1 | sed 's/^@{\([^}]*\)}.*/\1/')

        printf "\n### @{%s}\n\n" "$first_var"
        if [[ -n "$section_comments" ]]; then
            printf "%s\n\n" "$section_comments"
        fi
        # shellcheck disable=SC2016
        printf '```\n%s\n```\n' "$section_variables"
    }

    while IFS= read -r line; do
        # Empty line marks end of a paragraph
        if [[ -z "$line" ]]; then
            in_header=false
            # Save pending section if any
            if [[ -n "$variables" ]]; then
                sections="$sections$(_format_section "$comments" "$variables")"
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
            # Skip header comments (copyright, license, etc.)
            if $in_header; then
                [[ "$line" =~ (Copyright|SPDX-License-Identifier|^#\ apparmor\.d\ ) ]] && continue
                [[ "$line" =~ ^#aa: ]] && continue
            fi
            in_header=false

            # Skip vim modeline and similar
            [[ "$line" =~ ^#\ vim: ]] && continue

            # Skip commented-out variable definitions
            [[ "$line" =~ ^#\ @\{ ]] && continue

            # Skip directive lines
            [[ "$line" =~ ^#aa: ]] && continue

            # Extract comment text
            local comment="${line#\# }"
            [[ "$comment" == "#" ]] && comment=""
            if [[ -n "$comments" ]]; then
                comments="$comments"$'\n'"$comment"
            else
                comments="$comment"
            fi
            continue
        fi

        in_header=false

        # Handle variable definitions
        if [[ "$line" =~ ^@\{[^}]+\}[[:space:]]*\+?= ]]; then
            # Add to variables list
            if [[ -n "$variables" ]]; then
                variables="$variables"$'\n'"$line"
            else
                variables="$line"
            fi
        fi
    done < "$file"

    # Save last section if any
    if [[ -n "$variables" ]]; then
        sections="$sections$(_format_section "$comments" "$variables")"
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
        # Only process files in multiarch.d
        [[ "$name" != multiarch.d/* ]] && continue

        category="$(_get_tunable_type "$name")"
        [[ -z "$category" ]] && continue

        _get_tunable_variables "$file" >> "$TUNABLES_DOCS_DIR/$category.md"
    done < <(find "$TUNABLES_DIR" -type f -print0 | sort -z)
}

# Generate documentation for abstractions
_generate_abstractions_docs() {
    echo "Generating abstraction documentation..."

    while IFS= read -r -d '' file; do
        name="${file#apparmor.d/abstractions/}"
        # Skip files inside .d directories (e.g., bash.d/complete)
        [[ "$name" == *".d/"* ]] && continue
        category="$(_get_abstraction_type "$name" "$file")"
        # echo "Processing: $file -> |$category|"

        docs="$(_get_docs "$file")"
        variables="$(_get_variables "$file")"

        if [[ -z "$docs" && -z "$variables" ]]; then
            continue
        fi
        printf "\n## %s\n\n" "$name" >> "$ABSTRACTIONS_DOCS_DIR/$category.md"

        if [[ -n "$docs" ]]; then
            printf "%s\n" "$docs" >> "$ABSTRACTIONS_DOCS_DIR/$category.md"
        fi
        if [[ -n "$variables" ]]; then
            printf "\n**Required variables:** %s\n" "$variables" >> "$ABSTRACTIONS_DOCS_DIR/$category.md"
        fi
    done < <(find "$ABSTRACTIONS_DIR" -type f -print0 | sort -z)
}

mkdir -p "$ABSTRACTIONS_DOCS_DIR" "$TUNABLES_DOCS_DIR"
_init_files
_generate_abstractions_docs
_generate_tunables_docs
