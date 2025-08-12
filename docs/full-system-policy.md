---
title: Full system policy (FSP)
---

!!! danger

    Full system policy is still under early development:
    
    - Do not run this outside of a development VM! 
    - This is an **advanced** feature, you should understand what you are doing before use.

    **You have been warned!!!**

!!! quote

    AppArmor is also capable of being used for full system policy where processes are by default not running under the `unconfined` profile. This might be useful for high security environments or embedded systems.

    *Source: [AppArmor Wiki][apparmor-wiki]*


## Overview

The default mode of `apparmor.d` is the more advanced confinement configuration we can achieve while being as simple as installing a package and doing some minor configuration on your system. By design, a full system confinement does not work this way. Before enabling you need to consider your use case and security objective.

Particularly:

- Every system application will be **blocked** if they do not have a profile.
- Any non-standard system app need to be explicitly profiled and allowed to run. For instance, if you want to use your own proxy or VPN software, you need to ensure it is correctly profiled and allowed to run in the `systemd` profile.
- Desktop environment must be explicitly supported, your UI will not start otherwise. Again, it is a **feature**.
- In FSP mode, all sandbox managers **must** have a profile. Then user sandboxed applications (flatpak, snap, etc) will work as expected.
- PID 1 is the last program that should be confined. It does not make sense to confine only PID. All other programs must be confined first.



## Installation


This feature is only enabled when the project is built with `just fsp`. [Early policy](https://gitlab.com/apparmor/apparmor/-/wikis/AppArmorInSystemd#early-policy-loads) load **must** also be enabled. Once `apparmor.d` has been installed in FSP mode, it is required to reboot to apply the changes.

In `/etc/apparmor/parser.conf` ensure you have:
```
write-cache
cache-loc /etc/apparmor/earlypolicy/
Optimize=compress-fast
```

=== ":material-arch: Archlinux"

    In `PKGBUILD`, replace `just complain` by `just fsp-complain`:

    ```diff
    -  just complain
    +  just fsp-complain
    ```

    Then, build the package with: `just pkg`

=== ":material-ubuntu: Ubuntu"

    In `debian/rules`, replace `just complain` by `just fsp-complain`:

    ```make
      override_dh_auto_build:
    -     just complain
      override_dh_auto_build:
    +     just fsp-complain
    ```

    Then, build the package with: `just dpkg`

=== ":material-debian: Debian"
    
    In `debian/rules`, replace `just complain` by `just fsp-complain`:

    ```make
      override_dh_auto_build:
    -     just complain
      override_dh_auto_build:
    +     just fsp-complain
    ```

    Then, build the package with: `just dpkg`

=== ":simple-suse: openSUSE"

    In `dists/apparmor.d.spec`, replace `just complain` by `just fsp-complain`:

    ```diff
       %build
    -  just complain
       %build
    +  just fsp-complain
    ```

    Then, build the package with: `just rpm`

=== ":material-home: Partial Install"

    Use the `just fsp-complain` command to build instead of `just complain`


## Structure

The profiles dedicated for full system policies are maintained in the **[`_full`][full]** group.

### Systemd

**`systemd`**

This profile aims to confine PID 1. Systemd is (kind of obviously) a highly privileged program. The purpose of this profile is to transition to other less privileged program as soon as possible. On high security environments, it can also be used to strictly limit the list of allowed privileged program.

- It allows internal systemd access,
- It allows starting all common root services.

To work as intended, all privileged services started by systemd **must** have a profile. For a given distribution, the list of these services can be found under:
```sh
/usr/lib/systemd/system-generators/*
/usr/lib/systemd/system-environment-generators/*
/usr/lib/systemd/system/*.service
```

The main [fallback](#fallback) profile (`default`) is not intended to be used by privileged program or service. Such programs must have a dedicated profile and will fail otherwise. This is a **feature**, not a bug.

**`systemd-user`**

This profile is for `systemd --user`, it aims to confine userland systemd. It does not require a lot of access and is only intended to handle user services.

- It allows internal systemd user access,
- It allows starting all common user services.

To work as intended, userland services started by `systemd --user` **should** have a profile. For a given distribution, the list of these services can be found under:

```sh
/usr/lib/systemd/user-environment-generators/*
/usr/lib/systemd/user-generators/*
/usr/lib/systemd/user/*.service
```

!!! info

    To be allowed to run, additional root or user services may need to add extra rules inside the `usr/systemd.d` or `usr/systemd-user.d` directory. For example, when installing a new privileged service `foo` with [stacking](development/internal.md#no-new-privileges) you may need to add the following to `/etc/apparmor.d/usr/systemd.d/foo`:
    ```
    @{lib}/foo rPx -> systemd//&foo,
    ```

### Role Based Access Control (RBAC)

In FSP, interactive shell from the user must be confined. This is done through [pam_apparmor](https://gitlab.com/apparmor/apparmor/-/wikis/pam_apparmor). It provides [Role-based access controls (RBAC)](https://en.wikipedia.org/wiki/Role-based_access_control) that can restrict interactive shell to well-defined role. The role needs to be defined. This project ship with a default set of roles, but you can create your own. The default roles are:

- **`user`**: This is the default role. It is used for any user that does not have a specific role defined. It has access to the user home directory and other sensitive files.

- **`admin`**: This role is used for any user that has administrative access. It has access to the system files and directories, but not to the user home directory.

- **`system`**: This role is used for any user that has system access. It has access to the system files and directories, but not to the user home directory.

### Fallback

In addition to the `systemd` profiles, a full system policy needs to ensure that no programs run in an unconfined state at any time. The fallback profiles consist of a set generic specialized profiles:

- **`default`** is used for any *classic* user application with a GUI. It has full access to user home directories.
- **`bwrap`, `bwrap-app`** are used for *classic* user application that are sandboxed with **bwrap**.

!!! warning

    The main fallback profile (`default`) is not intended to be used by privileged program or service. Such programs **must** have they dedicated profile and would break otherwise.

Additionally, special user access can be setup using PAM rules set such as a random shell interactively opened (as user or as root).

[apparmor-wiki]: https://gitlab.com/apparmor/apparmor/-/wikis/FullSystemPolicy
[full]: https://github.com/roddhjav/apparmor.d/blob/main/apparmor.d/groups/_full
