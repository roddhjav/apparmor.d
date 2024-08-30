---
title: Enforce Mode
---

The default package configuration installs all profiles in *complain* mode. This is a safety measure to ensure you are not going to break your system on initial installation. Once you have tested it, and it works fine, you can easily switch to *enforce* mode. The profiles that are not considered stable are kept in complain mode, they can be tracked in the [`dists/flags`](https://github.com/roddhjav/apparmor.d/tree/main/dists/flags) directory.

!!! danger

    - You **must** test in complain mode first and ensure your system works as expected.
    - You **must** regularly check AppArmor log with [`aa-log`](usage.md#apparmor-log) and [report](report.md) issues first.
    - When reporting an issue, you **must** ensure the affected profiles are in complain mode.


=== ":material-arch: Archlinux"

    In the `PKGBUILD`, replace `make` by `make enforce`:

    ```diff
    -  make DISTRIBUTION=arch
    +  make enforce DISTRIBUTION=arch
    ```

    Then, build the package with: `make pkg`

=== ":material-ubuntu: Ubuntu"

    In `debian/rules`, add the following lines:

    ```make
    override_dh_auto_build:
        make enforce
    ```

    Then, build the package with: `make dpkg`

=== ":material-debian: Debian"
    
    In `debian/rules`, add the following lines:

    ```make
    override_dh_auto_build:
        make enforce
    ```

    Then, build the package with: `make dpkg`

=== ":simple-suse: openSUSE"

    In `dists/apparmor.d.spec`, replace `%make_build` by `%make_build enforce`

    ```diff
    -  %make_build
    +  %make_build enforce
    ```

    Then, build the package with: `make rpm`

=== ":material-home: Partial Install"

    Use the `make enforce` command to build instead of `make`

[aur]: https://aur.archlinux.org/packages/apparmor.d-git
