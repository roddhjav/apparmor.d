---
title: Variables References
---

## XDG directories

### User directories

| Description | Name | Default Value(s) |
|-------------|:----:|---------------|
| Desktop | `@{XDG_DESKTOP_DIR}` | `Desktop` |
| Documents | `@{XDG_DOCUMENTS_DIR}` | `Documents` |
| Downloads | `@{XDG_DOWNLOAD_DIR}` | `Downloads` |
| Music | `@{XDG_MUSIC_DIR}` | `Music` |
| Pictures | `@{XDG_PICTURES_DIR}` | `Pictures` |
| Videos | `@{XDG_VIDEOS_DIR}` | `Videos` |
| Screenshots | `@{XDG_SCREENSHOTS_DIR}` | `@{XDG_PICTURES_DIR}/Screenshots` |
| Wallpapers | `@{XDG_WALLPAPERS_DIR}` | `@{XDG_PICTURES_DIR}/Wallpapers` |
| Books | `@{XDG_BOOKS_DIR}` | `Books` |
| Games | `@{XDG_GAMES_DIR}` | `.games` |
| Templates | `@{XDG_TEMPLATES_DIR}` | `Templates` |
| Public | `@{XDG_PUBLICSHARE_DIR}` | `Public` |
| Projects | `@{XDG_PROJECTS_DIR}` | `Projects` |
| Private | `@{XDG_PRIVATE_DIR}` | `.{p,P}rivate {p,P}rivate` |
| Work | `@{XDG_WORK_DIR}` | `Work` |
| Mail | `@{XDG_MAIL_DIR}` | `Mail .{m,M}ail` |
| Sync | `@{XDG_SYNC_DIR}` | `Sync` |
| Torrents | `@{XDG_TORRENTS_DIR}` | `Torrents` |
| Vm | `@{XDG_VM_DIR}` | `.vm`
| Vm Shares | `@{XDG_VM_SHARES_DIR}` | `VM_Shares`
| Disk images | `@{XDG_IMG_DIR}` | `images` |

### Dotfiles

| Description | Name | Default Value(s) |
|-------------|:----:|---------------|
| Cache | ` @{XDG_CACHE_DIR}` | `.cache` |
| Config | `@{XDG_CONFIG_DIR}` | `.config` |
| Data | `@{XDG_DATA_DIR}` | `.local/share` |
| State | `@{XDG_STATE_DIR}` | `.local/state` |
| Bin | `@{XDG_BIN_DIR}` | `.local/bin` |
| Lib | `@{XDG_LIB_DIR}` | `.local/lib` |
| GPG | `@{XDG_GPG_DIR}` | `.gnupg` |
| SSH | `@{XDG_SSH_DIR}` | `.ssh` |
| Private | `@{XDG_PRIVATE_DIR}` | `.{p,P}rivate {p,P}rivate` |
| Passwords | `@{XDG_PASSWORD_STORE_DIR}` | `.password-store` |
| Mail | `@{XDG_MAIL_DIR}` | `Mail .{m,M}ail` |

### Full configuration path

| Description | Name | Default Value(s) |
|-------------|:----:|---------------|
| Cache | `@{user_cache_dirs}` | `@{HOME}/@{XDG_CACHE_DIR}` |
| Config | `@{user_config_dirs}` | `@{HOME}/@{XDG_CONFIG_DIR}` |
| Bin | `@{user_bin_dirs}` | `@{HOME}/@{XDG_BIN_DIR}` |
| Lib | `@{user_lib_dirs}` | `@{HOME}/@{XDG_LIB_DIR}` |
| Share | `@{user_share_dirs}` | ` @{HOME}/@{XDG_DATA_DIR}` |
| State | `@{user_state_dirs}` | ` @{HOME}/@{XDG_STATE_DIR}` |
| Build | `@{user_build_dirs}` | `/tmp/` |
| Packages | `@{user_pkg_dirs}` | `/tmp/pkg/` |
| Tmp | `@{user_tmp_dirs}` | `@{run}/user/@{uid} /tmp/` |

### Full user path

