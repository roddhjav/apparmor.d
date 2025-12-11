---
title: Variables References
---

## XDG directories

### User directories

<figure markdown>

| Description | Name | Default Value(s) |
|-------------|------|------------------|
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
| Vm | `@{XDG_VM_DIR}` | `.vm` |
| Vm Shares | `@{XDG_VMSHARE_DIR}` | `VM_Shares` |
| Disk images | `@{XDG_IMG_DIR}` | `images` |
| Games Studio | `@{XDG_GAMESSTUDIO_DIR}` | `.unity3d` |

</figure>

### Dotfiles

<figure markdown>

| Description | Name | Default Value(s) |
|-------------|------|------------------|
| Cache | ` @{XDG_CACHE_DIR}` | `.cache` |
| Config | `@{XDG_CONFIG_DIR}` | `.config` |
| Data | `@{XDG_DATA_DIR}` | `.local/share` |
| State | `@{XDG_STATE_DIR}` | `.local/state` |
| Bin | `@{XDG_BIN_DIR}` | `.local/bin` |
| Lib | `@{XDG_LIB_DIR}` | `.local/lib` |
| GPG | `@{XDG_GPG_DIR}` | `.gnupg` |
| SSH | `@{XDG_SSH_DIR}` | `.ssh` |
| Private | `@{XDG_PRIVATE_DIR}` | `.{p,P}rivate {p,P}rivate` |
| Passwords | `@{XDG_PASSWORDSTORE_DIR}` | `.password-store` |

</figure>

### Full configuration path

<figure markdown>

| Description | Name | Default Value(s) |
|-------------|------|------------------|
| Cache | `@{user_cache_dirs}` | `@{HOME}/@{XDG_CACHE_DIR}` |
| Config | `@{user_config_dirs}` | `@{HOME}/@{XDG_CONFIG_DIR}` |
| Bin | `@{user_bin_dirs}` | `@{HOME}/@{XDG_BIN_DIR}` |
| Lib | `@{user_lib_dirs}` | `@{HOME}/@{XDG_LIB_DIR}` |
| Share | `@{user_share_dirs}` | ` @{HOME}/@{XDG_DATA_DIR}` |
| State | `@{user_state_dirs}` | ` @{HOME}/@{XDG_STATE_DIR}` |
| Build | `@{user_build_dirs}` | `/tmp/build/` |
| Packages | `@{user_pkg_dirs}` | `/tmp/pkg/` |

</figure>

### Full user path

<figure markdown>

| Description | Name | Default Value(s) |
|-------------|------|------------------|
| Documents | `@{user_documents_dirs}` | `@{HOME}/@{XDG_DOCUMENTS_DIR} @{MOUNTS}/@{XDG_DOCUMENTS_DIR}` |
| Downloads | `@{user_download_dirs}` | `@{HOME}/@{XDG_DOWNLOAD_DIR} @{MOUNTS}/@{XDG_DOWNLOAD_DIR}` |
| Music | `@{user_music_dirs}` | `@{HOME}/@{XDG_MUSIC_DIR} @{MOUNTS}/@{XDG_MUSIC_DIR}` |
| Pictures | `@{user_pictures_dirs}` | `@{HOME}/@{XDG_PICTURES_DIR} @{MOUNTS}/@{XDG_PICTURES_DIR}` |
| Videos | `@{user_videos_dirs}` | `@{HOME}/@{XDG_VIDEOS_DIR} @{MOUNTS}/@{XDG_VIDEOS_DIR}` |
| Books | `@{user_books_dirs}` | `@{HOME}/@{XDG_BOOKS_DIR} @{MOUNTS}/@{XDG_BOOKS_DIR}` |
| Games | `@{user_games_dirs}` | `@{HOME}/@{XDG_GAMES_DIR} @{MOUNTS}/@{XDG_GAMES_DIR}` |
| Private | `@{user_private_dirs}` | `@{HOME}/@{XDG_PRIVATE_DIR} @{MOUNTS}/@{XDG_PRIVATE_DIR}` |
| Passwords | `@{user_passwordstore_dirs}` | `@{HOME}/@{XDG_PASSWORDSTORE_DIR} @{MOUNTS}/@{XDG_PASSWORDSTORE_DIR}` |
| Work | `@{user_work_dirs}` | `@{HOME}/@{XDG_WORK_DIR} @{MOUNTS}/@{XDG_WORK_DIR}` |
| Mail | `@{user_mail_dirs}` | `@{HOME}/@{XDG_MAIL_DIR} @{MOUNTS}/@{XDG_MAIL_DIR}` |
| Projects | `@{user_projects_dirs}` | `@{HOME}/@{XDG_PROJECTS_DIR} @{MOUNTS}/@{XDG_PROJECTS_DIR}` |
| Public | `@{user_publicshare_dirs}` | `@{HOME}/@{XDG_PUBLICSHARE_DIR} @{MOUNTS}/@{XDG_PUBLICSHARE_DIR}` |
| Templates | `@{user_templates_dirs}` | `@{HOME}/@{XDG_TEMPLATES_DIR} @{MOUNTS}/@{XDG_TEMPLATES_DIR}` |
| Torrents | `@{user_torrents_dirs}` | `@{HOME}/@{XDG_TORRENTS_DIR} @{MOUNTS}/@{XDG_TORRENTS_DIR}` |
| Sync | `@{user_sync_dirs}` | `@{HOME}/@{XDG_SYNC_DIR} @{MOUNTS}/*/@{XDG_SYNC_DIR}` |
| Vm | `@{user_vm_dirs}` | `@{HOME}/@{XDG_VM_DIR} @{MOUNTS}/@{XDG_VM_DIR}` |
| Vm Shares | `@{user_vmshare_dirs}` | `@{HOME}/@{XDG_VMSHARE_DIR} @{MOUNTS}/@{XDG_VMSHARE_DIR}` |
| Disk images | `@{user_img_dirs}` | `@{HOME}/@{XDG_IMG_DIR} @{MOUNTS}/@{XDG_IMG_DIR}` |

