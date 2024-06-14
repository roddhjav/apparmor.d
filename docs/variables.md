---
title: Variables References
---

## XDG directories

### User directories

| Description | Name | Default Value(s) |
|-------------|:----:|---------------|
| Books | `@{XDG_BOOKS_DIR}` | `Books` |
| Desktop | `@{XDG_DESKTOP_DIR}` | `Desktop` |
| Disk images | `@{XDG_IMG_DIR}` | `images` |
| Documents | `@{XDG_DOCUMENTS_DIR}` | `Documents` |
| Download | `@{XDG_DOWNLOAD_DIR}` | `Downloads` |
| Games | `@{XDG_GAMES_DIR}` | `.games` |
| Music | `@{XDG_MUSIC_DIR}` | `Music` |
| Pictures | `@{XDG_PICTURES_DIR}` | `Pictures` |
| Projects | `@{XDG_PROJECTS_DIR}` | `Projects` |
| Public | `@{XDG_PUBLICSHARE_DIR}` | `Public` |
| Screenshots | `@{XDG_SCREENSHOTS_DIR}` | `@{XDG_PICTURES_DIR}/Screenshots` |
| Sync | `@{XDG_SYNC_DIR}` | `Sync` |
| Templates | `@{XDG_TEMPLATES_DIR}` | `Templates` |
| Torrents | `@{XDG_TORRENTS_DIR}` | `Torrents` |
| Videos | `@{XDG_VIDEOS_DIR}` | `Videos` |
| Vm | `@{XDG_VM_DIR}` | `.vm`
| Wallpapers | `@{XDG_WALLPAPERS_DIR}` | `@{XDG_PICTURES_DIR}/Wallpapers` |

### Dotfiles

| Description | Name | Default Value(s) |
|-------------|:----:|---------------|
| Bin | `@{XDG_BIN_DIR}` | `.local/bin` |
| Cache | ` @{XDG_CACHE_DIR}` | `.cache` |
| Config | `@{XDG_CONFIG_DIR}` | `.config` |
| Data | `@{XDG_DATA_DIR}` | `.local/share` |
| GPG | `@{XDG_GPG_DIR}` | `.gnupg` |
| Lib | `@{XDG_LIB_DIR}` | `.local/lib` |
| Passwords | `@{XDG_PASSWORD_STORE_DIR}` | `.password-store` |
| SSH | `@{XDG_SSH_DIR}` | `.ssh` |
| State | `@{XDG_STATE_DIR}` | `.local/state` |

### Full configuration path

| Description | Name | Default Value(s) |
|-------------|:----:|---------------|
| Bin | `@{user_bin_dirs}` | `@{HOME}/@{XDG_BIN_DIR}` |
| Build | `@{user_build_dirs}` | `/tmp/` |
| Cache | `@{user_cache_dirs}` | `@{HOME}/@{XDG_CACHE_DIR}` |
| Config | `@{user_config_dirs}` | `@{HOME}/@{XDG_CONFIG_DIR}` |
| Lib | `@{user_lib_dirs}` | `@{HOME}/@{XDG_LIB_DIR}` |
| Packages | `@{user_pkg_dirs}` | `/tmp/pkg/` |
| Share | `@{user_share_dirs}` | ` @{HOME}/@{XDG_DATA_DIR}` |
| State | `@{user_state_dirs}` | ` @{HOME}/@{XDG_STATE_DIR}` |
| Tmp | `@{user_tmp_dirs}` | `@{run}/user/@{uid} /tmp/` |

### Full user path

| Description | Name | Default Value(s) |
|-------------|:----:|---------------|
| Books | `@{user_books_dirs}` | `@{HOME}/@{XDG_BOOKS_DIR} @{MOUNTS}/@{XDG_BOOKS_DIR}` |
| Disk images | `@{user_img_dirs}` | `@{HOME}/@{XDG_IMG_DIR} @{MOUNTS}/@{XDG_IMG_DIR}` |
| Documents | `@{user_documents_dirs}` | `@{HOME}/@{XDG_DOCUMENTS_DIR} @{MOUNTS}/@{XDG_DOCUMENTS_DIR}` |
| Download | `@{user_download_dirs}` | `@{HOME}/@{XDG_DOWNLOAD_DIR} @{MOUNTS}/@{XDG_DOWNLOAD_DIR}` |
| Games | `@{user_games_dirs}` | `@{HOME}/@{XDG_GAMES_DIR} @{MOUNTS}/@{XDG_GAMES_DIR}` |
| Music | `@{user_music_dirs}` | `@{HOME}/@{XDG_MUSIC_DIR} @{MOUNTS}/@{XDG_MUSIC_DIR}` |
| Password | `@{user_password_store_dirs}` | `@{HOME}/@{XDG_PASSWORD_STORE_DIR} @{MOUNTS}/@{XDG_PASSWORD_STORE_DIR}` |
| Pictures | `@{user_pictures_dirs}` | `@{HOME}/@{XDG_PICTURES_DIR} @{MOUNTS}/@{XDG_PICTURES_DIR}` |
| Projects | `@{user_projects_dirs}` | `@{HOME}/@{XDG_PROJECTS_DIR} @{MOUNTS}/@{XDG_PROJECTS_DIR}` |
| Public | `@{user_publicshare_dirs}` | `@{HOME}/@{XDG_PUBLICSHARE_DIR} @{MOUNTS}/@{XDG_PUBLICSHARE_DIR}` |
| Sync | `@{user_sync_dirs}` | `@{HOME}/@{XDG_SYNC_DIR} @{MOUNTS}/*/@{XDG_SYNC_DIR}` |
| Templates | `@{user_templates_dirs}` | `@{HOME}/@{XDG_TEMPLATES_DIR} @{MOUNTS}/@{XDG_TEMPLATES_DIR}` |
| Torrents | `@{user_torrents_dirs}` | `@{HOME}/@{XDG_TORRENTS_DIR} @{MOUNTS}/@{XDG_TORRENTS_DIR}` |
| Videos | `@{user_videos_dirs}` | `@{HOME}/@{XDG_VIDEOS_DIR} @{MOUNTS}/@{XDG_VIDEOS_DIR}` |
| Vm | `@{user_vm_dirs}` | `@{HOME}/@{XDG_VM_DIR} @{MOUNTS}/@{XDG_VM_DIR}`


