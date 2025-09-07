---
title: Abstractions
---

This project and the official apparmor-profiles project provide a large selection of abstractions to be included in profiles. They should always be used as they target wide compatibility across hardware and distributions while only allowing the bare minimum access.

!!! example

    For instance, to allow download directory access instead of read and write permissions:
    ```sh
    owner @{HOME}/@{XDG_DOWNLOAD_DIR}/{,**} rw,
    ```

    You should write:
    ```sh
    include <abstractions/user-download-strict>
    ```


All of these abstractions can be extended by a system admin by adding rules in a file under `/etc/apparmor.d/<name>.d` where `<name>` is the name of one of these abstractions.

## Architecture

Abstraction are structured in layers as follows:

- **Layer 0:** for core atomic functionalities. They cannot include other abstractions.

    E.g.: *this resource uses* `mesa`, `openssl`, `bash-strict`, `gtk`...

- **Layer 1:** for generic access. Cannot be architecture or device specific. Needs to be agnostic.

    E.g.: *This program needs/has this resource.* `nameservice`, `authentication`, `base`, `shell`, `graphics`, `audio-client`, `desktop`, `kde`, `gnome`...

- **Layer 2:** for common kind of program. Only present inside `abstraction/common`. Multiple layer 2 can be used alongside with layer 1 and 0 abstractions.

    E.g.: *This program kind is* is a game, an electron app, a gnome app, sandboxed with bwrap app, a systemd app...

- **Layer 3:** for application. Only present inside `abstraction/app`. The use of a layer 3 abstraction usually means you should not use any other abstractions (but base). Not a strict rule, but a good practice. Mostly used to provide common rules for subprofiles where the subprofiles only need to add rules for the specific use case.

    E.g.: *This program is* `firefox`, `sudo`, `systemctl`, `pgrep`, `editor`, `chromium`...


## Application helper

Abstraction that aims at including a complete set of rules for a given program. The calling profile only needs to add rules dependant of its use case/program.

It is mostly useful for program often used in sub profile or for forks based on the same upstream.

### **`app/chromium`**

A full set of rules for all chromium based browsers. It works as a *function* and requires some variables to be provided as *arguments* and to be set in the header of the calling profile:

