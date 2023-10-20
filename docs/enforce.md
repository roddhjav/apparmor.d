---
title: Enforce Mode
---

# Enforce Mode

The default package configuration installs all profiles in *complain* mode. This is a safety measure to ensure you are not going to break your system on initial installation. Once you have tested it, and it works fine, you can easily switch to *enforce* mode. The profiles that are not considered stable are kept in complain mode, they can be tracked in the [`dists/flags`](https://github.com/roddhjav/apparmor.d/tree/main/dists/flags) directory.

!!! warning

    When reporting issue. Please ensure the profiles are in complain mode

## Install

#### :material-arch: Archlinux

In `PKGBUILD`, replace `make` by `make enforce`:
```diff
-  make
+  make enforce
```

#### :material-ubuntu: Ubuntu & :material-debian: Debian

In `debian/rules`, add the following lines:

```make
override_dh_auto_build:
	make enforce
```

#### :simple-suse: OpenSUSE & Partial install

Use the `make enforce` command to build instead of `make`

## Track profiles in complain mode

The [`dists/flags`](https://github.com/roddhjav/apparmor.d/tree/main/dists/flags) directory tracks the profile that have been forced in complain mode. It is used for profile that are not considered stable. Files in this directory should respect the following format: `<profile> <flags>`, flags should be comma separated.

For instance, to move `adb` in complain mode, edit **[`dists/flags/main.flags`](https://github.com/roddhjav/apparmor.d/blob/main/dists/flags/main.flags)** and add the following line:
```sh
adb complain
```

Beware, flags defined in this file overwrite flags in the profile. So you may need to add other flags. Example for `gnome-shell`:
```sh
gnome-shell attach_disconnected,mediate_deleted,complain
```
