# Contributing

You want to contribute to `apparmor.d`, **thank a lot for this.** Feedbacks, 
contributors, pull requests are all very welcome. You will find in this page all
the useful information needed to contribute.

## How to contribute?

1. If you don't have git on your machine, [install it][git].
2. Fork this repo by clicking on the fork button on the top of this page.
3. Clone the repository and go to the directory:
   ```sh
   git clone  https://github.com/this-is-you/apparmor.d.git
   cd apparmor.d
   ```
4. Create a branch:
   ```
   git checkout -b my_contribution
   ```
5. Make the changes and commit:
   ```
   git add <files changed>
   git commit -m "A message for sum up my contribution"
   ```
6. Push changes to GitHub:
   ```
   git push origin my_contribution
   ```
7. Submit your changes for review: If you go to your repository on GitHub,
you'll see a Compare & pull request button, fill and submit the pull request.


## Projects rules

A few rules:
1. As these are mandatory access control policies only what it explicitly required
   should be authorized. Meaning, you should not allow everything (or a large area)
   and blacklist some sub area.
2. A profile **should not break a normal usage of the confined software**. It can
   be complex as simply running the program for your own use case is not alway
   exhaustive of the program features and required permissions.


## Add a profile

1. To add a new profile `foo`, add the file `foo` in `apparmor.d/profile-a-f`. 
   If your profile is part of a large group of profiles, it can also go in
   `apparmor.d/groups`.

2. Write the profile content, the rules depend of the confined program,
   Here is the bare minimum for the program `foo`:
```
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2022 You <your@email>
# SPDX-License-Identifier: GPL-2.0-only

abi <abi/3.0>,

include <tunables/global>

@{exec_path} = /{usr/,}bin/foo
profile foo @{exec_path} {
  include <abstractions/base>

  @{exec_path} mr,

  include if exists <local/foo>
}
```

3. You can automatically set the complain flag on your profile by editing the file `dists/flags/main.flags` and adding a new line with: `foo complain`

4. Build & install for your distribution.


## Profile Guidelines

**A common structure**

AppArmor profiles can be written without any specific guidelines. However, when 
you work with over 1200 profiles, you need a common structure among all the profiles. 

The logic behind it is that if a rule is present in a profile, it should only be
in one place, making profile review easier. 

For example, if a program needs to run executables binary. The rules allowing it
can only be in a specific rule block (just after the `@{exec_path} mr,` rule). It
is therefore easy to ensure some profile features such as: 
* A profile has access to a given resource 
* A profile enforces a strict [write xor execute] (W^X) policy. 

It also improves compatibilities and makes personalization easier thanks to the use of more variables 
 
**Guidelines**

> **Note**: This profile guideline is still evolving, feel free to propose improvment
> as long as it does not vary too much from the existing rules.

In order to ensure a common structure across the profiles, all new profile should
try to follow the guideline presented here.

The rules in the profile should be sorted in rule *block* as follow:
- include
- set rlimit
- capability
- network
- mount
- remount
- umount
- pivot_root
- change_profile
- signal
- ptrace
- unix
- dbus
- file
- Local include

This rule order is taken from AppArmor with minor changes as we tend to:
- Divide the file block in multiple sub categories
- Put the block with the longer rules (files, dbus) after the other blocks

**The file block**

Try to sort the file rules as follow:
- `@{exec_path} mr`, the entry point of the profile
- The binaries and library required: `/{usr/,}bin/`, `/{usr/,}lib/`, `/opt/`...
  It is the only place where you can have `mr`, `rix`, `rPx`, `rUx`, `rPUX` rules.
- The shared resources: `/usr/share`...
- The system configuration: `/etc`...
- The system data: `/var`...
- The user data: `owner @{HOME}/`...
- The user configuration, cache and in general all dotfiles
- Temporary and runtime data: `/tmp/`, `@{run}/`, `/dev/shm/`...
- Sys files: `@{sys}/`...
- Proc files: `@{PROC}/`... 
- Dev files: `/dev/`...
- Deny rules: `deny`...

**The dbus block**

Try to sort the dbus rules as follow:
- The system bus should be sorted *before* the session bus
- The bind rules should be sorred *after* the send & receive rules

For DBus, try to determine peer's label when possible. E.g.:
```
dbus send bus=session path=/org/freedesktop/DBus
     interface=org.freedesktop.DBus
     member={RequestName,ReleaseName}
     peer=(name=org.freedesktop.DBus, label=dbus-daemon),
```
If there is no predictable label it can be omited.

**Other rules**
* Do not use: `/usr/lib` or `/usr/bin` but `/{usr/,}bin/` or `/{usr/,}lib/`.
* Do not use: `/usr/sbin` or `/sbin` but `/{usr/,}{s,}bin/`.
* Always use the apparmor variables.
* In a rule block, the rule shall be alphabetically sorted.
* Subprofile should comes at the end of a profile.
* When some file access share similar purpose, they may be sorted together. Eg:
  ```
  /etc/machine-id r,
  /var/lib/dbus/machine-id r,
  ```

The included tool `aa-log` can be useful to explore the apparmor log 

## Abstractions

This project and the apparmor profile official project provide a large selection
of abstractions to be included in profiles. They should be used.

For instance, instead of writting:
```sh
owner @{HOME}/@{XDG_DOWNLOAD_DIR}/{,**} rw,
```
to allow download directory access, you should write

```sh
include <abstractions/user-download-strict>
```

## AppArmor variables

**Included variables:**