!!! note ""

    [apparmor.d/groups/browsers/chromium](https://github.com/roddhjav/apparmor.d/blob/e979fe05b06f525e5a65c767b4eabe5600147355/apparmor.d/groups/browsers/chromium#L10-L14)
    ``` sh linenums="10"
    @{name} = chromium
    @{domain} = org.chromium.Chromium
    @{lib_dirs} = @{lib}/@{name}
    @{config_dirs} = @{user_config_dirs}/@{name}
    @{cache_dirs} = @{user_cache_dirs}/@{name}
    ```

If your application requires chromium to run use [`common/chromium`](#commonchromium) or [`common/electron`](#commonelectron)
instead.

### **`app/firefox`**

Similar to `app/chromium` but for Firefox based browsers (and thunderbird). It requires the same *arguments* as `app/chromium`:


## Context helper

These are context helper to be used for in sub profile, they aim at providing a minimal set of rules for a given program. The calling profile only needs to add rules dependant of its use case.

### **`app/editor`**

A minimal set of rules for profiles including terminal editor. It is intended to be used in profiles or sub-profiles that need to edit file using the user editor of choice. The following editors are supported:

- neo vim
- vim
- nano

```sh
  @{editor_path} rCx -> editor,

  profile editor {
    include <abstractions/base>
    include <abstractions/app/editor>

    include if exists <local/<profile_name>_editor>
  }
```

### **`app/kmod`**

A minimal set of rules for profiles that need to load kernel modules. It is intended to be used in profiles or sub-profiles that need to load kernel modules for a very specific action:

```sh
  @{bin}/modprobe rCx -> kmod,

  profile kmod {
    include <abstractions/base>
    include <abstractions/app/kmod>

    include if exists <local/<profile_name>_kmod>
  }
``` 

### **`app/open`**

Set of rules for `child-open-*` profiles. It should usually not be used directly in a profile.

### **`app/pgrep`**

 Minimal set of rules for pgrep/pkill. It is intended to be used in profiles or sub-profiles that need to use `pgrep` or `pkill` for a very specific action:

 ```sh
  @{bin}/pgrep rCx -> pgrep,

  profile pgrep {
    include <abstractions/base>
    include <abstractions/app/pgrep>

    include if exists <local/<profile_name>_pgrep>
  }
 ```

### **`app/sudo`**

A minimal set of rules for profiles including internal `sudo`. Interactive sudo needs more rules. It is intended to be used in profiles or sub-profiles that need to elevate their privileges using `sudo` or `su` for a very specific action:
```sh
  @{bin}/sudo rCx  -> root,

  profile root {
    include <abstractions/base>
    include <abstractions/app/sudo>

    include if exists <local/<profile_name>_root>
  }
```


### **`app/pkexec`**

A minimal set of rules for profiles including internal `pkexec`. Like `app/sudo`, it should be used in profiles or sub-profiles that need to elevate their privileges using `pkexec` for a very specific action:

```sh
  @{bin}/pkexec rCx -> pkexec,

  profile pkexec {
    include <abstractions/base>
    include <abstractions/app/pkexec>

    include if exists <local/<profile_name>_pkexec>
  }
```

### **`app/systemctl`**

An alternative solution for [child-systemctl](internal.md#children-profiles), when the child profile provides too much/not enough access. This abstraction should be used by a sub profile as follows:

```sh
  @{bin}/systemctl  rCx ->  systemctl,

  profile systemctl {
    include <abstractions/base>
    include <abstractions/app/systemctl>

    include if exists <local/<profile_name>_systemctl>
  }
```

### **`app/udevadm`**

A minimal set of rules for profiles including internal `udevadm` as read-only. It is intended to be used in profiles or sub-profiles that need to use `udevadm` for a very specific action:

```sh
  @{bin}/udevadm rCx -> udevadm,

  profile udevadm {
    include <abstractions/base>
    include <abstractions/app/udevadm>

    include if exists <local/<profile_name>_udevadm>
  }
```

## Common Dependencies

On the contrary of [`abstractions/app/`](#application-helper), abstractions in this directory are expected to provide a minimal set of rules to make a program using a dependency work. 

### **`common/app`**

Common rules for unknown userland UI applications that are sandboxed using `bwrap`.

!!! warning

    This abstraction is wide on purpose. It is meant to be used by sandboxed applications that have no way to restrict access depending on the application being confined.

    **Do not use it for classic profile.**


### **`common/apt`**

Minimal access to apt sources, preferences, and configuration.

### **`common/bwrap`**

Minimal set of rules for sandboxed programs using `bwrap`. A profile using this abstraction still needs to set:

- The flag: `attach_disconnected`
- Bwrap execution: `@{bin}/bwrap rix,`


### **`common/chromium`**

A minimal set of rules for chromium based application. Handle access for internal sandbox.

It works as a *function* and requires some variables to be provided as *arguments* and set in the header of the calling profile:

!!! note ""

    [apparmor.d/profile-s-z/spotify](https://github.com/roddhjav/apparmor.d/blob/main/apparmor.d/groups/steam/steam#L24-L25)
    ``` sh linenums="24"
    @{domain} = org.chromium.Chromium
    ```

### **`common/electron`**

A minimal set of rules for all electron based UI applications. It works as a *function* and requires some variables to be provided as *arguments* and set in the header of the calling profile:

!!! note ""

    [apparmor.d/profile-s-z/spotify](https://github.com/roddhjav/apparmor.d/blob/7d1380530aa56f31589ccc6a360a8144f3601731/apparmor.d/profiles-s-z/spotify#L10-L13)
    ``` sh linenums="10"
    @{name} = spotify
    @{domain} = org.chromium.Chromium
    @{lib_dirs} = /opt/@{name}
    @{config_dirs} = @{user_config_dirs}/@{name}
    @{cache_dirs} = @{user_cache_dirs}/@{name}
    ```

### **`common/game`**

Core set of resources for any games on Linux. Runtimes such as sandboxing, wine, proton, game launchers should use this abstraction. 

This abstraction uses the following tunables:

- `@{XDG_GAMESSTUDIO_DIR}` for game studio and game engines specific directories (Default: `@{XDG_GAMESSTUDIO_DIR}="unity3d"`)
- `@{user_games_dirs}` for user specific game directories (e.g.: steam storage dir)

### **`common/systemd`**

Common set of rules for internal systemd suite.

!!! warning

    It should **only** be used by the systemd suite.


## Audio

### **`audio-client`**

Most programs do not need access to audio devices, `audio-client` only includes configuration files to be used by client applications.

### **`audio-server`**

Provides access to audio devices. It should only be used by audio servers that need direct access to them. 


## Dbus

See the [Dbus](dbus.md#abstractions) page.


## User files

### **`user-read`**

This abstraction gives read access on all defined user directories. It should only be used if access to **ALL** folders under [xdg directories](../variables.md#xdg-directories) is required.


### **`user-download-strict`**

Provides write access to all user download directories


### **`deny-sensitive-home`**

Denies access to some sensitive directories under `/home/`. It is intended to be used by the few profiles that legitimately require full unrestricted access over all user directories (file managers and search engines). It allows to us to block access to really sensitive data to such profiles.

!!! danger

    **Do not use this abstraction for other profiles without explicit authorisation from the project maintainer**

    Per the **[Rule :material-numeric-1-circle:](index.md#rule-mandatory-access-control)** of this project:

    > Only what is explicitly required should be authorized. Meaning, you should **not** allow everything (or a large area) and deny some sub areas.


### **`dconf-write`**

Permissions for querying dconf settings with write access.


## Shell

!!! warning

    These abstractions are only required when an interactive shell is started. Classic shell scripts do not need them.


Only use [`shells`](#shells), other abstractions are software dependant and should not usually be used directly.

### **`shells`**

Common rules for interactive shells.

### **`bash-strict`**

Common rules for interactive shell using bash.

### **`zsh`**

Common rules for interactive shell using zsh.

### **`fish`**

Common rules for interactive shell using fish.

## System



### **`nameservice-strict`**

Many programs wish to perform nameservice like operations, such as looking up users by name or ID, groups by name or ID, hosts by name or IP, etc.

Use this abstraction instead of upstream `abstractions/nameservice` as upstream abstraction also provide full network access which is not needed for a lot of programs.

### **`app-open`**

Instead of allowing the run of all software under `@{bin}` or `@{lib}` the purpose of this abstraction is to list all GUI program that can open resources. Ultimately, only sandbox manager program such as `bwrap`, `snap`, `flatpak`, `firejail` should be present here. Until this day, this profile will be a controlled mess.

### **`app-launcher-root`**

### **`app-launcher-user`**


## Devices

### **`devices-usb`**

Provides access to USB devices

### **`disks-write`**

Provides read write access to disks devices

### **`disks-read`**

Provides read-only access to disks devices


## Desktop Environment

### **`desktop`**

Unified minimal abstraction for all UI applications regardless of the desktop environment. When supported in apparmor, the condition will be used in this abstraction to filter resources specific for supported DE.

It is safe to use this in GUI applications as well as minimal desktop resource files, it includes access to configuration for: `fonts`, `gtk` & `qt`, `wayland` & `xorg`.

### **`gnome-strict`**

Same as `abstractions/desktop` but limited to gnome.

### **`kde-strict`**

Same as `abstractions/desktop` but limited to KDE.

## Graphics

Use either [`graphics`](#graphics) or [`graphics-full`](#graphics-full). The other abstractions are hardware/software dependent and should not usually be used directly.

### **`graphics`**

Unified abstraction for GPU access regardless of the hardware used.

Replace and highly restrict `<abstractions/opencl>`

### **`graphics-full`**

Identical to [`graphics`](#graphics) with more direct access to nvidia GPU devices.

### **`dri`**

Linux's graphics stack which allows unprivileged user-space programs to issue commands to graphics hardware without conflicting with other programs. Mostly used by Intel (integrated or not) and AMD GPUs.

Modernized equivalent of both `dri-common` and `dri-enumerate`

### **`nvidia-strict`**

Modernized equivalent of `abstractions/nvidia`

### **`vulkan-strict`**

Modernized equivalent of `abstractions/vulkan`

