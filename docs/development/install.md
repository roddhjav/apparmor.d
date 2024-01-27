---
title: Installation
---

## Development Install

!!! warning

    Do **not** install this project *"manually"* (with `make`, `sudo make install`). The distribution specific packages are intended to be used in development as they include additional rule to ensure compatibility with upstream. 

    See `debian/`, `PKGBUILD` and `dists/apparmor.d.spec`.


**:material-docker: Docker**

From any distribution, if you have docker installed, you can simply build the package with:
```sh
make package dist=<distribution>
```
Then you can install the package with `dpkg`, `pacman` or `rpm`.

**:material-arch: Archlinux**
```sh
make pkg
```

**:material-ubuntu: Ubuntu & :material-debian: Debian**
```sh
make dpkg
```

**:simple-suse: OpenSUSE**
```sh
make rpm
```


## Profile flags

Flags for all profiles in this project are tracked under the [`dists/flags`](https://github.com/roddhjav/apparmor.d/tree/main/dists/flags) directory. It is used for profile that are not considered stable. Files in this directory should respect the following format: `<profile> <flags>`, flags should be comma separated.

For instance, to move `adb` in complain mode, edit **[`dists/flags/main.flags`](https://github.com/roddhjav/apparmor.d/blob/main/dists/flags/main.flags)** and add the following line:
```sh
adb complain
```

Beware, flags defined in this file overwrite flags in the profile. So you may need to add other flags. Example for `gnome-shell`:
```sh
gnome-shell attach_disconnected,mediate_deleted,complain
```


## Ignore profiles

It can be handy to not install a profile for a given distribution. Profile or directory to ignore are tracked under the [`dists/ignore`](https://github.com/roddhjav/apparmor.d/tree/main/dists/ignore) directory. Files in this directory should respect the following format: `<profile or path>`. One ignore by line. It can be a profile name or a directory to ignore (relative to the project root).
