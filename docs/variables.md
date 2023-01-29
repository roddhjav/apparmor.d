---
title: Variables References
---

## XDG directories

### User directories

| Description | Name | Value |
|-------------|:----:|---------|
| Desktop | `@{XDG_DESKTOP_DIR}` | `Desktop` |
| Download | `@{XDG_DOWNLOAD_DIR}` | `Downloads` |
| Templates | `@{XDG_TEMPLATES_DIR}` | `Templates` |
| Public | `@{XDG_PUBLICSHARE_DIR}` | `Public` |
| Documents | `@{XDG_DOCUMENTS_DIR}` | `Documents` |
| Music | `@{XDG_MUSIC_DIR}` | `Music` |
| Pictures | `@{XDG_PICTURES_DIR}` | `Pictures` |
| Videos | `@{XDG_VIDEOS_DIR}` | `Videos` |
| Books | `@{XDG_BOOKS_DIR}` | `Books` |
| Projects | `@{XDG_PROJECTS_DIR}` | `Projects` |
| Screenshots | `@{XDG_SCREENSHOTS_DIR}` | `@{XDG_PICTURES_DIR}/Screenshots` |
| Sync | `@{XDG_SYNC_DIR}` | `Sync` |
| Torrents | `@{XDG_TORRENTS_DIR}` | `Torrents` |
| Vm | `@{XDG_VM_DIR}` | `.vm`
| Wallpapers | `@{XDG_WALLPAPERS_DIR}` | `@{XDG_PICTURES_DIR}/Wallpapers` |

### Dotfiles

| Description | Name | Value |
|-------------|:----:|---------|
| SSH | `@{XDG_SSH_DIR}` | `.ssh` |
| GPG | `@{XDG_GPG_DIR}` | `.gnupg` |
| Passwords | `@{XDG_PASSWORD_STORE_DIR}` | `.password-store` |
| Cache | ` @{XDG_CACHE_HOME}` | `.cache` |
| Config | `@{XDG_CONFIG_HOME}` | `.config` |
| Data | `@{XDG_DATA_HOME}` | `.local/share` |
| Bin | `@{XDG_BIN_HOME}` | `.local/bin` |
| Lib | `@{XDG_LIB_HOME}` | `.local/lib` |

### Full configuration path

| Description | Name | Value |
|-------------|:----:|---------|
| Cache | `@{user_cache_dirs}` | `@{HOME}/@{XDG_CACHE_HOME}` |
| Config | `@{user_config_dirs}` | `@{HOME}/@{XDG_CONFIG_HOME}` |
| Share | `@{user_share_dirs}` | ` @{HOME}/.local/share/` |
| Bin | `@{user_bin_dirs}` | `@{HOME}/@{XDG_BIN_HOME}` |
| Lib | `@{user_lib_dirs}` | `@{HOME}/@{XDG_LIB_HOME}` |
| Build | `@{user_build_dirs}` | `/tmp/` |
| Tmp | `@{user_tmp_dirs}` | `@{run}/user/@{uid} /tmp/` |
| Packages | `@{user_pkg_dirs}` | `/tmp/pkg/` |

### Full user path

| Description | Name | Value |
|-------------|:----:|---------|
| Books | `@{user_books_dirs}` | `@{HOME}/@{XDG_BOOKS_DIR} @{MOUNTS}/@{XDG_BOOKS_DIR}` |
| Documents | `@{user_documents_dirs}` | `@{HOME}/@{XDG_DOCUMENTS_DIR} @{MOUNTS}/@{XDG_DOCUMENTS_DIR}` |
| Download | `@{user_download_dirs}` | `@{HOME}/@{XDG_DOWNLOAD_DIR} @{MOUNTS}/@{XDG_DOWNLOAD_DIR}` |
| Music | `@{user_music_dirs}` | `@{HOME}/@{XDG_MUSIC_DIR} @{MOUNTS}/@{XDG_MUSIC_DIR}` |
| Pictures | `@{user_pictures_dirs}` | `@{HOME}/@{XDG_PICTURES_DIR} @{MOUNTS}/@{XDG_PICTURES_DIR}` |
| Projects | `@{user_projects_dirs}` | `@{HOME}/@{XDG_PROJECTS_DIR} @{MOUNTS}/@{XDG_PROJECTS_DIR}` |
| Public | `@{user_publicshare_dirs}` | `@{HOME}/@{XDG_PUBLICSHARE_DIR} @{MOUNTS}/@{XDG_PUBLICSHARE_DIR}` |
| Sync | `@{user_sync_dirs}` | `@{HOME}/@{XDG_SYNC_DIR} @{MOUNTS}/*/@{XDG_SYNC_DIR}` |
| Templates | `@{user_templates_dirs}` | `@{HOME}/@{XDG_TEMPLATES_DIR} @{MOUNTS}/@{XDG_TEMPLATES_DIR}` |
| Torrents | `@{user_torrents_dirs}` | `@{HOME}/@{XDG_TORRENTS_DIR} @{MOUNTS}/@{XDG_TORRENTS_DIR}` |
| Videos | `@{user_videos_dirs}` | `@{HOME}/@{XDG_VIDEOS_DIR} @{MOUNTS}/@{XDG_VIDEOS_DIR}` |
| Vm | `@{user_vm_dirs}` | `@{HOME}/@{XDG_VM_DIR} @{MOUNTS}/@{XDG_VM_DIR}`
| Password | `@{user_password_store_dirs}` | `@{HOME}/@{XDG_PASSWORD_STORE_DIR} @{MOUNTS}/@{XDG_PASSWORD_STORE_DIR}` |


## System variables

!!! warning

    Do not modify these variables unless you know what you are doing

| Description | Name | Value |
|-------------|:----:|---------|
| Root Home | `@{HOMEDIRS}` | `/home/` |
| Home directories | `@{HOME}` | `@{HOMEDIRS}/*/ /root/` |
| Current Process id | `@{pid}` | `[0-9]*` |
| Processes ids | `@{pids}` | `[0-9]*` |
| User id | `@{uid}` | `[0-9]*` |
| Thread id | `@{tid}` | `[0-9]*` |
| Root Mountpoints | `@{MOUNTDIRS}` | `/media/ @{run}/media/ /mnt/` |
| Mountpoints directories | `@{MOUNTS}` | `@{MOUNTDIRS}/*/` |
| Universally unique identifier | `@{uuid}` | `[0-9a-fA-F]*-[0-9a-fA-F]*-[0-9a-fA-F]*-[0-9a-fA-F]*-[0-9a-fA-F]*` |
| Hexadecimal | `@{hex}` | `[0-9a-fA-F]*` |
| Libexec *(Archlinux)* | `@{libexec}` | `/{usr/,}lib` |
| Libexec *(Debian/Ubuntu)* | `@{libexec}` | `/{usr/,}libexec` |
| multi-arch library | `@{multiarch}` | `*-linux-gnu*` |
| Proc | `@{PROC}` | `/proc/` |
| Run | `@{run}` | `/run/ /var/run/` |
| Sys | `@{sys}` | `/sys/` |
| Flatpack export | `@{flatpak_exports_root}` | `{flatpak/exports,flatpak/{app,runtime}/*/*/*/*/export}` |
| System wide share | `@{system_share_dirs}` | `/{usr,usr/local,var/lib/@{flatpak_exports_root}}/share` |
