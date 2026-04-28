---
title: Bin
---

# `bin` / `sbin`

Use of incorrect binary path in rules.

## Problematic rule

```sh
# WRONG
@{bin}/cron Px,
```

```sh
# WRONG
@{sbin}/pass Px,
```

## Correct rule

```sh
@{sbin}/cron Px,
```

```sh
@{bin}/pass Px,
```

## Rationale

To differentiate between system binaries and administrator binaries, `apparmor.d` uses two separate variables: `@{bin}` for regular binaries and `@{sbin}` for system binaries.

The list of known path in `/usr/sbin` is maintained under the `sbin.list` file.

## Exceptions

Some binaries may be installed in both @{bin} and @{sbin} depending on the package it is installed from. For instance, upstream docker package installs `dockerd` in `/usr/bin/` while the distribution package installs it in `/usr/sbin/`. In such cases, both paths is required.
