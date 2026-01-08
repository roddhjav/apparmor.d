#!/usr/bin/env bash
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

export BATS_LIB_PATH=${BATS_LIB_PATH:-/usr/lib/bats}
load "$BATS_LIB_PATH/bats-support/load"

export SYSTEMD_PAGER=

# Ignore the profile not managed by apparmor.d
IGNORE=(/usr/sbin/mysqld php-fpm snapd/snap-confine snap.vault.vaultd /dev/pts/0)

# User password for sudo commands
export PASSWORD=${PASSWORD:-user}

export XDG_CACHE_DIR=".cache"
export XDG_CONFIG_DIR=".config"
export XDG_DATA_DIR=".local/share"
export XDG_STATE_DIR=".local/state"
export XDG_BIN_DIR=".local/bin"
export XDG_LIB_DIR=".local/lib"

# Define extended user directories not defined in the XDG standard but commonly
# used in profiles
export XDG_SCREENSHOTS_DIR="Pictures/Screenshots"
export XDG_WALLPAPERS_DIR="Pictures/Wallpapers"
export XDG_BOOKS_DIR="Books"
export XDG_GAMES_DIR="Games"
export XDG_PROJECTS_DIR="Projects"
export XDG_WORK_DIR="Work"
export XDG_MAIL_DIR="Mail"
export XDG_SYNC_DIR="Sync"
export XDG_TORRENTS_DIR="Torrents"
export XDG_GAMESSTUDIO_DIR="unity3d"

# Define user directories for virtual machines, shared folders and disk images
export XDG_VM_DIR=".vm"
export XDG_VMSHARE_DIR=".vmshare"
export XDG_IMG_DIR=".img"

# Define user build directories and artifacts output
export XDG_BUILD_DIR=".build"
export XDG_PKG_DIR=".pkg"

# Define user personal keyrings
export XDG_GPG_DIR=".gnupg"
export XDG_SSH_DIR=".ssh"
export XDG_PASSWORDSTORE_DIR=".password-store"

# Define user personal private directories
export XDG_PRIVATE_DIR=".private"

# Full path of the XDG Base Directory
export user_cache_dirs=$HOME/$XDG_CACHE_DIR
export user_config_dirs=$HOME/$XDG_CONFIG_DIR
export user_state_dirs=$HOME/$XDG_STATE_DIR
export user_bin_dirs=$HOME/$XDG_BIN_DIR
export user_lib_dirs=$HOME/$XDG_LIB_DIR

# Other user directories
export user_desktop_dirs=$HOME/$XDG_DESKTOP_DIR
export user_download_dirs=$HOME/$XDG_DOWNLOAD_DIR
export user_templates_dirs=$HOME/$XDG_TEMPLATES_DIR
export user_publicshare_dirs=$HOME/$XDG_PUBLICSHARE_DIR
export user_documents_dirs=$HOME/$XDG_DOCUMENTS_DIR
export user_music_dirs=$HOME/$XDG_MUSIC_DIR
export user_pictures_dirs=$HOME/$XDG_PICTURES_DIR
export user_videos_dirs=$HOME/$XDG_VIDEOS_DIR
export user_books_dirs=$HOME/$XDG_BOOKS_DIR
export user_games_dirs=$HOME/$XDG_GAMES_DIR
export user_projects_dirs=$HOME/$XDG_PROJECTS_DIR
export user_work_dirs=$HOME/$XDG_WORK_DIR
export user_mail_dirs=$HOME/$XDG_MAIL_DIR
export user_sync_dirs=$HOME/$XDG_SYNC_DIR
export user_torrents_dirs=$HOME/$XDG_TORRENTS_DIR
export user_vm_dirs=$HOME/$XDG_VM_DIR
export user_vmshare_dirs=$HOME/$XDG_VMSHARE_DIR
export user_img_dirs=$HOME/$XDG_IMG_DIR
export user_build_dirs=$HOME/$XDG_BUILD_DIR
export user_pkg_dirs=$HOME/$XDG_PKG_DIR
export user_gpg_dirs=$HOME/$XDG_GPG_DIR
export user_ssh_dirs=$HOME/$XDG_SSH_DIR
export user_passwordstore_dirs=$HOME/$XDG_PASSWORDSTORE_DIR
export user_private_dirs=$HOME/$XDG_PRIVATE_DIR

_START="$(date +%s)"
PROGRAM="$(basename "$BATS_TEST_FILENAME")"
PROGRAM="${PROGRAM%.*}"
export _START PROGRAM

skip_if_not_installed() {
    if ! which "$PROGRAM" &>/dev/null; then
        skip "$PROGRAM is not installed"
    fi
}

aa_setup() {
    aa_start
    skip_if_not_installed
}

aa_start() {
    _START=$(date +%s)
}

aa_check() {
    local now duration logs

    now=$(date +%s)
    duration=$((now - _START + 1))
    logs=$(aa-log --raw --systemd --since "-${duration}s")
    for pattern in "${IGNORE[@]}"; do
        logs=$(echo "$logs" | grep -v "$pattern")
    done

    aa_start
    if [[ -n "$logs" ]]; then
        fail "profile $PROGRAM raised logs: $logs"
    fi
}

_timeout() {
    local duration="2s"
    timeout --preserve-status --kill-after="$duration" "$duration" "$@"
}

# Bats setup and teardown hooks

setup_file() {
    aa_setup
}

teardown() {
	aa_check
}