* `@{PROC}=/proc/`
* `@{run}=/run/ /var/run/`
* `@{sys}=/sys/`
* The home root: `@{HOMEDIRS}=/home/`
* The home directories: `@{HOME}=@{HOMEDIRS}/*/ /root/`
* Process id(s): `@{pid}`, `@{pids}`
* User id: `@{uid}`
* Thread id: `@{tid}`
* Classic XDG user directories: 
  - Desktop: `@{XDG_DESKTOP_DIR}="Desktop"`
  - Download: `@{XDG_DOWNLOAD_DIR}="Downloads"`
  - Templates: `@{XDG_TEMPLATES_DIR}="Templates"`
  - Public: `@{XDG_PUBLICSHARE_DIR}="Public"`
  - Documents: `@{XDG_DOCUMENTS_DIR}="Documents"`
  - Music: `@{XDG_MUSIC_DIR}="Music"`
  - Pictures: `@{XDG_PICTURES_DIR}="Pictures"`
  - Videos: `@{XDG_VIDEOS_DIR}="Videos"`

**Additional variables available with this project:**

* Libexec:
  - On Archlinux: `@{libexec}=/{usr/,}lib`
  - On Debian/Ubuntu: `@{libexec}=/{usr/,}libexec`
* Mountpoints root: `@{MOUNTDIRS}=/media/ @{run}/media/ /mnt/`
* Common mountpoints: `@{MOUNTS}=@{MOUNTDIRS}/*/`
* Universally unique identifier: `@{uuid}=[0-9a-fA-F]*-[0-9a-fA-F]*-[0-9a-fA-F]*-[0-9a-fA-F]*-[0-9a-fA-F]*`
* Hexadecimal: `@{hex}=[0-9a-fA-F]*`
* Extended XDG user directories: 
  - Books: `@{XDG_BOOKS_DIR}="Books"`
  - Projects: `@{XDG_PROJECTS_DIR}="Projects"`
  - Screenshots: `@{XDG_SCREENSHOTS_DIR}="@{XDG_PICTURES_DIR}/Screenshots"`
  - Sync: `@{XDG_SYNC_DIR}="Sync"`
  - Torrents: `@{XDG_TORRENTS_DIR}="Torrents"`
  - Vm: `@{XDG_VM_DIR}=".vm"`
  - Wallpapers: `@{XDG_WALLPAPERS_DIR}="@{XDG_PICTURES_DIR}/Wallpapers"`
* Extended XDG dotfiles:
  - SSH: `@{XDG_SSH_DIR}=".ssh"`
  - GPG: `@{XDG_GPG_DIR}=".gnupg"`
  - Cache:` @{XDG_CACHE_HOME}=".cache"`
  - Config: `@{XDG_CONFIG_HOME}=".config"`
  - Data: `@{XDG_DATA_HOME}=".local/share"`
  - Bin: `@{XDG_BIN_HOME}=".local/bin"`
  - Lib: `@{XDG_LIB_HOME}=".local/lib"`
* Full path of the user configuration directories
  - Cache: `@{user_cache_dirs}=@{HOME}/@{XDG_CACHE_HOME}`
  - Config: `@{user_config_dirs}=@{HOME}/@{XDG_CONFIG_HOME}`
  - Bin: `@{user_bin_dirs}=@{HOME}/@{XDG_BIN_HOME}`
  - Lib: `@{user_lib_dirs}=@{HOME}/@{XDG_LIB_HOME}`
* Full path user directories
  - Books: `@{user_books_dirs}=@{HOME}/@{XDG_BOOKS_DIR} @{MOUNTS}/@{XDG_BOOKS_DIR}`
  - Documents: `@{user_documents_dirs}=@{HOME}/@{XDG_DOCUMENTS_DIR} @{MOUNTS}/@{XDG_DOCUMENTS_DIR}`
  - Download: `@{user_download_dirs}=@{HOME}/@{XDG_DOWNLOAD_DIR} @{MOUNTS}/@{XDG_DOWNLOAD_DIR}`
  - Music: `@{user_music_dirs}=@{HOME}/@{XDG_MUSIC_DIR} @{MOUNTS}/@{XDG_MUSIC_DIR}`
  - Pictures: `@{user_pictures_dirs}=@{HOME}/@{XDG_PICTURES_DIR} @{MOUNTS}/@{XDG_PICTURES_DIR}`
  - Projects: `@{user_projects_dirs}=@{HOME}/@{XDG_PROJECTS_DIR} @{MOUNTS}/@{XDG_PROJECTS_DIR}`
  - Public: `@{user_publicshare_dirs}=@{HOME}/@{XDG_PUBLICSHARE_DIR} @{MOUNTS}/@{XDG_PUBLICSHARE_DIR}`
  - Sync: `@{user_sync_dirs}=@{HOME}/@{XDG_SYNC_DIR} @{MOUNTS}/*/@{XDG_SYNC_DIR}`
  - Templates: `@{user_templates_dirs}=@{HOME}/@{XDG_TEMPLATES_DIR} @{MOUNTS}/@{XDG_TEMPLATES_DIR}`
  - Torrents: `@{user_torrents_dirs}=@{HOME}/@{XDG_TORRENTS_DIR} @{MOUNTS}/@{XDG_TORRENTS_DIR}`
  - Videos: `@{user_videos_dirs}=@{HOME}/@{XDG_VIDEOS_DIR} @{MOUNTS}/@{XDG_VIDEOS_DIR}`
  - Vm: `@{user_vm_dirs}=@{HOME}/@{XDG_VM_DIR} @{MOUNTS}/@{XDG_VM_DIR}`

## Additional documentation

* https://gitlab.com/apparmor/apparmor/-/wikis/AppArmor_Core_Policy_Reference
* https://man.archlinux.org/man/apparmor.d.5
* https://presentations.nordisch.org/apparmor/#/

[git]: https://help.github.com/articles/set-up-git/
[write xor execute]: https://en.wikipedia.org/wiki/W%5EX
