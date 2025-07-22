---
title: Enforce Mode
---

The default package configuration installs all profiles in *complain* mode. This is a safety measure to ensure you are not going to break your system on initial installation. Once you have tested it, and it works fine, you can easily switch to *enforce* mode. The profiles that are not considered stable are kept in complain mode, they can be tracked in the [`dists/flags`](https://github.com/roddhjav/apparmor.d/tree/main/dists/flags) directory.

!!! danger

    - You **must** test in complain mode first and ensure your system works as expected.
    - You **must** regularly check AppArmor log with [`aa-log`](usage.md#apparmor-log) and [report](report.md) issues first.
    - When reporting an issue, you **must** ensure the affected profiles are in complain mode.


=== ":material-arch: Archlinux"

    In the `PKGBUILD`, replace `just complain` by `just enforce`:

    ```diff
    -  just complain
    +  just enforce
    ```

    Then, build the package with: `just pkg`

=== ":material-ubuntu: Ubuntu"

    In `debian/rules`, replace `just complain` by `just enforce`:

    ```diff
      override_dh_auto_build:
    -     just complain
      override_dh_auto_build:
    +     just enforce
    ```

    Then, build the package with: `just dpkg`

=== ":material-debian: Debian"
    
    In `debian/rules`, replace `just complain` by `just enforce`:

    ```diff
      override_dh_auto_build:
    -     just complain
      override_dh_auto_build:
    +     just enforce
    ```

    Then, build the package with: `just dpkg`

=== ":simple-suse: openSUSE"

    In `dists/apparmor.d.spec`, replace `just complain` by `just enforce`:

    ```diff
       %build
    -  just complain
       %build
    +  just enforce
    ```

    Then, build the package with: `just rpm`

=== ":material-home: Partial Install"

    Use the `just enforce` command to build instead of `just complain`

[aur]: https://aur.archlinux.org/packages/apparmor.d-git
