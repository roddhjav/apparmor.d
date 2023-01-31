---
title: Installation
---

!!! danger

    In order to not break your system, the default package configuration installs
    all profiles in complain mode. They can be enforced later.
    See the [Enforce Mode](/enforce) page.

## Requirements

**AppArmor**

An `apparmor` based Linux distribution is required. The basic profiles and
abstractions shipped with AppArmor must be installed.

**Desktop environment**

The following desktop environments are supported:

  - [x] :material-gnome: Gnome

Also, please note wayland has better support than xorg.

**Build dependencies**

* Go
* Rsync

## :material-arch: Archlinux

`apparmor.d-git` is available in the [Arch User Repository][aur]:
```sh
git clone https://aur.archlinux.org/apparmor.d-git.git
cd apparmor.d-git
makepkg -s
sudo pacman -U apparmor.d-*.pkg.tar.zst \
  --overwrite etc/apparmor.d/tunables/global \
  --overwrite etc/apparmor.d/tunables/xdg-user-dirs \
  --overwrite etc/apparmor.d/abstractions/trash
```

The overwrite options are only required on the first install. You can use `yay`
or your preferred AUR install method to update it.

!!! note

    The following Archlinux based distributions are supported:

    - [x] CachyOS
    - [x] EndeavourOS
    - [x] :material-manjaro: Manjaro Linux


## :material-ubuntu: Ubuntu & :material-debian: Debian


Build the package from sources:
```sh
sudo apt install apparmor-profiles build-essential config-package-dev debhelper golang-go rsync git
git clone https://github.com/roddhjav/apparmor.d.git
cd apparmor.d
dpkg-buildpackage -b -d --no-sign
sudo dpkg -i ../apparmor.d_*_all.deb
```


## Partial install

!!! warning

    Partial installation is discouraged because profile dependencies are
    not fetched. You may need to either switch desired `rPx` rules to `rPUx`
    (fallback to unconfined) or install these related profiles.
    (PR is welcome see [#77](https://github.com/roddhjav/apparmor.d/issues/77))

For test purposes, you can install a specific profile with the following commands.
Abstractions, tunables, and most of the OS dependent post-processing is managed.

```sh
./configure --complain
make
sudo make profile-names...
```

[aur]: https://aur.archlinux.org/packages/apparmor.d-git
[repo]: https://repo.pujol.io/
[keys]: https://repo.pujol.io/gpgkey
