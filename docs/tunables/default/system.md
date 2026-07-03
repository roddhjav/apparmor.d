---
title: System
tags:
  - tunables
  - default
---

## apparmorfs

### @{apparmorfs}

[:material-file-eye:](https://gitlab.com/apparmor/apparmor/-/blob/9a402e9e6a04b45dd788c7dfb0106dae45443b01/profiles/apparmor.d/tunables/apparmorfs#L11 "View source"){ .abs-source }

```
@{apparmorfs}=@{securityfs}/apparmor/
```

## etc

@{etc_ro} contains a space-separated list of the system configuration directories.
Traditionally this means /etc/, but when using a read-only / filesystem and/or
with the goal of having only user-modified config files in /etc/, directories
like /usr/etc/ get introduced for storing the default config.

### @{etc_ro}

[:material-file-eye:](https://gitlab.com/apparmor/apparmor/-/blob/9a402e9e6a04b45dd788c7dfb0106dae45443b01/profiles/apparmor.d/tunables/etc#L18 "View source"){ .abs-source }

@{etc_ro} contains directories with configuration files, including read-only directories.
Do not use @{etc_ro} in rules that allow write access.

```
@{etc_ro}=/etc/ /usr/etc/
```
### @{etc_rw}

[:material-file-eye:](https://gitlab.com/apparmor/apparmor/-/blob/9a402e9e6a04b45dd788c7dfb0106dae45443b01/profiles/apparmor.d/tunables/etc#L25 "View source"){ .abs-source }

@{etc_rw} contains directories where writing to configuration files is allowed.
@{etc_rw} should always be a subset of @{etc_ro}.

Only use @{etc_rw} if the profile allows writing to a configuration file.
For rules that only allows read access, use @{etc_ro}.

```
@{etc_rw}=/etc/
```

## kernelvars

This file should contain declarations to kernel vars or variables
that will become kernel vars at some point

### @{pid}

[:material-file-eye:](https://gitlab.com/apparmor/apparmor/-/blob/9a402e9e6a04b45dd788c7dfb0106dae45443b01/profiles/apparmor.d/tunables/kernelvars#L16 "View source"){ .abs-source }

until kernel vars are implemented
and until the parser supports nested groupings like
use

```
@{pid}={[1-9],[1-9][0-9],[1-9][0-9][0-9],[1-9][0-9][0-9][0-9],[1-9][0-9][0-9][0-9][0-9],[1-9][0-9][0-9][0-9][0-9][0-9],[1-4][0-9][0-9][0-9][0-9][0-9][0-9]}
```
### @{tid}

[:material-file-eye:](https://gitlab.com/apparmor/apparmor/-/blob/9a402e9e6a04b45dd788c7dfb0106dae45443b01/profiles/apparmor.d/tunables/kernelvars#L19 "View source"){ .abs-source }

same pattern as @{pid} for now

```
@{tid}=@{pid}
```
### @{pids}

[:material-file-eye:](https://gitlab.com/apparmor/apparmor/-/blob/9a402e9e6a04b45dd788c7dfb0106dae45443b01/profiles/apparmor.d/tunables/kernelvars#L22 "View source"){ .abs-source }

A pattern for pids that can appear

```
@{pids}=@{pid}
```
### @{uid}

[:material-file-eye:](https://gitlab.com/apparmor/apparmor/-/blob/9a402e9e6a04b45dd788c7dfb0106dae45443b01/profiles/apparmor.d/tunables/kernelvars#L27 "View source"){ .abs-source }

Placeholder for user id until kernel var is implemented to match
current user of the confined application.
Values are 0...4,294,967,295 (32-bit unsigned, 10 digits).

```
@{uid}={[0-9],[1-9][0-9],[1-9][0-9][0-9],[1-9][0-9][0-9][0-9],[1-9][0-9][0-9][0-9][0-9],[1-9][0-9][0-9][0-9][0-9][0-9],[1-9][0-9][0-9][0-9][0-9][0-9][0-9],[1-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9],[1-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9],[1-4][0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9]}
```
### @{uids}

[:material-file-eye:](https://gitlab.com/apparmor/apparmor/-/blob/9a402e9e6a04b45dd788c7dfb0106dae45443b01/profiles/apparmor.d/tunables/kernelvars#L30 "View source"){ .abs-source }

same pattern as @{uid} for now

```
@{uids}=@{uid}
```
### @{sys}

[:material-file-eye:](https://gitlab.com/apparmor/apparmor/-/blob/9a402e9e6a04b45dd788c7dfb0106dae45443b01/profiles/apparmor.d/tunables/kernelvars#L33 "View source"){ .abs-source }

until kernel var is implemented

```
@{sys}=/sys/
```

## multiarch

### @{multiarch}

[:material-file-eye:](https://gitlab.com/apparmor/apparmor/-/blob/9a402e9e6a04b45dd788c7dfb0106dae45443b01/profiles/apparmor.d/tunables/multiarch#L13 "View source"){ .abs-source }

@{multiarch} is the set of patterns matching multi-arch library
install prefixes.

```
@{multiarch}=*-linux-gnu*
```

## proc

### @{PROC}

[:material-file-eye:](https://gitlab.com/apparmor/apparmor/-/blob/9a402e9e6a04b45dd788c7dfb0106dae45443b01/profiles/apparmor.d/tunables/proc#L12 "View source"){ .abs-source }

@{PROC} is the location where procfs is mounted.

```
@{PROC}=/proc/
```

## run

### @{run}

[:material-file-eye:](https://gitlab.com/apparmor/apparmor/-/blob/9a402e9e6a04b45dd788c7dfb0106dae45443b01/profiles/apparmor.d/tunables/run#L1 "View source"){ .abs-source }

```
@{run}=/run/ /var/run/
```

## securityfs

### @{securityfs}

[:material-file-eye:](https://gitlab.com/apparmor/apparmor/-/blob/9a402e9e6a04b45dd788c7dfb0106dae45443b01/profiles/apparmor.d/tunables/securityfs#L10 "View source"){ .abs-source }

@{securityfs} is the location where securityfs is mounted.

```
@{securityfs}=@{sys}/kernel/security/
```
