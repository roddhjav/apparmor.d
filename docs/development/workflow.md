---
title: Workflow
---

**Workflow to write profiles**

<div class="grid cards" markdown>

-   :material-file-document: &nbsp; **[Write a blanck profile](#add-a-blank-profile)**

</div>
<div class="grid cards" markdown>

-   :material-download: &nbsp; **[Install the profile](#individual-profile)**

</div>
<div class="grid cards" markdown>

-   :material-test-tube: &nbsp; **[Profile the program](#program-profiling)**

</div>
<div class="grid cards" markdown>

-   :octicons-law-16: &nbsp; **[Respect the profile guidelines](guidelines.md)**

</div>


## Add a blank profile

1. To add a new profile `foo`, add the file `foo` in [`apparmor.d/profile-a-f`](https://github.com/roddhjav/apparmor.d/blob/main/apparmor.d/profiles-a-f).
   If your profile is part of a large group of profiles, it can also go in
   [`apparmor.d/groups`](https://github.com/roddhjav/apparmor.d/blob/main/apparmor.d/groups).

2. Write the profile content, the rules depend on the confined program,
   Here is the bare minimum for the program `foo`:
``` sh
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2024 You <your@email>
# SPDX-License-Identifier: GPL-2.0-only

abi <abi/4.0>,

include <tunables/global>

@{exec_path} = @{bin}/foo
profile foo @{exec_path} {
  include <abstractions/base>

  @{exec_path} mr,

  include if exists <local/foo>
}

# vim:syntax=apparmor
```

## Development Install

It is not recommended installing the full project *"manually"* (with `make`, `sudo make install`). The distribution specific packages are intended to be used in development as they include additional rule to ensure compatibility with upstream (see `debian/`, `PKGBUILD` and `dists/apparmor.d.spec`).

Instead, install an individual profile or the development package, the following way.

### Development package

=== ":material-arch: Archlinux"

    ```sh
    make pkg
    ```

=== ":material-ubuntu: Ubuntu"

    ```sh
    make dpkg
    ```

=== ":material-debian: Debian"

    ```sh
    make dpkg
    ```

=== ":simple-suse: openSUSE"

    ```sh
    make rpm
    ```

=== ":material-docker: Docker"

    For any system with docker installed you can simply build the package with:

    ```sh
    make package dist=<distribution>
    ```

    Then you can install the package with `dpkg`, `pacman` or `rpm`.

### Individual profile

**Format**

```sh
make dev name=<profile-name>
```

**Exampe**

:   Testing the profile `pass`

    ```
    make dev name=pass
    ```

    This:

    - Prebuild the `pass` profile in complain mode to `.build`,
    - Install the profile to `/etc/apparmor.d/`
    - Load the profile by restarting the AppArmor service.


More advanced development, like editing the abstractions or working over multiple profiles at the same time requires installing the full development package.

For this individual profile installation to work, the full package needs to be installed, regardless of the installation method ([dev](#development-package) or [stable](../install.md)).

## Program Profiling

### Workflow

To discover the access needed by a program, you can use the following tools:

1. Star the program in *complain* mode, let it initialize itself, then close it.

1. Run **[`aa-log -r`](../usage.md#apparmor-log)**. It will:
    - Convert the logs to AppArmor rules.
    - Detect if flags such as `attach_disconnected` are needed.
    - Convert all common paths to **[variables](../variables.md)**.

1. From `aa-log` output, you can:
    - Copy the rules to the profile.
    - Replace some rules with **[abstractions](abstractions.md)** as 80% of the rules should already be covered by an abstraction.

1. Then, [update the profile](#individual-profile) and start the program again. Use the program as you would normally do, but also try to run all the features of the program, e.g.: open the help, settings, etc.

1. Run **[`aa-log`](../usage.md#apparmor-log)**. Stop the program as long as you get over 100 new rules. Add the rules to the profile.

After 2 or 3 iterations, you should have a working profile.

### Recommendations

<div class="grid cards" markdown>

-   :material-function: &nbsp; **[Use the abstractions](abstractions.md)**
-   :simple-files: &nbsp; **[Learn how to open resources](internal.md#open-resources)**
-   :fontawesome-solid-bus-simple: &nbsp; **[Learn how Dbus rules are handled](dbus.md)**
-   :material-sign-direction: &nbsp; **[Learn about directives `#aa:`](directives.md)**
-   :octicons-law-16: &nbsp; **[Follow the profile guidelines](guidelines.md)**
-   :octicons-light-bulb-16: &nbsp; **[See other recommendations](recommendations.md)**

</div>

!!! danger "Warning"

    Following the [profile guidelines](guidelines.md) is **mandatory** for all profiles. PRs that do not follow the guidelines will not get merged.

### Tools

* **[aa-notify](https://wiki.archlinux.org/title/AppArmor#Get_desktop_notification_on_DENIED_actions)** is a tool that will allow you to get notified on every apparmor log.

* **[aa-logprof](https://man.archlinux.org/man/aa-logprof.8)** is another tool that will help you to generate a profile from logs. However, the logs generated by `aa-logprof` need to be rewritten to comply with the profile [guidelines](guidelines.md).

* **[aa-complain](https://man.archlinux.org/man/aa-complain.8), aa-enforce** are tools to quickly change the mode of a profile.


## Development Settings

### Profile flags

Flags for all profiles in this project are tracked under the [`dists/flags`](https://github.com/roddhjav/apparmor.d/tree/main/dists/flags) directory. It is used for profile that are not considered stable. Files in this directory should respect the following format: `<profile> <flags>`, flags should be comma separated.

For instance, to move `adb` in *complain* mode, edit **[`dists/flags/main.flags`](https://github.com/roddhjav/apparmor.d/blob/main/dists/flags/main.flags)** and add the following line:
```sh
adb complain
```

Beware, flags defined in this file overwrite flags in the profile. So you may need to add other flags. Example for `gnome-shell`:
```sh
gnome-shell attach_disconnected,mediate_deleted,complain
```


### Ignore profiles

It can be handy to not install a profile for a given distribution. Profiles and directories to ignore are tracked under the [`dists/ignore`](https://github.com/roddhjav/apparmor.d/tree/main/dists/ignore) directory. Files in this directory should respect the following format: `<profile or path>`. One ignore by line. It can be a profile name or a directory to ignore (relative to the project root).
