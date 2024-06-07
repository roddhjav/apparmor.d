---
title: Installation
---

!!! warning

    To prevent the risk of breaking your system, the default package configuration installs all profiles in complain mode. They can be enforced later. See the [Enforce Mode](enforce.md) page.

!!! danger

    Do **not** expect this project to work correctly if your Desktop Environment and Display Manager are not supported. Your Desktop Environment or Display Manager might not load, and that would be a feature.

## Requirements

**AppArmor**

An `AppArmor` supported Linux distribution is required. The default profiles and abstractions shipped with AppArmor must be installed.

**Desktop environment**

The following desktop environments are supported:

  - [x] :material-gnome: Gnome
  - [x] :simple-kde: KDE
  - [ ] :simple-xfce: XFCE *(work in progress)*

**Build dependency**

* Go >= 1.18

## :material-arch: Arch Linux

`apparmor.d-git` is available in the [Arch User Repository][aur]:
```
yay -S apparmor.d-git  # or your preferred AUR install method
```

Or without an AUR helper:
```sh
git clone https://aur.archlinux.org/apparmor.d-git.git
cd apparmor.d-git
makepkg -si
```


## :material-ubuntu: Ubuntu & :material-debian: Debian

Build the package from sources:
```sh
sudo apt install apparmor-profiles build-essential config-package-dev debhelper golang-go rsync git
git clone https://github.com/roddhjav/apparmor.d.git
cd apparmor.d
dpkg-buildpackage -b -d --no-sign
sudo dpkg -i ../apparmor.d_*.deb
```

!!! tip

    If you have `devscripts` installed, you can use the one liner:
    ```sh
    make dpkg
    ```

!!! note

    Debian user may need golang from the backports repository to build:
    ```sh
    echo 'deb http://deb.debian.org/debian bookworm-backports main contrib non-free' | sudo tee -a /etc/apt/sources.list
    sudo apt update
    sudo apt install -t bookworm-backports golang-go
    ```

!!! warning

    **Beware**: do not install a `.deb` made for Debian on Ubuntu, the packages are different.

    If your distribution is based on Ubuntu or Debian, you may want to manually set the target distribution by exporting `DISTRIBUTION=debian` if is Debian based, or `DISTRIBUTION=ubuntu` if it is Ubuntu based.

## :simple-suse: OpenSUSE

OpenSUSE users need to add [cboltz](https://en.opensuse.org/User:Cboltz) repo on OBS
```sh
zypper addrepo https://download.opensuse.org/repositories/home:cboltz/openSUSE_Factory/home:cboltz.repo
zypper refresh
zypper install apparmor.d
```


## Partial install

For test purposes, you can install specific profiles with the following commands. Abstractions, tunable, and most of the OS dependent post-processing is managed.

```sh
make
sudo make profile-names...
```

!!! warning

    Partial installation is discouraged because profile dependencies are not fetched. To prevent some AppArmor issues, the dependencies are automatically switched to unconfined (`rPx` -> `rPUx`). The installation process warns on the missing profiles so that you can easily install them if desired. (PR is welcome see [#77](https://github.com/roddhjav/apparmor.d/issues/77))

    For instance, `sudo make pass` gives:
    ```sh
    Warning: profile dependencies fallback to unconfined.
    @{bin}/wl-{copy,paste} rPx,
    @{bin}/xclip           rPx,
    @{bin}/python3.@{int} rPx -> pass-import,  # pass-import
        @{bin}/pager         rPx -> child-pager,
        @{bin}/less          rPx -> child-pager,
        @{bin}/more          rPx -> child-pager,
    '.build/apparmor.d/pass' -> '/etc/apparmor.d/pass'
    ```
    So, you can install the additional profiles `wl-copy`, `xclip`, `pass-import`, and `child-pager` if desired.


## Uninstall

- :material-arch: Arch Linux `sudo pacman -R apparmor.d`
- :material-ubuntu: Ubuntu & :material-debian: Debian `sudo apt purge apparmor.d`
- :simple-suse: OpenSUSE `sudo zypper remove apparmor.d`

[aur]: https://aur.archlinux.org/packages/apparmor.d-git
[repo]: https://repo.pujol.io/
[keys]: https://repo.pujol.io/gpgkey
