---
title: Full system policy (FSP)
---

!!! danger

    Full system policy is still under development:

    - It is experimental, and it has only been tested on server.
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
- User interactive shell must be confined. This is done through PAM and Role Based Access Control (RBAC).

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

## Systemd

The profiles dedicated for full system policies are maintained in the **[`_full`][full]** group. Systemd (as PID 1) is the entrypoint of the system, thus in FSP mode, it is also the entry point of the confinement.

```sh
systemd                                # PID 1, entrypoint, requires "Early policy"
├── systemd                            # To restart itself
├── systemd-generators-*               # Systemd system and environment generators
└── sd                                 # Internal service starter and config handler, handles all services
    ├── Px or px,                      # Any service with profile
    ├── Px ->                          # Any service without profile defined in the unit file (see systemd/full/systemd)
    ├── &*                             # Stacked service as defined in the unit file (see systemd/full/systemd)
    ├── sd-mount                       # Handles mount operations from services
    ├── sd-umount                      # Handles unmount operations from services
    ├── sd//systemctl                  # Internal system systemctl
    └── systemd-user                   # Profile for 'systemd --user'
        ├── systemd-user               # To restart itself
        ├── systemd-user-generators-*  # Systemd user and environment generators
        └── sdu                        # Handles all user services
            ├── Px or px,              # Any user service with profile
            ├── Px ->                  # Any user service without profile defined in the unit file (see systemd/full/systemd)
            ├── &*                     # Stacked user service as defined in the unit file (see systemd/full/systemd)
            └── sdu//systemctl         # Internal user systemctl
```
<figure>
    <figcaption>Overall architecture of the systemd profiles stack</figcaption>
</figure>

### Design rationale

The systemd profiles design aims at providing a flexible and secure confinement for systemd and its services while addressing several challenges:

- Differentiate systemd (PID 1) and `system --user`
- Keep `systemd` and `systemd-user` as mininal as possible, and transition to less privileged profiles.
- Allow the executor profiles to handled stacked profiles.
- Most additions need to be done in the `sd`/`sdu` profile, not in `systemd`/`systemd-user`.
- Dedicated `sd-mount` / `sd-umount` profiles for most mount from the unit services.

### Profile `systemd`

The profile for `systemd` (PID 1) does not specify an attachment path because it is directly loaded by systemd thanks to the [early policy](https://gitlab.com/apparmor/apparmor/-/wikis/AppArmorInSystemd#early-policy-loads) feature.

Systemd is (kind of obviously) a highly privileged program. The purpose of this profile is to transition to other less privileged program as soon as possible. It only allows transition to two kinds of profiles:

- The systemd executor (profile named `sd`) it is the systemd internal service starter and config handler.
- The system generators (profiles named `systemd-generators-*`), they are used at boot time to generate systemd unit files based on the current system configuration.

!!! note "Profile requirement"

    To work as intended, all system generators **must** have a profile. For a given distribution, the list of these generators can be found under `/usr/lib/systemd/system-generators/*`

### Profile `sd`

`sd` is a profile for SystemD-executor run as root, it is used to run all services files and to encapsulate stacked services profiles (hence the short name to keep security attribute easy to read). It aims at reducing the size of the main systemd profile.

In an even more secure environment, it can also be used to strictly limit the list of allowed services that can be started by systemd(1).
{ .annotate }

1. **:construction: Work in Progress :construction:**

!!! note "Profile requirement"

    To work as intended, all privileged services **must** have a profile. For a given distribution, the list of these services can be found under:
    ```
    /usr/lib/systemd/system/*.service
    /usr/lib/systemd/system-environment-generators/*
    ```

### Profile `systemd-user`

`sd` allow transition to `systemd-user`, the profile for `systemd --user`. It is only intended to handle user based sessions.

Similarly to `systemd`, it only allows transition to two kinds of profiles:

- The systemd executor (profile named `sdu`) it is the systemd internal **user** service starter and config handler.
- The user generators, they are used at user session start to generate systemd unit files based on the current user configuration.

!!! note "Profile requirement"

    To work as intended, all userland generators **must** have a profile For a given distribution, the list of these services can be found under:
    ```
    /usr/lib/systemd/user-environment-generators/*
    /usr/lib/systemd/user-generators/*
    ```

!!! info "Future Improvements"

    To differentiate user session started with `systemd --user` and a root session also started with `systemd --user`, future improvements will use apparmor namespace and will allow further restrictions of this profile.

### Profile `sdu`

`sdu` is a profile for SystemD-executor run as User, it is used to run all services files and to encapsulate stacked services profiles (hence the short name). It aims at reducing the size of the systemd-user profile.

!!! note "Profile requirement"

    To work as intended, all userland services **must** have a profile For a given distribution. If it is to complex to ensure all services are profiled, you can add rules in a local addition file under `/etc/apparmor.d/usr/sdu.d`.

## Role Based Access Control (RBAC)

In FSP, interactive shell from the user must be confined. This is done through [pam_apparmor](https://gitlab.com/apparmor/apparmor/-/wikis/pam_apparmor). It provides [Role-based access controls (RBAC)](https://en.wikipedia.org/wiki/Role-based_access_control) that can restrict interactive shell to well-defined role. The role needs to be defined. This project ship with a default set of roles, but you can create your own. The default roles are:

- **`user`**: This is the default role. It is used for any user that does not have a specific role defined. It has access to the user home directory and other sensitive files.

- **`admin`**: This role is used for any user that has administrative access. It has access to the system files and directories, but not to the user home directory.

The profiles dedicated for the roles definition are maintained in the **[`_roles`][role]** group.

!!! note

    The roles provided are only examples. It is recommended to create your own roles based on your needs.
    For example, the play machine provides three roles: `root`, `play`, and `deploy`. See the [play machine](play.md) page for more details.

[apparmor-wiki]: https://gitlab.com/apparmor/apparmor/-/wikis/FullSystemPolicy
[full]: https://github.com/roddhjav/apparmor.d/blob/main/apparmor.d/groups/_full
[role]: https://github.com/roddhjav/apparmor.d/blob/main/apparmor.d/groups/_roles
