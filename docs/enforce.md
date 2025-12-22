---
title: Enforce Mode
---

The default package configuration installs all profiles in *complain* mode. This is a safety measure to ensure you are not going to break your system on initial installation. Once you have tested it, and it works fine, you can easily switch to *enforce* mode. The profiles that are not considered stable are kept in complain mode, they can be tracked in the [`dists/flags`](https://github.com/roddhjav/apparmor.d/tree/main/dists/flags) directory.

!!! danger

    - You **must** test in complain mode first and ensure your system works as expected.
    - You **must** regularly check AppArmor log with [`aa-log`](usage.md#apparmor-log) and [report](report.md) issues first.
    - When reporting an issue, you **must** ensure the affected profiles are in complain mode.


**Prerequisite**

As the `enforced` version of the package conficts with the default `apparmor.d` package, you need to uninstall it first:

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


**Installation**

=== ":material-arch: Archlinux"

    `apparmor.d.enforced` is available in the [Arch User Repository][aur]:

    ```sh
    yay -S apparmor.d.enforced  # or your preferred AUR install method
    ```

=== ":material-ubuntu: Ubuntu"

    Using the [pkg.pujol.io][repo] debian repository, install the package:
    ```sh
    sudo apt install apparmor.d.enforced
    ```


=== ":material-debian: Debian"

    Using the [pkg.pujol.io][repo] debian repository, install the package:
    ```sh
    sudo apt install apparmor.d.enforced
    ```

=== ":simple-suse: openSUSE"

    openSUSE users need to add [cboltz](https://en.opensuse.org/User:Cboltz) repo on OBS:

    ```sh
    zypper install apparmor.d.enforced
    ```

=== ":material-home: Partial Install"

    Use the `just enforce` command to build instead of `just complain`

[aur]: https://aur.archlinux.org/packages/apparmor.d-git
[repo]: https://pkg.pujol.io
