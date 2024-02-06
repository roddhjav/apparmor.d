# apparmor.d - Full set of apparmor profiles
# Extended user XDG directories definition
# Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

# To allow extended personalisation by the user without breaking everything.
# All apparmor profiles should always use the variables defined here.

# XDG_*_DIR variables are relative pathnames from the user home directory. 
# user_*_dirs variables are absolute path.

# Second part. First part in /etc/apparmor.d/tunables/home.d/apparmor.d

# Other user directories
@{user_documents_dirs}=@{HOME}/@{XDG_DOCUMENTS_DIR} @{MOUNTS}/@{XDG_DOCUMENTS_DIR}
@{user_download_dirs}=@{HOME}/@{XDG_DOWNLOAD_DIR} @{MOUNTS}/@{XDG_DOWNLOAD_DIR}
@{user_music_dirs}=@{HOME}/@{XDG_MUSIC_DIR} @{MOUNTS}/@{XDG_MUSIC_DIR}
@{user_pictures_dirs}=@{HOME}/@{XDG_PICTURES_DIR} @{MOUNTS}/@{XDG_PICTURES_DIR}
@{user_publicshare_dirs}=@{HOME}/@{XDG_PUBLICSHARE_DIR} @{MOUNTS}/@{XDG_PUBLICSHARE_DIR}
@{user_templates_dirs}=@{HOME}/@{XDG_TEMPLATES_DIR} @{MOUNTS}/@{XDG_TEMPLATES_DIR}
@{user_videos_dirs}=@{HOME}/@{XDG_VIDEOS_DIR} @{MOUNTS}/@{XDG_VIDEOS_DIR}
@{user_vm_shares}=@{HOME}/@{XDG_VM_SHARES_DIR} @{MOUNTS}/@{XDG_VM_SHARES_DIR}

include if exists <tunables/xdg-user-dirs.d/apparmor.d.d>