</figure>


## System variables

!!! danger

    Do not modify these variables unless you know what you are doing

#### Base variables

<figure markdown>

| Description | Name | Default Value(s) |
|-------------|------|------------------|
| Any digit | `@{d}` | `[0-9]` |
| Any letter | `@{l}` | `[a-zA-Z]` |
| Single alphanumeric character | `@{c}` | `[0-9a-zA-Z]` |
| Word character: matches any letter, digit or underscore. | `@{w}` | `[0-9a-zA-Z_]` |
| Single hexadecimal character | `@{h}` | `[0-9a-fA-F]` |
| Integer up to 10 digits (0-9999999999) | `@{int}` | `@{d}{@{d},}{@{d},}{@{d},}{@{d},}{@{d},}{@{d},}{@{d},}{@{d},}{@{d},}` |
| Unsigned integer over 8 bits (0-255) | `@{u8}` | `[0-9]{[0-9],} 1[0-9][0-9] 2[0-4][0-9] 25[0-5]` |
| Unsigned integer over 16 bits (0-65535, 5 digits) | `@{u16}` | `@{d}{@{d},}{@{d},}{@{d},}{@{d},}` |
| Hexadecimal up to 64 characters | `@{hex}` | |
| Alphanumeric up to 64 characters | `@{rand}` | |
| Word up to 64 characters | `@{word}` | |

</figure>

#### Basic variables of a given length

<figure markdown>

| Description | Name |
|-------------|------|
| Any x digits characters | `@{int2}` `@{int4}` `@{int6}` `@{int8}` `@{int9}` `@{int10}` `@{int12}` `@{int15}` `@{int16}` `@{int32}` `@{int64}` |
| Any x hexadecimal characters | `@{hex2}` `@{hex4}` `@{hex6}` `@{hex8}` `@{hex9}` `@{hex10}` `@{hex12}` `@{hex15}` `@{hex16}` `@{hex32}` `@{hex38}` `@{hex64}` |
| Any x alphanumeric characters | `@{rand2}` `@{rand4}` `@{rand6}` `@{rand8}` `@{rand9}` `@{rand10}` `@{rand12}` `@{rand15}` `@{rand16}` `@{rand32}` `@{rand64}` |
| Any x word characters | `@{word2}` `@{word4}` `@{word6}` `@{word8}` `@{word9}` `@{word10}` `@{word12}` `@{word15}` `@{word16}` `@{word32}` `@{word64}` |

</figure>

#### System Variables

<figure markdown>

