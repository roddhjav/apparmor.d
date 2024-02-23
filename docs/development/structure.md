---
title: Structure
---

Description of common structure found across various AppArmor profiles


## Programs to not confine

Some programs should not be confined by themselves. For example, tools such as `ls`, `rm`, `diff` or `cat` do not have profiles in this project. Let's see why.

These are general tools that in a general context can legitimately access any file in the system. Therefore, the confinement of such tools by a global profile would at best be minimal at worst be a security theater.

It gets even worse. Let's say, we write a profile for `cat`. Such a profile would need access to `/etc/`. We will add the following rule:
```sh
  /etc/{,**} rw,
```

However, as `/etc` can contain sensitive files, we now want to explicitly prevent access to these sensitive files. Problems:

1. How do we know the exhaustive list of *sensitive files* in `/etc`?
2. How do we ensure access to these sensitive files are not required?
3. This breaks the principle of mandatory access control.
   See the [first rule of this project](index.md#project-rules) that is to only allow
   what is required. Here we allow everything and blacklist some paths.

It creates even more issues when we want to use this profile in other profiles. Let's take the example of `diff`. Using this rule: `@{bin}/diff rPx,` will restrict access to the very generic and not very confined `diff` profile. Whereas most of the time, we want to restrict `diff` to some specific file in our profile:

* In `dpkg`, an internal child profile (`rCx -> diff`), allows `diff` to only
  access etc config files:

!!! note ""

    [apparmor.d/apparmor.d/groups/apt/dpkg](https://github.com/roddhjav/apparmor.d/blob/accf5538bdfc1598f1cc1588a7118252884df50c/apparmor.d/groups/apt/dpkg#L123)
    ``` aa linenums="123"
    profile diff {
      include <abstractions/base>
      include <abstractions/consoles>

      @{bin}/       r,
      @{bin}/pager mr,
      @{bin}/less  mr,
      @{bin}/more  mr,
      @{bin}/diff  mr,

      owner @{HOME}/.lesshs* rw,

      # Diff changed config files
      /etc/** r,

      # For shell pwd
      /root/ r,

    }
    ```

* In `pass`, as it is a dependency of pass. Here `diff` inherits pass' profile 
  and has the same access than the pass profile, so it will be allowed to diff
  password files because more than a generic `diff` it is a `diff` for the pass
  password manager:

!!! note ""

    [apparmor.d/apparmor.d/profiles-m-r/pass](https://github.com/roddhjav/apparmor.d/blob/accf5538bdfc1598f1cc1588a7118252884df50c/apparmor.d/profiles-m-r/pass#L20
    )
    ``` aa linenums="20"
      @{bin}/diff      rix,
    ```

**What if I still want to protect these programs?**

You do not protect these programs. *Protect the usage you have of these programs*.
In practice, it means that you should put your development's terminal in a
sandbox managed with [Toolbox].

!!! example "To sum up"

    1. Do not a create profile for programs such as: `rm`, `ls`, `diff`, `cd`, `cat`
    2. Do not a create profile for the shell: `bash`, `sh`, `dash`, `zsh`
    3. Use [Toolbox].

[Toolbox]: https://containertoolbx.org/



## Abstractions

This project and the apparmor profile official project provide a large selection of abstractions to be included in profiles. They should be used.

For instance, to allow download directory access, instead of writing:
```sh
owner @{HOME}/@{XDG_DOWNLOAD_DIR}/{,**} rw,
```

You should write:
```sh
include <abstractions/user-download-strict>
```


## Children profiles

Usually, a child profile is in the [`children`][children] group. They have the following note:

!!! quote

    Note: This profile does not specify an attachment path because it is
    intended to be used only via `"Px -> child-open"` exec transitions
    from other profiles. 

[children]: https://github.com/roddhjav/apparmor.d/blob/main/apparmor.d/groups/children

Here is an overview of the current children profile:

1. **`child-open`**: To open resources. Instead of allowing the run of all
   software in `@{bin}/`, the purpose of this profile is to list all GUI
   programs that can open resources. Ultimately, only sandbox manager programs
   such as `bwrap`, `snap`, `flatpak`, `firejail` should be present here. Until
   this day, this profile will be a controlled mess.

2. **`child-pager`**: Simple access to pager such as `pager`, `less` and `more`.
   This profile supposes the pager is reading its data from stdin, not from a
   file on disk.

3. **`child-systemctl`**: Common `systemctl` action. Do not use it too much as most
   of the time you will need more privilege than what this profile is giving you.


## Browsers

Chromium based browsers share a similar structure. Therefore, they share the same abstraction: [`abstractions/chromium`][chromium] that includes most of the profile content.

This abstraction requires the following variables defined in the profile header:
```sh
@{name} = chromium
@{domain} = org.chromium.Chromium
@{lib_dirs} = @{lib}/chromium
@{config_dirs} = @{user_config_dirs}/chromium
@{cache_dirs} = @{user_cache_dirs}/chromium
```

If your application requires chromium to run (like electron) use [`abstractions/chromium-common`][chromium-common] instead.

[chromium]: https://github.com/roddhjav/apparmor.d/blob/main/apparmor.d/abstractions/chromium
[chromium-common]: https://github.com/roddhjav/apparmor.d/blob/main/apparmor.d/abstractions/chromium-common

## Udev rules

See the **[kernel docs][kernel]** to check the major block and char numbers used in `/run/udev/data/`.

Special care must be given as sometimes udev numbers are allocated dynamically by the kernel. Therefore, the full range must be allowed:

!!! note ""

    [apparmor.d/groups/virt/libvirtd](https://github.com/roddhjav/apparmor.d/blob/b2af7a631a2b8aca7d6bdc8f7ff4fdd5ec94220e/apparmor.d/groups/virt/libvirtd#L188)
    ``` aa linenums="179"
      @{run}/udev/data/c@{dynamic}:@{int} r,  # For dynamic assignment range 234 to 254, 384 to 511
    ```

[kernel]: https://raw.githubusercontent.com/torvalds/linux/master/Documentation/admin-guide/devices.txt


## No New Privileges

[**No New Privileges**](https://www.kernel.org/doc/html/latest/userspace-api/no_new_privs.html) is a flag preventing a newly-started program to get more privileges that its parent. So it is a **good thing** for security. And it is commonly used in systemd unit files (when possible). This flag also prevents transition to other profile because it could be less restrictive than the parent profile (no `Px` or `Ux` allowed).

The possible solutions are:

* The easiest (and unfortunately less secure) workaround is to ensure the programs do not run with no new privileges flag by disabling `NoNewPrivileges` in the systemd unit (or any other [options implying it](https://man.archlinux.org/man/core/systemd/systemd.exec.5.en#SECURITY)).
* Inherit the current confinement (`ix`)
* [Stacking](https://gitlab.com/apparmor/apparmor/-/wikis/AppArmorStacking)
