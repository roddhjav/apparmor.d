---
title: Dbus
---

All dbus rules are labelled under the name of the given profiles that provide dbus data. It is one of the value added by this project, as we have profiles for *everything*, we can restrict the bus further by limiting connection to a given peer label (the profile name). In the case of renaming a profile, all dbus rules related in this profile need to be updated accordingly.

## Profiles

Regardless of the Dbus implementation used (`dbus-daemon` or `dbus-broker`), all dbus daemons are handled under the same set of profiles: [`dbus-system`](https://github.com/roddhjav/apparmor.d/blob/main/apparmor.d/groups/bus/dbus-system), [`dbus-session`](https://github.com/roddhjav/apparmor.d/blob/main/apparmor.d/groups/bus/dbus-session), and [`dbus-accessibility`](https://github.com/roddhjav/apparmor.d/blob/main/apparmor.d/groups/bus/dbus-accessibility). This structure largely improves the confinement of each profile.

To ensure the system and session bus are handled by a different profile, a [systemd drop-in](https://github.com/roddhjav/apparmor.d/blob/main/systemd/default/system/dbus.service) configuration file is used to set the specific dbus profile that a dbus service must use.

## Abstractions

### Base

Default **system**, **session**, and **accessibility** bus access are provided with the following abstractions:

- `abstractions/bus-system`
- `abstractions/bus-session`
- `abstractions/bus-accessibility`

Do not use the dbus abstractions from apparmor in this project, they won't work as expected as the dbus daemon is confined. Furthermore, in `apparmor.d` there is no such thing as a strict dbus abstraction (`abstractions/dbus-strict`) as they are strict by default: bus access needs to be explicitly allowed using an interface abstraction or a directive.

### Interfaces

Access to common dbus interfaces is done using the abstractions under **[`abstractions/bus/`](https://github.com/roddhjav/apparmor.d/tree/main/apparmor.d/abstractions/bus)**. They are kept minimal on purpose. The goal is not to give full talk access an interface but to provide a *read-only* like view of it. It may be required to have a look at the dbus interface documentation to check what method can be safely allowed.

For more access, simply use the [`aa:dbus talk`](#dbus-directive) directive.

There is a trade of between security and maintenance to make:

- `aa:dbus talk` will generate less issue as it gives full talk access
- `abstractions/bus/*` will provide more restriction, and possibly more issue. In the future, these rules will be automatically generated from the interface documentation.

## Dbus Directive

We use a special [directive](directives.md) to generate more advanced dbus access. The directive format is on purpose very similar to the AppArmor dbus rule.

**Format**

```sh
#aa:dbus <access> bus=<bus> name=<name> [label=AARE] [interface=AARE] [path=AARE]
```

**`<access>`**

:    Access type. Can be `own` or `talk`:

     - `own` means the profile owns the dbus interface. It is allowed to send and receive from anyone on this interface. It should only be used for profile owning the dbus interface.
     - `talk` means the profile can talk on a given interface to the profile that owns it (a label must be given under the `label` option). It should only be used when full access to an interface is required.

**`<bus>`**

:    Dbus bus, can be `system`, `session` or `accessibility`.

**`<name>`**

:    Dbus interface name.

**`[label=AARE]`**

:    Name of the profile. Mandatory for `talk` access.

**`[interface=AARE]`**

:    Can optionally be given when it is different to the dbus path.

**`[path=AARE]`**

:    Can optionally be given when it is different to the dbus name.


Note: `<access>`, `<bus>`, and `<name>` are mandatory and will break the build if ignored.


**Example**

Allow owning a dbus interface:

!!! note ""

    [apparmor.d/groups/network/NetworkManager](https://github.com/roddhjav/apparmor.d/blob/f81ceb91855f194dc53c10c17cbe1d7b50434a1e/apparmor.d/groups/network/NetworkManager#L45)
    ``` sh linenums="45"
    #aa:dbus own bus=system name=org.freedesktop.NetworkManager
    ```

Allow talking to a dbus interface on a given profile:

!!! note ""

    [apparmor.d/groups/gnome/gdm](https://github.com/roddhjav/apparmor.d/blob/f81ceb91855f194dc53c10c17cbe1d7b50434a1e/apparmor.d/groups/gnome/gdm#L44)
    ``` sh linenums="34"
    #aa:dbus talk bus=system name=org.freedesktop.login1 label=systemd-logind
    ```

**Generate**

`#aa:dbus own bus=system name=org.freedesktop.NetworkManager`

:    
     ```sh
     dbus bind bus=system name=org.freedesktop.NetworkManager{,.*},
     dbus receive bus=system path=/org/freedesktop/NetworkManager{,/**}
          interface=org.freedesktop.NetworkManager{,.*}
          peer=(name=":1.@{int}"),
     dbus receive bus=system path=/org/freedesktop/NetworkManager{,/**}
          interface=org.freedesktop.DBus.Properties
          peer=(name=":1.@{int}"),
     dbus receive bus=system path=/org/freedesktop/NetworkManager{,/**}
          interface=org.freedesktop.DBus.ObjectManager
          peer=(name=":1.@{int}"),
     dbus send bus=system path=/org/freedesktop/NetworkManager{,/**}
          interface=org.freedesktop.NetworkManager{,.*}
          peer=(name="{:1.@{int},org.freedesktop.DBus}"),
     dbus send bus=system path=/org/freedesktop/NetworkManager{,/**}
          interface=org.freedesktop.DBus.Properties
          peer=(name="{:1.@{int},org.freedesktop.DBus}"),
     dbus send bus=system path=/org/freedesktop/NetworkManager{,/**}
          interface=org.freedesktop.DBus.ObjectManager
          peer=(name="{:1.@{int},org.freedesktop.DBus}"),
     dbus receive bus=system path=/org/freedesktop/NetworkManager{,/**}
          interface=org.freedesktop.DBus.Introspectable
          member=Introspect
          peer=(name=":1.@{int}"),
     ```

`#aa:dbus talk bus=system name=org.freedesktop.login1 label=systemd-logind`

:    
     ```sh
     dbus send bus=system path=/org/freedesktop/login1{,/**}
          interface=org.freedesktop.login1{,.*}
          peer=(name="{:1.@{int},org.freedesktop.login1{,.*}}", label=systemd-logind),
     dbus send bus=system path=/org/freedesktop/login1{,/**}
          interface=org.freedesktop.DBus.Properties
          peer=(name="{:1.@{int},org.freedesktop.login1{,.*}}", label=systemd-logind),
     dbus send bus=system path=/org/freedesktop/login1{,/**}
          interface=org.freedesktop.DBus.ObjectManager
          peer=(name="{:1.@{int},org.freedesktop.login1{,.*}}", label=systemd-logind),
     dbus receive bus=system path=/org/freedesktop/login1{,/**}
          interface=org.freedesktop.login1{,.*}
          peer=(name="{:1.@{int},org.freedesktop.login1{,.*}}", label=systemd-logind),
     dbus receive bus=system path=/org/freedesktop/login1{,/**}
          interface=org.freedesktop.DBus.Properties
          peer=(name="{:1.@{int},org.freedesktop.login1{,.*}}", label=systemd-logind),
     dbus receive bus=system path=/org/freedesktop/login1{,/**}
          interface=org.freedesktop.DBus.ObjectManager
          peer=(name="{:1.@{int},org.freedesktop.login1{,.*}}", label=systemd-logind),
     dbus send bus=system path=/org/freedesktop/Accounts{,/**}
          interface=org.freedesktop.Accounts{,.*}
     ```
