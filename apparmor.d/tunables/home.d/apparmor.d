# apparmor.d - Full set of apparmor profiles
# Extended user XDG directories definition
# Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

# To allow extended personalisation by the user without breaking everything.
# All apparmor profiles should always use the variables defined here.

# XDG_*_DIR variables are relative pathnames from the user home directory. 
# user_*_dirs variables are absolute path.

# First part, second part in /etc/apparmor.d/tunables/xdg-user-dirs.d/apparmor.d

# Extra user personal directories
@{XDG_BOOKS_DIR}="Books"
@{XDG_PROJECTS_DIR}="Projects"
@{XDG_WORK_DIR}="Work"
@{XDG_SYNC_DIR}="Sync"
@{XDG_TORRENTS_DIR}="Torrents"
@{XDG_GAMES_DIR}=".games"
@{XDG_VM_DIR}=".vm"
@{XDG_VM_SHARES_DIR}="VM_Shares"
@{XDG_IMG_DIR}="images"
@{XDG_MAIL_DIR}="Mail"
@{XDG_SCREENSHOTS_DIR}="Pictures/Screenshots"
@{XDG_WALLPAPERS_DIR}="Pictures/Wallpapers"

# User personal keyrings
@{XDG_SSH_DIR}=".ssh"
@{XDG_GPG_DIR}=".gnupg"
@{XDG_PASSWORD_STORE_DIR}=".password-store"

# Definition of local user configuration directories
@{XDG_CACHE_DIR}=".cache"
@{XDG_CONFIG_DIR}=".config"
@{XDG_DATA_DIR}=".local/share"
@{XDG_STATE_DIR}=".local/state"
@{XDG_BIN_DIR}=".local/bin"
@{XDG_LIB_DIR}=".local/lib"

# Full path of the user configuration directories
@{user_cache_dirs}=@{HOME}/@{XDG_CACHE_DIR}
@{user_config_dirs}=@{HOME}/@{XDG_CONFIG_DIR}
@{user_state_dirs}=@{HOME}/@{XDG_STATE_DIR}
@{user_bin_dirs}=@{HOME}/@{XDG_BIN_DIR}
@{user_lib_dirs}=@{HOME}/@{XDG_LIB_DIR}

# User build directories and output
@{user_build_dirs}="/tmp/build/"
@{user_pkg_dirs}="/tmp/pkg/"
@{user_tmp_dirs}=@{run}/user/@{uid} /tmp/
@{user_img_dirs}=@{HOME}/@{XDG_IMG_DIR} @{MOUNTS}/@{XDG_IMG_DIR}

# Other user directories
@{user_books_dirs}=@{HOME}/@{XDG_BOOKS_DIR} @{MOUNTS}/@{XDG_BOOKS_DIR}
@{user_games_dirs}=@{HOME}/@{XDG_GAMES_DIR} @{MOUNTS}/@{XDG_GAMES_DIR}
@{user_mail_dirs}=@{HOME}/@{XDG_MAIL_DIR} @{MOUNTS}/@{XDG_MAIL_DIR}
@{user_projects_dirs}=@{HOME}/@{XDG_PROJECTS_DIR} @{MOUNTS}/@{XDG_PROJECTS_DIR}
@{user_sync_dirs}=@{HOME}/@{XDG_SYNC_DIR} @{MOUNTS}/*/@{XDG_SYNC_DIR}
@{user_torrents_dirs}=@{HOME}/@{XDG_TORRENTS_DIR} @{MOUNTS}/@{XDG_TORRENTS_DIR}
@{user_vm_dirs}=@{HOME}/@{XDG_VM_DIR} @{MOUNTS}/@{XDG_VM_DIR}
@{user_work_dirs}=@{HOME}/@{XDG_WORK_DIR} @{MOUNTS}/@{XDG_WORK_DIR}
@{user_password_store_dirs}=@{HOME}/@{XDG_PASSWORD_STORE_DIR} @{MOUNTS}/@{XDG_PASSWORD_STORE_DIR}