| Description | Name | Default Value(s) |
|-------------|------|------------------|
| Common architecture names | `@{arch}` | `x86_64 amd64 i386 i686` |
| Dbus unique name | `@{busname}` | `:1.@{u16} :not.active.yet` |
| Universally unique identifier | `@{uuid}` | `@{h}@{h}@{h}@{h}@{h}@{h}@{h}@{h}[-_]@{h}@{h}@{h}@{h}[-_]@{h}@{h}@{h}@{h}[-_]@{h}@{h}@{h}@{h}[-_]@{h}@{h}@{h}@{h}@{h}@{h}@{h}@{h}@{h}@{h}@{h}@{h}` |
| Username valid characters | `@{user}` | `[a-zA-Z_]{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}` |
| Group valid characters | `@{group}` | `@{user}` |
| Semantic version | `@{version}` | `@{int}{.@{int},}{.@{int},}{-@{rand},}` |
| Current Process Id | `@{pid}` | `{[1-9],[1-9][0-9],[1-9][0-9][0-9],[1-9][0-9][0-9][0-9],[1-9][0-9][0-9][0-9][0-9],[1-9][0-9][0-9][0-9][0-9][0-9],[1-4][0-9][0-9][0-9][0-9][0-9][0-9]}` |
| Processes Ids | `@{pids}` | `@{pid}` |
| Thread Id | `@{tid}` | `@{pid}` |
| User Id (equivalent to `@{int}`) | `@{uid}` | `{[0-9],[1-9][0-9],[1-9][0-9][0-9],[1-9][0-9][0-9][0-9],[1-9][0-9][0-9][0-9][0-9],[1-9][0-9][0-9][0-9][0-9][0-9],[1-9][0-9][0-9][0-9][0-9][0-9][0-9],[1-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9],[1-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9],[1-4][0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9]}` |

</figure>

#### System Paths

<figure markdown>

| Description | Name | Default Value(s) |
|-------------|------|------------------|
| Root Home | `@{HOMEDIRS}` | `/home/` |
| Home directories | `@{HOME}` | `@{HOMEDIRS}/*/ /root/` |
| Root Mountpoints | `@{MOUNTDIRS}` | `/media/ @{run}/media/@{user}/ /mnt/` |
| Mountpoints directories | `@{MOUNTS}` | `@{MOUNTDIRS}/*/ @{run}/user/@{uid}/gvfs/` |
| Bin | `@{bin}` |  `/{usr/,}bin` |
| Sbin | `@{sbin}` |  `/{usr/,}sbin` |
| Lib | `@{lib}` |  `/{usr/,}lib{,exec,32,64}` |
| multi-arch library | `@{multiarch}` | `*-linux-gnu*` |
| Proc | `@{PROC}` | `/proc/` |
| Run | `@{run}` | `/run/ /var/run/` |
| Sys | `@{sys}` | `/sys/` |
| System wide share | `@{system_share_dirs}` | `/{usr,usr/local,var/lib/@{flatpak_exports_root}}/share` |
| Flatpak export | `@{flatpak_exports_root}` | `{flatpak/exports,flatpak/{app,runtime}/*/*/*/*/export}` |

</figure>

#### System Internal

| Description | Name | Default Value(s) |
|-------------|------|------------------|
| PCI Devices | `@{pci}` | `@{pci_bus}/**/` |
| PCI Bus | `@{pci_bus}` | `pci@{h}@{h}@{h}@{h}:@{h}@{h}` |
| PCI Id | `@{pci_id}` | `@{h}@{h}@{h}@{h}:@{h}@{h}:@{h}@{h}.@{h}` |
| HCI devices | `@{hci_id}` | `dev_@{c}@{c}_@{c}@{c}_@{c}@{c}_@{c}@{c}_@{c}@{c}_@{c}@{c}` |
| Udev data dynamic assignment ranges (234 to 254 then 384 to 511) | `@{dynamic}` | `23[4-9] 24[0-9] 25[0-4] 38[4-9] 39[0-9] 4[0-9][0-9] 50[0-9] 51[0-1]` |

#### Program paths

<figure markdown>

| Description | Name | Default Value(s) |
|-------------|------|------------------|
| All the shells | `@{shells}` | `sh zsh bash dash fish rbash ksh tcsh csh` |
| Shells path | `@{shells_path}` | `@{bin}/@{shells}` |
| Coreutils programs that should not have dedicated profile | `@{coreutils}` | See [tunables/multiarch.d/paths](https://github.com/roddhjav/apparmor.d/blob/c2d88c9bffc626fcf7d9b15b42b50706afb29562/apparmor.d/tunables/multiarch.d/paths#L46) |
| Coreutils paths | `@{coreutils_path}` | `@{bin}/@{coreutils}` |
| Launcher paths | `@{open_path}` | `@{bin}/exo-open @{bin}/xdg-open @{lib}/@{multiarch}/glib-@{version}/gio-launch-desktop @{lib}/gio-launch-desktop`
| All browser paths | `@{*_path}` | See [tunables/multiarch.d/paths](https://github.com/roddhjav/apparmor.d/blob/c2d88c9bffc626fcf7d9b15b42b50706afb29562/apparmor.d/tunables/multiarch.d/paths#L11)

</figure>
