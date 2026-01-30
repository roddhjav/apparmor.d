---
title: Internal
---

## Profile Context

These are context helper to be used for in sub profile, they aim at providing a minimal set of rules for a given program. The calling profile only needs to add rules dependant of its use case.

See [abstractions/app](abstractions.md#context-helper) for more information.


## Open Resources

The standard way to allow opening resources such as URL, pictures, video, in this project is to use one of the `child-open` profile available in the [`children`](https://github.com/roddhjav/apparmor.d/blob/main/apparmor.d/groups/children) group.

**Example:**
```sh
@{open_path} rPx -> child-open,
```


### Manual

Directly using any of the following:

- `@{bin}/* PUx,`
- `include <abstractions/app-launcher-user>`
- `include <abstractions/app-launcher-root>`

Allow every installed program to be started from the current program with or without profile. This is a very permissive rule and should be avoided if possible. They are however legitimately needed for program launcher.

### **`child-open`**

Instead of allowing the ability to run all software in `@{bin}/`, the purpose of this profile is to list all GUI programs that can open resources. Ultimately, only sandbox manager programs such as `bwrap`, `snap`, `flatpak`, `firejail` should be present here. Until this day, this profile will be a controlled mess.

??? quote "[children/child-open](https://github.com/roddhjav/apparmor.d/blob/main/apparmor.d/groups/children/child-open)"

    ``` aa
      # Sandbox managers
      @{bin}/bwrap                  rPUx,
      @{bin}/firejail               rPUx,
      @{bin}/flatpak                rPx,
      @{bin}/snap                   rPx,

      # Labelled programs
      @{archive_viewers_path}       rPUx,
      @{browsers_path}              rPx,
      @{document_viewers_path}      rPUx,
      @{emails_path}                rPUx,
      @{file_explorers_path}        rPx,
      @{help_path}                  rPx,
      @{image_viewers_path}         rPUx,
      @{offices_path}               rPUx,
      @{text_editors_path}          rPUx,

      # Others
      @{bin}/blueman-tray           rPx,
      @{bin}/discord{,-ptb}         rPx,
      @{bin}/draw.io                rPUx,
      @{bin}/dropbox                rPx,
      @{bin}/element-desktop        rPx,
      @{bin}/extension-manager      rPx,
      @{bin}/filezilla              rPx,
      @{bin}/flameshot              rPx,
      @{bin}/gimp*                  rPUx,
      @{bin}/gnome-calculator       rPUx,
      @{bin}/gnome-disk-image-mounter rPx,
      @{bin}/gnome-disks            rPx,
      @{bin}/gnome-software         rPx,
      @{bin}/gwenview               rPUx,
      @{bin}/kgx                    rPx,
      @{bin}/qbittorrent            rPx,
      @{bin}/qpdfview               rPx,
      @{bin}/smplayer               rPx,
      @{bin}/steam-runtime          rPUx,
      @{bin}/telegram-desktop       rPx,
      @{bin}/transmission-gtk       rPx,
      @{bin}/viewnior               rPUx,
      @{bin}/vlc                    rPUx,
      @{bin}/xbrlapi 	              rPx,

      # Backup
       @{lib}/deja-dup/deja-dup-monitor rPx,
    ```

### **`child-open-browsers`**

 This version of child-open only allow to open browsers.

??? quote "[children/child-open-browsers](https://github.com/roddhjav/apparmor.d/blob/main/apparmor.d/groups/children/child-open-browsers)"

    ``` aa
      @{browsers_path}              rPx,
    ```

### **`child-open-help`**

This version of child-open only allow to open browsers and help programs.

??? quote "[children/child-open-help](https://github.com/roddhjav/apparmor.d/blob/main/apparmor.d/groups/children/child-open-help)"

    ``` aa
      @{browsers_path}              rPx,
      @{help_path}                  rPx,
    ```

### **`child-open-strict`**

This version of child-open only allow to open browsers & folders:

??? quote "[children/child-open-strict](https://github.com/roddhjav/apparmor.d/blob/main/apparmor.d/groups/children/child-open-strict)"

    ``` aa
      @{browsers_path}               Px,
      @{file_explorers_path}         Px,
    ```


!!! warning

    Although needed to not break a program, wrongly used these profiles can lead to confinment escape.


## Children profiles

Usually, a child profile is in the [`children`](https://github.com/roddhjav/apparmor.d/blob/main/apparmor.d/groups/children) group. They have the following note:

!!! quote

    Note: This profile does not specify an attachment path because it is intended to be used only via `"Px -> child-open"` exec transitions from other profiles. 

<!-- ### **`child-dpkg`**

### **`child-dpkg-divert`** -->

### **`child-modprove-nvidia`**

Used internally by the `nvidia` abstraction.

### **`child-pager`**

Simple access to pagers such as `pager`, `less` and `more`. This profile assumes the pager is reading its data from stdin, not from a file on disk. Supported pagers are: `sensible-pager`, `pager`, `less`, and `more`.
It can be as follows in a profile:
```
  @{pager_path} rPx -> child-pager,
```

### **`child-systemctl`**

Common `systemctl` action. Do not use it too much as most of the time you will need more privilege than what this profile is giving you.

It is recommended to transition [in a subprofile](abstractions.md#appsystemctl) everything that is not generic and that may require some access (so restart, enable...), while `child-systemctl` can handle the more basic tasks.


## Labelled programs

All common programs are tracked and labelled in the [`apparmor.d/tunables/multiarch.d/programs`](https://github.com/roddhjav/apparmor.d/blob/main/apparmor.d/tunables/multiarch.d/programs) and 
[`apparmor.d/tunables/multiarch.d/paths`](https://github.com/roddhjav/apparmor.d/blob/main/apparmor.d/tunables/multiarch.d/paths) files. They can be used in a `child-open` profile or directly in a profile. They are useful to allow opening resources using a kind of program (browsers, image viewer, text editor...), instead of allowing a given program path.

## Re-attached path

**[<span class="pg-red">:material-tag-heart-outline: abi/4.0</span>]("Minimum version")**

The flag `attach_disconnect` control how disconnected paths are handled. It determines if pathnames resolved to be outside the namespace are attached to the root (ie. have the `/` character prepended). 
It is a security issue as it allows disconnected paths to alias to other files that exist in the file name. Therefore, it is only provided to work around problems that can arise with sandboxed programs.

AppAmor 4.0 provides the `attach_disconnect.path` flag allowing to reattach this path to a prefix that is not `/`. When used it provides an important security improvement from AppArmor 3.0.

**`apparmor.d`** uses `attach_disconnect.path` by **default and automatically** on all profiles with the `attach_disconnect` flag. The attached path is set to `@{att}` a new dynamically generated variable set at build time in the preamble of all profile to be:

- `@{att}=/att/<profile_name>` for profile with `attach_disconnect` flag.
- `@{att}=/` for other profiles


## User Confinement

[:material-police-badge-outline:{ .pg-red }](../full-system-policy.md "Full System Policy only (FSP)")

!!! warning "TODO"


## No New Privileges

[**No New Privileges**](https://www.kernel.org/doc/html/latest/userspace-api/no_new_privs.html) is a flag preventing a newly started program to get more privileges than its parent process. This is a **good thing** for security. And it is commonly used in systemd unit files (when possible). This flag also prevents transitions to other profiles because it could be less restrictive than the parent profile (no `Px` or `Ux` allowed).

The possible solutions are:

* The easiest (and unfortunately less secure) workaround is to ensure the programs do not run with no new privileges flag by disabling `NoNewPrivileges` in the systemd unit (or any other [options implying it](https://man.archlinux.org/man/core/systemd/systemd.exec.5.en#SECURITY)).
* Inherit the current confinement (`ix`)
* [Stacking](#stacking)

## Stacking

[Stacking](https://gitlab.com/apparmor/apparmor/-/wikis/AppArmorStacking) of two or more profile is the strict intersection them. It is a way to ensure that a profile never becomes more permissive than the intersection of all profiles in the stack. It provides several abilities to the policy author: 

- It can be used to ensure that confinement never becomes more permissive.
- To reduce the permissions of a generic profile on a specific task.
- To provide both system level and container and user level policy (when combined with policy namespaces).

!!! note ""

    [apparmor.d/groups/browsers/chromium](https://github.com/roddhjav/apparmor.d/blob/b51576139b3ed3125aaa3ea4d737a77baac0f00e/apparmor.d/groups/browsers/chromium#L25)
    ``` aa linenums="23"
    profile chromium @{exec_path} {
      ...
      @{lib_dirs}/chrome_crashpad_handler  rPx -> chromium//&chromium-crashpad-handler,
      ...
    }
    ```

## Udev rules
c
See the **[kernel docs](https://www.kernel.org/doc/html/latest/admin-guide/devices.html)** to check the major block and char numbers used in `/run/udev/data/`.

Special care must be given as sometimes udev numbers are allocated dynamically by the kernel. Therefore, the full range must be allowed:

!!! note ""

    [apparmor.d/groups/virt/libvirtd](https://github.com/roddhjav/apparmor.d/blob/b2af7a631a2b8aca7d6bdc8f7ff4fdd5ec94220e/apparmor.d/groups/virt/libvirtd#L188)
    ``` aa linenums="179"
      @{run}/udev/data/c@{dynamic}:@{int} r,  # For dynamic assignment range 234 to 254, 384 to 511
    ```
