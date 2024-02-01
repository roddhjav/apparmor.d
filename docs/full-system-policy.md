---
title: Full system policy
---

!!! danger

    Full system policy is still under early development:
    
    - Do not run it outside a development VM! 
    - This is an **advanced** feature, you should understand what you are doing

    **You have been warned!!!**

!!! quote

    AppArmor is also capable of being used for full system policy where processes are by default not running under the `unconfined` profile. This might be useful for high security environments or embedded systems.

    *Source: [AppArmor Wiki][apparmor-wiki]*


## Install


This feature is only enabled when the project is built with `make full`. [Early policy](https://gitlab.com/apparmor/apparmor/-/wikis/AppArmorInSystemd#early-policy-loads) load **must** also be enabled. Once `apparmor.d` has been installed in FSP mode, it is required to reboot to apply the changes.

In `/etc/apparmor/parser.conf` ensure you have:
```
write-cache
cache-loc /etc/apparmor/earlypolicy/
Optimize=compress-fast
```

**:material-arch: Archlinux**

In `PKGBUILD`, replace `make` by `make full`:
```diff
-  make
+  make full
```

**:material-ubuntu: Ubuntu & :material-debian: Debian**

In `debian/rules`, add the following lines:

```make
override_dh_auto_build:
	make full
```

**:simple-suse: OpenSUSE**

In `dists/apparmor.d.spec`, replace `%make_build` by `make full`
```diff
-  %make_build
+  %make_build full
```

**Partial install**

Use the `make full` command to build instead of `make`



## Structure

The profiles dedicated for full system policies are maintained in the **[`_full`][full]** group.

### Systemd

In addition to systemd services (`systemd-*`) that have their own profiles, systemd itself, is confined using:

- [x] **`systemd`**: For systemd as PID 1, designed such as:
     - It allows internal systemd access,
     - It allows starting all common root services.

- [ ] **`systemd-user`**: For `systemd --user`, designed such as:
     - It allows internal systemd user access,
     - It allows starting all common user services.

Both `systemd` and `systemd-user` should not fallback at all. I am working on some new profiles to ensure that 99% of program started by system have a profile.

These profiles are only intended to confine themselves. Any services started by systemd must have their corresponding profile. It means that for a given distribution, the following services must have profiles:

- [x] For `systemd`:
```sh
/usr/lib/systemd/system-generators/*
/usr/lib/systemd/system-environment-generators/*
/usr/lib/systemd/system/*.service
```

- [ ] For `systemd-user`
```sh
/usr/lib/systemd/user-environment-generators/*
/usr/lib/systemd/user-generators/*
/usr/lib/systemd/user/*.service
```

To be allowed to run, additional root or user services may need to add extra rules inside the `usr/systemd.d` or `usr/systemd-user.d` directory. For example, when installing a new privileged service `foo` with [stacking](#no-new-privileges) you may need to add the following to `/etc/apparmor.d/usr/systemd.d/foo`:
```
  @{lib}/foo rPx -> systemd//&foo,
  ...
```

### Fallback

In addition to the `systemd` profiles, a full system policy needs to ensure that no program run in an unconfined state at any time. The fallbacks profiles consist of a set generic specialized profiles:

- **`default`** is used for any *classic* user application with a GUI. It has full access to user home directories.
- **`bwrap`, `bwrap-app`** are used for *classic* user application that are sandboxed with **bwrap**.

!!! warning

    The main fallback profile (`default`) is not intended to be used by priviligied program or service. Such programs **must** have they dedicaded profile and would break otherwise.

Additionally, special user access can be setup using PAM rules set such as a random shell interactively opened (as user or as root). 

[apparmor-wiki]: https://gitlab.com/apparmor/apparmor/-/wikis/FullSystemPolicy
[full]: https://github.com/roddhjav/apparmor.d/blob/main/apparmor.d/groups/_full
