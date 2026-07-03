---
title: User
tags:
  - tunables
  - default
---

## home

### @{HOMEDIRS}

[:material-file-eye:](https://gitlab.com/apparmor/apparmor/-/blob/9a402e9e6a04b45dd788c7dfb0106dae45443b01/profiles/apparmor.d/tunables/home#L15 "View source"){ .abs-source }

@{HOMEDIRS} is a space-separated list of where user home directories
are stored, for programs that must enumerate all home directories on a
system.

```
@{HOMEDIRS}=/home/
```
### @{HOME}

[:material-file-eye:](https://gitlab.com/apparmor/apparmor/-/blob/9a402e9e6a04b45dd788c7dfb0106dae45443b01/profiles/apparmor.d/tunables/home#L21 "View source"){ .abs-source }

@{HOME} is a space-separated list of all user home directories. While
it doesn't refer to a specific home directory (AppArmor doesn't
enforce discretionary access controls) it can be used as if it did
refer to a specific home directory

```
@{HOME}=@{HOMEDIRS}/*/ /root/
```

## rygel

### @{rygel_media_dirs}

[:material-file-eye:](https://gitlab.com/apparmor/apparmor/-/blob/9a402e9e6a04b45dd788c7dfb0106dae45443b01/profiles/apparmor.d/tunables/rygel#L10 "View source"){ .abs-source }

allow standard XDG media paths by default

```
@{rygel_media_dirs}=@{HOME}/@{XDG_MUSIC_DIR} @{HOME}/@{XDG_VIDEOS_DIR} @{HOME}/@{XDG_PICTURES_DIR}
```

## share

### @{flatpak_exports_root}

[:material-file-eye:](https://gitlab.com/apparmor/apparmor/-/blob/9a402e9e6a04b45dd788c7dfb0106dae45443b01/profiles/apparmor.d/tunables/share#L1 "View source"){ .abs-source }

```
@{flatpak_exports_root} = {flatpak/exports,flatpak/{app,runtime}/*/*/*/*/export}
```
### @{system_share_dirs}

[:material-file-eye:](https://gitlab.com/apparmor/apparmor/-/blob/9a402e9e6a04b45dd788c7dfb0106dae45443b01/profiles/apparmor.d/tunables/share#L8 "View source"){ .abs-source }

System-wide directories with behaviour analogous to /usr/share
in patterns like the freedesktop.org basedir spec. These are
owned by root or a system user, appear in XDG_DATA_DIRS, and
are the parent directory for `applications`, `themes`,
`dbus-1/services`, etc.

```
@{system_share_dirs} = /{usr,usr/local,var/lib/@{flatpak_exports_root}}/share
```
### @{user_share_dirs}

[:material-file-eye:](https://gitlab.com/apparmor/apparmor/-/blob/9a402e9e6a04b45dd788c7dfb0106dae45443b01/profiles/apparmor.d/tunables/share#L15 "View source"){ .abs-source }

Per-user/personal directories with behaviour analogous to
~/.local/share in patterns like the freedesktop.org basedir spec.
These are owned by the user running an application, appear in
XDG_DATA_DIRS or XDG_DATA_HOME, and are the parent directory
for the same subdirectories as @{system_share_dirs}

```
@{user_share_dirs} = @{HOME}/.local{,/share/@{flatpak_exports_root}}/share
```

## xdg-user-dirs

### @{XDG_DESKTOP_DIR}

[:material-file-eye:](https://gitlab.com/apparmor/apparmor/-/blob/9a402e9e6a04b45dd788c7dfb0106dae45443b01/profiles/apparmor.d/tunables/xdg-user-dirs#L13 "View source"){ .abs-source }

Define the common set of XDG user directories (usually defined in
/etc/xdg/user-dirs.defaults)

```
@{XDG_DESKTOP_DIR}="Desktop"
@{XDG_DOWNLOAD_DIR}="Downloads"
@{XDG_TEMPLATES_DIR}="Templates"
@{XDG_PUBLICSHARE_DIR}="Public"
@{XDG_DOCUMENTS_DIR}="Documents"
@{XDG_MUSIC_DIR}="Music"
@{XDG_PICTURES_DIR}="Pictures"
@{XDG_VIDEOS_DIR}="Videos"
```
