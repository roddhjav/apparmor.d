---
title: Abstractions
---

# `abstractions`

Use of dangerous or deprecated abstractions

## Problematic rule

```sh
# WRONG
include <abstractions/nameservice>
```

## Correct rule

```sh
include <abstractions/nameservice-strict>
```

## Rationale

Some abstractions provide more access than required, do not integrate with profiles defined in apparmor.d or with non-Ubuntu systems.

The following abstractions are considered dangerous:

- `dbus`: Full dbus access
- `dbus-accessibility`: Full dbus accessibility access
- `dbus-session`: Full dbus session access
- `dbus-system`: Full dbus system access
- `user-tmp`: Full access to user temporary files (See [too-wide](too-wide.md) check)

Deprecated abstractions:

- `bash` -> `shell`: `bash` does not cover all shells.
- `nameservice` -> `nameservice-strict`: `nameservice` gives network access which is not required in most cases.

Deprecated abstractions, would conflict with apparmor.d integration

- `dbus-accessibility-strict` -> `bus-accessibility`
- `dbus-network-manager-strict` -> `network-manager-observe`
- `dbus-session-strict` -> `bus-session`
- `dbus-system-strict` -> `bus-system`
- `gnome` -> `gnome-strict`
- `kde` -> `kde-strict`
- `X` -> `X-strict`

## Exceptions

None
