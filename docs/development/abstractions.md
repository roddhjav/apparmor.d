---
title: Abstractions
---

This project and the apparmor profile official project provide a large selection of abstractions to be included in profiles. They should always be used as they target wide compatibility across hardware and distribution wile only allowing the bare minimum access.

!!! example

    For instance, to allow download directory access, instead of writing:
    ```sh
    owner @{HOME}/@{XDG_DOWNLOAD_DIR}/{,**} rw,
    ```

    You should write:
    ```sh
    include <abstractions/user-download-strict>
    ```


All of these abstractions can be extended by a system admin by adding rules in a file under `/etc/apparmor.d/<name>.d` where `<name>` is the name of one of these abstractions.


## Application helper

### **`bwrap`**

Minimal set of rules for sandboxed program using `bwrap`. A profile using this abstraction still needs to set:

- The flag: `attach_disconnected`
- Bwrap execution: `@{bin}/bwrap rix,`

### **`bwrap-app`**

Common rules for unknown userland UI applications sandboxed using `bwrap`.

!!! warning

    This abstraction is wide on purpose. It is meant to be used by sandboxed applications that have no way to restrict access depending on the application being confined.

    **Do not use it for classic profile.**


### **`chromium`**

Full set of rules for all chromium based browsers. It works as a *function* and requires some variables to be provided as *arguments* and set in the header of the calling profile:

!!! note ""

    [apparmor.d/groups/browsers/chromium](https://github.com/roddhjav/apparmor.d/blob/e979fe05b06f525e5a65c767b4eabe5600147355/apparmor.d/groups/browsers/chromium#L10-L14)
    ``` sh linenums="10"
    @{name} = chromium
    @{domain} = org.chromium.Chromium
    @{lib_dirs} = @{lib}/@{name}
    @{config_dirs} = @{user_config_dirs}/@{name}
    @{cache_dirs} = @{user_cache_dirs}/@{name}
    ```

If your application requires chromium to run (like electron) use [`chromium-common`](#chromium-common) instead.

### **`chromium-common`**

Minimal set of rules for chromium based application such as electron. Handle access for internal sandbox.

### **`sudo`**

Minimal set of rules for profile including internal `sudo`. Interactive sudo need more rules. It is intended to be used in profile or sub profile that need to elevate their privileges using `sudo` or `su` for a very specific action:
```sh
  @{bin}/sudo rCx  -> root,

  profile root {
    include <abstractions/base>
    include <abstractions/sudo>

    @{bin}/sudo rm,
  
    include if exists <local/<profile_name>_root>
  }
```

### **`systemctl`**

Alternative solution for [child-systemctl](structure.md#children-profiles), when the child profile provide too much/not enough access. This abstraction should be used by a sub profile as follows:
```sh
  @{bin}/systemctl  rCx ->  systemctl,

  profile systemctl {
    include <abstractions/base>
    include <abstractions/systemctl>

    include if exists <local/<profile_name>_systemctl>
  }
```


## Audio

### **`audio-client`**

Most programs do not need access to audio devices, `audio-client` only includes configuration files to be used by client applications.

### **`audio-server`**

Provide access to audio devices. It should only be used by audio servers that need direct access to them. 


## Dbus

See the [Dbus](dbus.md#abstractions) page.


## User files

### **`user-read`**

This abstraction gives read access on all defined user directories. It should only be used if access to **ALL** folders under [xdg directories](../variables.md#xdg-directories) is required.


### **`user-download-strict`**

Provide write access to all user download directories


### **`deny-sensitive-home`**

Deny access to some sensitive directories under `/home/`. It is intended to be used by the few profiles that legitimately require full unrestricted access over all user directories (file browser and search engines). It allows to us to block access to really sensitive data to such profiles.

!!! danger

    **Do not use this abstraction for other profile without explicit authorisation from the project maintainer**

    Per the **[Rule :material-numeric-1-circle:](index.html#rule-mandatory-access-control)** of this project:

    > Only what is explicitly required should be authorized. Meaning, you should **not** allow everything (or a large area) and deny some sub areas.


### **`dconf-write`**

Permissions for querying dconf settings with write access.


## Shell

!!! warning

    This abstractions are only required when an interactive shell is started. Classic shell scripts do not need them.

### **`bash-strict`**

Common rules for interactive shell using bash.

### **`zsh`**

Common rules for interactive shell using zsh.


## System



### **`nameservice-strict`**

Many programs wish to perform nameservice like operations, such as looking up users by name or Id, groups by name or Id, hosts by name or IP, etc.

Use this abstraction instead of upstream `abstractions/nameservice` as upstream abstraction also provide full network access which is not needed for a lot of programs.

### **`systemd-common`**

Common set of rules for internal systemd suite.

!!! warning

    It should **only** be used by the systemd suite.


### **`app-open`**

Instead of allowing the run of all software under `@{bin}` or `@{lib}` the purpose of this abstraction is to list all GUI program that can open resources. Ultimately, only sandbox manager program such as `bwrap`, `snap`, `flatpak`, `firejail` should be present here. Until this day, this profile will be a controlled mess.


## Devices

### **`devices-usb`**

Provide access to USB devices

### **`disks-write`**

Provide read write access to disks devices

### **`disks-read`**

Provide read only access to disks devices


## Desktop Environment

### **`desktop`**

Unified minimal abstraction for all UI application regardless of the desktop environment. When supported in apparmor, condition will be used in this abstraction to filter resources specific for supported DE.

It is safe to use it in GUI application. As well as minimal desktop resource files, it includes access to configuration for: `fonts`, `gtk` & `qt`, `wayland` & `xorg`.

### **`gnome-strict`**

Same than `abstractions/desktop` but limited to gnome.

### **`kde-strict`**

Same than `abstractions/desktop` but limited to KDE.

## Graphics

Use either [`graphics`](#graphics) or [`graphics-full`](#graphics-full). The other abstractions are hardware/software dependant and should not usually be used directly.

### **`graphics`**

Unified abstraction for GPU access regardless of the hardware used.

Replace and highly restrict `<abstractions/opencl>`

### **`graphics-full`**

Identical to [`graphics`](#graphics) with more direct access to nvidia GPU devices.

### **`dri`**

Linux's graphics stack which allows unprivileged user-space programs to issue commands to graphics hardware without conflicting with other programs. Mostly used by Intel (integrated or not) and AMD GPU.

Modernized equivalent of both `dri-common` and `dri-enumerate`

### **`nvidia-strict`**

Modernized equivalent of `abstractions/nvidia`

### **`vulkan-strict`**

Modernized equivalent of `abstractions/vulkan`

