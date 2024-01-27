---
title: Dbus
---

All dbus rules are labelled under the name of the given profiles that provide dbus data. If the profiles were going to change (a renaming, an architectural change), the dbus rules need to be updated accordingly.

Default **system**, **session** and **accessibility** bus access are provided with the abstraction:

- `abstractions/bus-system`
- `abstractions/bus-session`
- `abstractions/bus-accessibility`

## Dbus Abstractions

Access to common dbus interface is done using the abstractions under **[`abstractions/bus/`](https://github.com/roddhjav/apparmor.d/tree/main/apparmor.d/abstractions/bus)**. They are kept minimal on purpose. The goal is not to give full talk access an interface but to provide a *read only* like view of it. For more access, use the dbus directive

## Dbus Directive

We use a special directive to generate (when running `make`) more advanced dbus access.

**Directive format**
```
# dbus: ( own | talk ) bus=( system | session ) name=AARE [label=AARE] [interface=AARE]
```

The directive format is on purpose very similar to apparmor dbus rules. However, there is some restrictions:

- `bus` and `name` are mandatory and will break the build if ignored.
- For the *talk* sub directive, profile name under a `label` is also mandatory
- `interface` can optionally be given when it is different to the dbus path.

**Example:**

Allow owning a dbus interface:

!!! note ""

    [apparmor.d/groups/network/NetworkManager](https://github.com/roddhjav/apparmor.d/blob/a3b15973640042af7da0ed540db690c711fbf6ec/apparmor.d/groups/network/NetworkManager#L46)
    ``` aa linenums="46"
    # dbus: own bus=system name=org.freedesktop.NetworkManager
    ```

Allow talking to a dbus interface on a given profile

!!! note ""

    [apparmor.d/groups/gnome/gdm](https://github.com/roddhjav/apparmor.d/blob/a3b15973640042af7da0ed540db690c711fbf6ec/apparmor.d/groups/gnome/gdm#L32)
    ``` aa linenums="32"
    # dbus: talk bus=system name=org.freedesktop.login1 label=systemd-logind
    ```