## System variables

!!! warning

    Do not modify these variables unless you know what you are doing

**Helper variables**

| Description | Name | Default Value(s) |
|-------------|:----:|---------------|
| Any 6, 8 or 10 characters | `@{rand6}`, `@{rand8}`, `@{rand10}` | |
| Current Process id | `@{pid}` | `[0-9]*` |
| Hexadecimal | `@{h}*@{h}` |  |
| Integer (up to 10 digits) | `@{int}` | `[0-9]{[0-9],}{[0-9],}{[0-9],}{[0-9],}{[0-9],}{[0-9],}{[0-9],}{[0-9],}{[0-9],}` |
| PCI Devices | `@{pci}` | `@{pci_bus}/**/` |
| PCI Bus | `@{pci_bus}` | `pci@{h}@{h}@{h}@{h}:@{h}@{h}` |
| PCI Id | `@{pci_id}` | `@{h}@{h}@{h}@{h}:@{h}@{h}:@{h}@{h}.@{h}` |
| Processes ids | `@{pids}` | `[0-9]*` |
| Single hexadecimal character | `@{h}` | `[0-9a-fA-F]` |
| Single alphanumeric character | `@{c}` | `[0-9a-zA-Z]` |
| Thread id | `@{tid}` | `[0-9]*` |
| Universally unique identifier | `@{uuid}` |  |
| User id | `@{uid}` | `[0-9]*` |

**System Paths**

| Description | Name | Default Value(s) |
|-------------|:----:|---------------|
| Bin | `@{bin}` |  `/{usr/,}{s,}bin` |
| Flatpack export | `@{flatpak_exports_root}` | `{flatpak/exports,flatpak/{app,runtime}/*/*/*/*/export}` |
| Home directories | `@{HOME}` | `@{HOMEDIRS}/*/ /root/` |
| Lib | `@{lib}` |  `/{usr/,}lib{,exec,32,64}` |
| Proc | `@{PROC}` | `/proc/` |
| Mountpoints directories | `@{MOUNTS}` | `@{MOUNTDIRS}/*/` |
| multi-arch library | `@{multiarch}` | `*-linux-gnu*` |
| Root Home | `@{HOMEDIRS}` | `/home/` |
| Root Mountpoints | `@{MOUNTDIRS}` | `/media/ @{run}/media/ /mnt/` |
| Run | `@{run}` | `/run/ /var/run/` |
| Sys | `@{sys}` | `/sys/` |
| System wide share | `@{system_share_dirs}` | `/{usr,usr/local,var/lib/@{flatpak_exports_root}}/share` |

**Program paths**

| Description | Name | Default Value(s) |
|-------------|:----:|---------------|
| All browser paths | `@{*_path}` | See [tunables/multiarch.d/paths](https://github.com/roddhjav/apparmor.d/blob/c2d88c9bffc626fcf7d9b15b42b50706afb29562/apparmor.d/tunables/multiarch.d/paths#L11)
| All the shells | `@{shells}` | `sh zsh bash dash fish rbash ksh tcsh csh` |
| Coreutils programs that should not have dedicated profile | `@{coreutils}` | See [tunables/multiarch.d/paths](https://github.com/roddhjav/apparmor.d/blob/c2d88c9bffc626fcf7d9b15b42b50706afb29562/apparmor.d/tunables/multiarch.d/paths#L46) |
| Coreutils paths | `@{coreutils_path}` | `@{bin}/@{coreutils}` |
| Launcher paths | `@{open_path}` | `@{bin}/exo-open @{bin}/xdg-open @{lib}/@{multiarch}/glib-[0-9]*/gio-launch-desktop @{lib}/gio-launch-desktop`
| Shells path | `@{shells_path}` | `@{bin}/@{shells}` |
