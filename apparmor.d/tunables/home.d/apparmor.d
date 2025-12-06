# apparmor.d - Full set of apparmor profiles
# Extended user XDG directories definition
# Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

# To allow extended personalisation by the user without breaking everything.
# All apparmor profiles should always use the variables defined here.

# XDG_*_DIR variables are relative pathnames from the user home directory. 
# user_*_dirs variables are absolute path.

# First part, second part in /etc/apparmor.d/tunables/xdg-user-dirs.d/apparmor.d

# Define the XDG Base Directory
@{XDG_CACHE_DIR}=".cache"
@{XDG_CONFIG_DIR}=".config"
@{XDG_DATA_DIR}=".local/share"
@{XDG_STATE_DIR}=".local/state"
@{XDG_BIN_DIR}=".local/bin"
@{XDG_LIB_DIR}=".local/lib"

# Define extended user directories not defined in the XDG standard but commonly
# used in profiles
@{XDG_BOOKS_DIR}="Books"
@{XDG_GAMES_DIR}="Games"
@{XDG_PROJECTS_DIR}="Projects"
@{XDG_WORK_DIR}="Work"
@{XDG_MAIL_DIR}="Mail" ".{m,M}ail"
@{XDG_SYNC_DIR}="Sync"
@{XDG_TORRENTS_DIR}="Torrents"
@{XDG_GAMESSTUDIO_DIR}="unity3d"

# Define user directories for virtual machines, shared folders and disk images
@{XDG_VM_DIR}=".vm"
@{XDG_VMSHARE_DIR}=".vmshare"
@{XDG_IMG_DIR}=".img"

# Define user build directories and artifacts output
@{XDG_BUILD_DIR}=".build"
@{XDG_PKG_DIR}=".pkg"

# Define user personal keyrings
@{XDG_GPG_DIR}=".gnupg"
@{XDG_SSH_DIR}=".ssh"
@{XDG_PASSWORD_STORE_DIR}=".password-store"

# Define user personal private directories
@{XDG_PRIVATE_DIR}=".{p,P}rivate" "{p,P}rivate"

# Full path of the XDG Base Directory 
@{user_cache_dirs}=@{HOME}/@{XDG_CACHE_DIR}
@{user_config_dirs}=@{HOME}/@{XDG_CONFIG_DIR}
@{user_state_dirs}=@{HOME}/@{XDG_STATE_DIR}
@{user_bin_dirs}=@{HOME}/@{XDG_BIN_DIR}
@{user_lib_dirs}=@{HOME}/@{XDG_LIB_DIR}

# Other user directories
@{user_books_dirs}=@{HOME}/@{XDG_BOOKS_DIR} @{MOUNTS}/@{XDG_BOOKS_DIR}
@{user_games_dirs}=@{HOME}/@{XDG_GAMES_DIR} @{MOUNTS}/@{XDG_GAMES_DIR}
@{user_projects_dirs}=@{HOME}/@{XDG_PROJECTS_DIR} @{MOUNTS}/@{XDG_PROJECTS_DIR}
@{user_work_dirs}=@{HOME}/@{XDG_WORK_DIR} @{MOUNTS}/@{XDG_WORK_DIR}
@{user_mail_dirs}=@{HOME}/@{XDG_MAIL_DIR} @{MOUNTS}/@{XDG_MAIL_DIR}
@{user_sync_dirs}=@{HOME}/@{XDG_SYNC_DIR} @{MOUNTS}/@{XDG_SYNC_DIR}
@{user_torrents_dirs}=@{HOME}/@{XDG_TORRENTS_DIR} @{MOUNTS}/@{XDG_TORRENTS_DIR}
@{user_vm_dirs}=@{HOME}/@{XDG_VM_DIR} @{MOUNTS}/@{XDG_VM_DIR}
@{user_vmshare_dirs}=@{HOME}/@{XDG_VMSHARE_DIR} @{MOUNTS}/@{XDG_VMSHARE_DIR}
@{user_img_dirs}=@{HOME}/@{XDG_IMG_DIR} @{MOUNTS}/@{XDG_IMG_DIR}
@{user_build_dirs}=@{HOME}/@{XDG_BUILD_DIR} @{MOUNTS}/@{XDG_BUILD_DIR}
@{user_pkg_dirs}=@{HOME}/@{XDG_PKG_DIR} @{MOUNTS}/@{XDG_PKG_DIR}
@{user_gpg_dirs}=@{HOME}/@{XDG_GPG_DIR} @{MOUNTS}/@{XDG_GPG_DIR}
@{user_ssh_dirs}=@{HOME}/@{XDG_SSH_DIR} @{MOUNTS}/@{XDG_SSH_DIR}
@{user_passwordstore_dirs}=@{HOME}/@{XDG_PASSWORD_STORE_DIR} @{MOUNTS}/@{XDG_PASSWORD_STORE_DIR}
@{user_private_dirs}=@{HOME}/@{XDG_PRIVATE_DIR} @{MOUNTS}/@{XDG_PRIVATE_DIR}

# Similar system-wide paths
@{system_games_dirs}=/usr/games /var/lib/games

# vim:syntax=apparmor
