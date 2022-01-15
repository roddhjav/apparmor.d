# Contributing

You want to contribute to `apparmor.d`, **thank a lot for this.** You will find
in this page all the useful information needed to contribute.


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

In order to ensure a common structure across the profiles, all new profile should try to follow the guideline presented here.

The rules in the profile should be sorted as follow: 
- include
- capability
- ptrace
- signal
- network 
- mount
- @{exec_path} mr,
- The binaries and library required: `/{usr/,}bin/`, `/{usr/,}lib/`, `/opt/`...
- The shared resources: `/usr/share`...
- The system configuration: `/etc`...
- The user data: `owner @{HOME}/`...
- The user configuration (all dotfiles)
- Temporary data: `/tmp/`, `@{run}/`...
- Sys files: `@{sys}/`...
- Proc files: `@{PROC}/`... 
- Dev files: `/dev/`...


**Other rules**
* Do not use: `/usr/lib` or `/usr/bin` but `/{usr/,}bin/` or `/{usr/,}lib/`.
* Always use the apparmor variables.
* In a rule block, the rule shall be alphabetically sorted.
* When some file access share similar purpose, they shall be sorted together. Eg:
    ```
    /etc/machine-id r,
    /var/lib/dbus/machine-id r,
    ```

## AppArmor variables

**Included variables:**

* `@{PROC}=/proc/`
* `@{run}=/run/ /var/run/`
* `@{sys}=/sys/`
* The Home directory: `@{HOME}`
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

* Common mountpoints: `@{MOUNTS}=/media/ @{run}/media /mnt`
* Extended XDG user directories: 
    - Projects: `@{XDG_PROJECTS_DIR}="Projects"`
    - Books: `@{XDG_BOOKS_DIR}="Books"`
    - Wallpapers: `@{XDG_WALLPAPERS_DIR}="Pictures/Wallpapers"`
    - Sync: `@{XDG_SYNC_DIR}="Sync"`
    - Vm: `@{XDG_VM_DIR}=".vm"`
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
* Other full path user directories
    - Sync: `@{user_sync_dirs}=@{HOME}/@{XDG_SYNC_DIR} @{MOUNTS}/*/@{XDG_SYNC_DIR}`

## Additional documentation

* https://gitlab.com/apparmor/apparmor/-/wikis/AppArmor_Core_Policy_Reference
* https://presentations.nordisch.org/apparmor/#/

[git]: https://help.github.com/articles/set-up-git/