| Description | Name | Default Value(s) |
|-------------|:----:|---------------|
| Documents | `@{user_documents_dirs}` | `@{HOME}/@{XDG_DOCUMENTS_DIR} @{MOUNTS}/@{XDG_DOCUMENTS_DIR}` |
| Downloads | `@{user_download_dirs}` | `@{HOME}/@{XDG_DOWNLOAD_DIR} @{MOUNTS}/@{XDG_DOWNLOAD_DIR}` |
| Music | `@{user_music_dirs}` | `@{HOME}/@{XDG_MUSIC_DIR} @{MOUNTS}/@{XDG_MUSIC_DIR}` |
| Pictures | `@{user_pictures_dirs}` | `@{HOME}/@{XDG_PICTURES_DIR} @{MOUNTS}/@{XDG_PICTURES_DIR}` |
| Videos | `@{user_videos_dirs}` | `@{HOME}/@{XDG_VIDEOS_DIR} @{MOUNTS}/@{XDG_VIDEOS_DIR}` |
| Books | `@{user_books_dirs}` | `@{HOME}/@{XDG_BOOKS_DIR} @{MOUNTS}/@{XDG_BOOKS_DIR}` |
| Games | `@{user_games_dirs}` | `@{HOME}/@{XDG_GAMES_DIR} @{MOUNTS}/@{XDG_GAMES_DIR}` |
| Private | `@{user_private_dirs}` | `@{HOME}/@{XDG_PRIVATE_DIR} @{MOUNTS}/@{XDG_PRIVATE_DIR}` |
| Passwords | `@{user_password_store_dirs}` | `@{HOME}/@{XDG_PASSWORD_STORE_DIR} @{MOUNTS}/@{XDG_PASSWORD_STORE_DIR}` |
| Work | `@{user_work_dirs}` | `@{HOME}/@{XDG_WORK_DIR} @{MOUNTS}/@{XDG_WORK_DIR}` |
| Mail | `@{user_mail_dirs}` | `@{HOME}/@{XDG_MAIL_DIR} @{MOUNTS}/@{XDG_MAIL_DIR}` |
| Projects | `@{user_projects_dirs}` | `@{HOME}/@{XDG_PROJECTS_DIR} @{MOUNTS}/@{XDG_PROJECTS_DIR}` |
| Public | `@{user_publicshare_dirs}` | `@{HOME}/@{XDG_PUBLICSHARE_DIR} @{MOUNTS}/@{XDG_PUBLICSHARE_DIR}` |
| Templates | `@{user_templates_dirs}` | `@{HOME}/@{XDG_TEMPLATES_DIR} @{MOUNTS}/@{XDG_TEMPLATES_DIR}` |
| Torrents | `@{user_torrents_dirs}` | `@{HOME}/@{XDG_TORRENTS_DIR} @{MOUNTS}/@{XDG_TORRENTS_DIR}` |
| Sync | `@{user_sync_dirs}` | `@{HOME}/@{XDG_SYNC_DIR} @{MOUNTS}/*/@{XDG_SYNC_DIR}` |
| Vm | `@{user_vm_dirs}` | `@{HOME}/@{XDG_VM_DIR} @{MOUNTS}/@{XDG_VM_DIR}`
| Vm Shares | `@{user_vm_shares}` | `@{HOME}/@{XDG_VM_DIR} @{MOUNTS}/@{XDG_VM_DIR}`
| Disk images | `@{user_img_dirs}` | `@{HOME}/@{XDG_VM_SHARES_DIR} @{MOUNTS}/@{XDG_VM_SHARES_DIR}` |


## System variables

!!! warning

    Do not modify these variables unless you know what you are doing

**Helper variables**

| Description | Name | Default Value(s) |
|-------------|:----:|---------------|
| Integer (up to 10 digits) | `@{int}` | `[0-9]{[0-9],}{[0-9],}{[0-9],}{[0-9],}{[0-9],}{[0-9],}{[0-9],}{[0-9],}{[0-9],}` |
| Any 6, 8 or 10 characters | `@{rand6}`, `@{rand8}`, `@{rand10}` | |
| Hexadecimal | `@{h}*@{h}` |  |
| Universally unique identifier | `@{uuid}` |  |
| Current Process id | `@{pid}` | `[0-9]*` |
| Processes ids | `@{pids}` | `[0-9]*` |
| User id | `@{uid}` | `[0-9]*` |
| Thread id | `@{tid}` | `[0-9]*` |
| Single hexadecimal character | `@{h}` | `[0-9a-fA-F]` |
| Single alphanumeric character | `@{c}` | `[0-9a-zA-Z]` |
| PCI Devices | `@{pci}` | `@{pci_bus}/**/` |
| PCI Bus | `@{pci_bus}` | `pci@{h}@{h}@{h}@{h}:@{h}@{h}` |
| PCI Id | `@{pci_id}` | `@{h}@{h}@{h}@{h}:@{h}@{h}:@{h}@{h}.@{h}` |

**System Paths**

| Description | Name | Default Value(s) |
|-------------|:----:|---------------|
| Root Home | `@{HOMEDIRS}` | `/home/` |
| Home directories | `@{HOME}` | `@{HOMEDIRS}/*/ /root/` |
| Root Mountpoints | `@{MOUNTDIRS}` | `/media/ @{run}/media/ /mnt/` |
| Mountpoints directories | `@{MOUNTS}` | `@{MOUNTDIRS}/*/` |
| Bin | `@{bin}` |  `/{usr/,}{s,}bin` |
| Lib | `@{lib}` |  `/{usr/,}lib{,exec,32,64}` |
| multi-arch library | `@{multiarch}` | `*-linux-gnu*` |
| Proc | `@{PROC}` | `/proc/` |
| Run | `@{run}` | `/run/ /var/run/` |
| Sys | `@{sys}` | `/sys/` |
| System wide share | `@{system_share_dirs}` | `/{usr,usr/local,var/lib/@{flatpak_exports_root}}/share` |
| Flatpak export | `@{flatpak_exports_root}` | `{flatpak/exports,flatpak/{app,runtime}/*/*/*/*/export}` |

**Program paths**

| Description | Name | Default Value(s) |
|-------------|:----:|---------------|
| All the shells | `@{shells}` | `sh zsh bash dash fish rbash ksh tcsh csh` |
| Shells path | `@{shells_path}` | `@{bin}/@{shells}` |
| Coreutils programs that should not have dedicated profile | `@{coreutils}` | See [tunables/multiarch.d/paths](https://github.com/roddhjav/apparmor.d/blob/c2d88c9bffc626fcf7d9b15b42b50706afb29562/apparmor.d/tunables/multiarch.d/paths#L46) |
| Coreutils paths | `@{coreutils_path}` | `@{bin}/@{coreutils}` |
| Launcher paths | `@{open_path}` | `@{bin}/exo-open @{bin}/xdg-open @{lib}/@{multiarch}/glib-[0-9]*/gio-launch-desktop @{lib}/gio-launch-desktop`
| All browser paths | `@{*_path}` | See [tunables/multiarch.d/paths](https://github.com/roddhjav/apparmor.d/blob/c2d88c9bffc626fcf7d9b15b42b50706afb29562/apparmor.d/tunables/multiarch.d/paths#L11)
