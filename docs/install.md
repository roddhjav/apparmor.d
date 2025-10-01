---
title: Installation
---

## Setup

!!! danger

    Do **not** expect this project to work correctly on your desktop if your Desktop Environment (DE) and Display Manager (DM) are not supported. Your DE/DM might not load, and that would be a **feature**.

Due to the development stage of this project, the default package configuration installs all profiles in **complain** mode. The recommended installation workflow is as-follow:

1. **[Configure AppArmor](#configure-apparmor)** AppArmor for *apparmor.d*.
1. **[Install](#installation)** *apparmor.d* in the (default) complain mode.
1. **[Configure your personal directories](configuration.md)**.
1. Reboot your system.
1. You **must** check for any AppArmor logs with [`aa-log`](usage.md#apparmor-log).
1. **[Report](https://apparmor.pujol.io/report/)** any raised logs.
1. Use the profiles in *complain* mode for a while (a week), regularly check for new AppArmor logs.
1. Only if there are no logs raised for your daily usage, install it in [enforce mode](enforce.md).


## Requirements

**AppArmor**

An `AppArmor` supported Linux distribution is required. The default profiles and abstractions shipped with AppArmor must be installed.

**Desktop environment**

The following desktop environments are supported:

- [x] :material-gnome: Gnome (GDM)
- [x] :simple-kde: KDE (SDDM)
- [ ] :simple-xfce: XFCE (Lightdm) *(work in progress)*

**Build dependency**

* Go >= 1.23
* [just](https://github.com/casey/just) >= 1.40.0


## Configure AppArmor

As there are a lot of rules (~100k lines), it is recommended to enable fast caching compression of AppArmor profiles. [Early policy](https://gitlab.com/apparmor/apparmor/-/wikis/AppArmorInSystemd#early-policy-loads) load **must** also be enabled.

In `/etc/apparmor/parser.conf` ensure you have:
```
write-cache
cache-loc /etc/apparmor/earlypolicy/
Optimize=compress-fast
```

Or run:
```sh
echo 'write-cache' | sudo tee -a /etc/apparmor/parser.conf
echo 'cache-loc /etc/apparmor/earlypolicy/' | sudo tee -a /etc/apparmor/parser.conf
echo 'Optimize=compress-fast' | sudo tee -a /etc/apparmor/parser.conf
```


## Installation

=== ":material-arch: Archlinux"

    `apparmor.d-git` is available in the [Arch User Repository][aur]:

    ```sh
    yay -S apparmor.d-git  # or your preferred AUR install method
    ```

    Or without an AUR helper:

    ```sh
    git clone https://aur.archlinux.org/apparmor.d-git.git
    cd apparmor.d-git
    makepkg -si
    ```

=== ":material-ubuntu: Ubuntu"

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
        just dpkg
        ```

    !!! note

        **Ubuntu 24.04 user will need to:**

        Install [just](https://github.com/casey/just). E.g:
        ```sh
        pipx install rust-just
        ```

    !!! warning

        **Beware**: do not install a `.deb` made for Debian on Ubuntu as the packages are different.

        If your distribution is based on Ubuntu, you may want to manually set the target distribution by exporting `DISTRIBUTION=ubuntu`.

=== ":material-debian: Debian"

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
        just dpkg
        ```

    !!! note

        **Debian 12 user will need to:**

        1. Install Golang from the backports repository:
        ```sh
        echo 'deb http://deb.debian.org/debian bookworm-backports main contrib non-free' | sudo tee -a /etc/apt/sources.list
        sudo apt update
        sudo apt install -t bookworm-backports golang-go
        ```

        2. Install [just](https://github.com/casey/just) locally, and ignore the dependence. E.g:
        ```sh
        pipx install rust-just
        sed '/just/d' -i debian/control
        ```

    !!! warning

        **Beware**: do not install a `.deb` made for Ubuntu on Debian as the packages are different.

        If your distribution is based on Debian, you may want to manually set the target distribution by exporting `DISTRIBUTION=debian`.

=== ":simple-suse: openSUSE"

    openSUSE users need to add [cboltz](https://en.opensuse.org/User:Cboltz) repo on OBS:

    ```sh
    zypper addrepo https://download.opensuse.org/repositories/home:cboltz/openSUSE_Factory/home:cboltz.repo
    zypper refresh
    zypper install apparmor.d
    ```

=== ":material-home: Partial"

    For test purposes, you can install specific profiles with the following commands. Abstractions, tunable, and most of the OS dependent post-processing is managed.

    ```sh
    just complain
    sudo just local profile-names...
    ```

    !!! warning

        Partial installation is discouraged because profile dependencies are not fetched. To prevent some AppArmor issues, the dependencies are automatically switched to unconfined (`rPx` -> `rPUx`). The installation process warns on the missing profiles so that you can easily install them if desired. (PR is welcome see [#77](https://github.com/roddhjav/apparmor.d/issues/77))

        For instance, `sudo just local pass` gives:
        ```sh
        Warning: profile dependencies fallback to unconfined.
        @{bin}/wl-{copy,paste} rPx,
        @{bin}/xclip           rPx,
        @{python_path}         rPx -> pass-import,  # pass-import
            @{pager_path}        rPx -> child-pager,
        '.build/apparmor.d/pass' -> '/etc/apparmor.d/pass'
        ```
        So, you can install the additional profiles `wl-copy`, `xclip`, `pass-import`, and `child-pager` if desired.


[Next: Configure your personal directories](configuration.md){ .md-button .md-button--primary }


## Uninstallation

=== ":material-arch: Archlinux"

    ```sh
    sudo pacman -R apparmor.d
    ```

=== ":material-ubuntu: Ubuntu"

    ```sh
    sudo apt purge apparmor.d
    ```

=== ":material-debian: Debian"

    ```sh
    sudo apt purge apparmor.d
    ```

=== ":simple-suse: openSUSE"

    ```sh
    sudo zypper remove apparmor.d
    ```

[aur]: https://aur.archlinux.org/packages/apparmor.d-git
