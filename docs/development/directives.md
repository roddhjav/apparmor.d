---
title: Directives
---

`apparmor.d` supports build directives, they are processed at build time of the project, when running `make`. They are valid apparmor comment, therefore, `apparmor_parser` can be used on a profile even if the directives have not been processed. They should not end with a comma. Multiline directive is not supported.

The directives follow the format:
```sh
#aa:<name> [options]
```

**`<name>`**

:   The name of the directive to apply


**`[options]`**

:   A (possibly empty) list or map of arguments. Exact format depend on the directive.

## Dbus

See the [dbus page](dbus.md#dbus-directive).

    
## Only, Exclude

The `only` and `exclude` directives can be used to filter individual rule or rule paragraph depending on the target distribution of distribution family.

**Format**

```sh
#aa:only <filter>
#aa:exclude <filter>
```

**`<filter>`**

:   The filter to apply. Can be:

    - A supported target distribution: `arch`, `debian`, `ubuntu`, `opensuse`, `whonix`.
    - A supported distribution family: `apt`, `pacman`, `zypper`.

**Example**

!!! note ""

    [apparmor.d/profiles-m-r/packagekitd](https://github.com/roddhjav/apparmor.d/blob/f81ceb91855f194dc53c10c17cbe1d7b50434a1e/apparmor.d/profiles-m-r/packagekitd#L99)
    ``` sh linenums="99"
      #aa:only opensuse
            @{run}/zypp.pid rwk,
      owner @{run}/zypp-rpm.pid rwk,
      owner @{run}/zypp/packages/ r,
    ```

**Generate**

`#aa:only pacman`

:   
    Remove the line/paragraph when the project is not compiled on the Archlinux family.


## Exec

The `exec` directive is useful to allow executing transition to a profile without having to manage the possible long list of profile attachment (it varies depending on the distribution). The directive parse and resolve the attachment variable (`@{exec_path}`) of the target profile and include it in the current profile.

**Format**

```sh
#aa:exec [transition] profiles...
```

**`profiles...`**

:   List of profile **file** that can be executed from the current profile.

**`[transition]`**

:   Optional transition mode (default: `P`). Can be any of: `P`, `U`, `p`, `u`, `PU`, `pu`.


**Example**

!!! note ""

    [apparmor.d/groups/kde/ksmserver](https://github.com/roddhjav/apparmor.d/blob/f81ceb91855f194dc53c10c17cbe1d7b50434a1e/apparmor.d/groups/kde/ksmserver#L29)
    ``` sh linenums="29"
    #aa:exec kscreenlocker_greet
    ```

**Generate**

`#aa:exec baloo`

:   
    ```sh
    @{bin}/baloo_file Px,
    @{lib}/@{multiarch}/{,libexec/}baloo_file Px,
    @{lib}/{,kf6/}baloo_file Px,
    ```


## Stack

[Stacked](https://gitlab.com/apparmor/apparmor/-/wikis/AppArmorStacking) profiles can be hard to maintain. The *parent* profile needs to manage its own rules as well as always include the stacked profile rules. This directive automatically include the stacked profile rules into the parent profile.

**Format**

```sh
#aa:stack profiles...
```

**`profiles...`**

:   List a profile **file** to stack at the end of the current profile.


**Example**

!!! note ""

    [apparmor.d/_full/systemd](https://github.com/roddhjav/apparmor.d/blob/f81ceb91855f194dc53c10c17cbe1d7b50434a1e/apparmor.d/groups/_full/systemd#L150)
    ``` sh linenums="150"
    #aa:stack systemd-networkd systemd-oomd systemd-resolved systemd-timesyncd
    ```

**Generate**

`#aa:stack systemd-oomd`

:   
    ```sh
    # Stacked profile: systemd-oomd
    include <abstractions/bus-system>
    include <abstractions/systemd-common>
    capability dac_override,
    capability kill,
    unix (bind) type=stream addr=@@{hex}/bus/systemd-oomd/bus-api-oom,
    #aa:dbus own bus=system name=org.freedesktop.oom1
    /etc/systemd/oomd.conf r,
    /etc/systemd/oomd.conf.d/{,**} r,
            @{run}/systemd/io.system.ManagedOOM rw,
            @{run}/systemd/io.systemd.ManagedOOM rw,
            @{run}/systemd/notify rw,
    owner @{run}/systemd/journal/socket w,
    @{sys}/fs/cgroup/cgroup.controllers r,
    @{sys}/fs/cgroup/memory.pressure r,
    @{sys}/fs/cgroup/user.slice/user-@{uid}.slice/user@@{uid}.service/memory.* r,
    @{PROC}/pressure/{cpu,io,memory} r,
    include if exists <local/systemd-oomd>
    ```
