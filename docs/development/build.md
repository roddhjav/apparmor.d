---
title: Building the profiles
---

The profiles in `apparmor.d` must not be used directly. They need to be prebuilt (by running `just complain`). This page documents all possibles prebuild tasks. It is not intended to be read by end user, and it is only targeted at developers and maintainers.

The build system is fully configurable, general usage can be seen with:
```sh
go run ./cmd/prebuild -h
```

```
aa-prebuild [-h] [--complain | --enforce] [--full] [--server] [--abi 3|4] [--version V] [--file FILE]

    Prebuild apparmor.d profiles for a given distribution and apply
    internal built-in directives.

Options:
    -h, --help        Show this help message and exit.
    -c, --complain    Set complain flag on all profiles.
    -e, --enforce     Set enforce flag on all profiles.
    -a, --abi ABI     Target apparmor ABI.
    -v, --version V   Target apparmor version.
    -f, --full        Set AppArmor for full system policy.
    -s, --server      Set AppArmor for server.
    -b, --buildir DIR Root build directory.
    -F, --file        Only prebuild a given file.
        --debug       Enable debug mode.

Prepare tasks:
    configure - Set distribution specificities
    setflags - Set flags on some profiles
    fsp - Configure AppArmor for full system policy
    merge - Merge profiles (from group/, profiles-*-*/) to a unified apparmor.d directory
    overwrite - Overwrite dummy upstream profiles
    synchronise - Initialize a new clean apparmor.d build directory
    ignore - Ignore profiles and files from:
    server - Configure AppArmor for server
    systemd-default - Configure systemd unit drop in files to a profile for some units
    systemd-early - Configure systemd unit drop in files to ensure some service start after apparmor
    attach - Configure tunable for re-attached path

Build tasks:
    userspace - Fix: resolve variable in profile attachments
    abi3 - Build: convert all profiles from abi 4.0 to abi 3.0
    attach - Feat: re-attach disconnected path
    base-strict - Feat: use 'base-strict' as base abstraction
    complain - Build: set complain flag on all profiles
    debug - Build: debug mode enabled
    enforce - Build: all profiles have been enforced
    fsp - Feat: prevent unconfined transitions in profile rules
    hotfix - Fix: temporary solution for #74, #80 & #235
    stacked-dbus - Fix: resolve peer label variable in dbus rules

Directive:
    #aa:dbus own bus=<bus> name=<name> [interface=AARE] [path=AARE]
    #aa:dbus talk bus=<bus> name=<name> label=<profile> [interface=AARE] [path=AARE]
    #aa:dbus common bus=<bus> name=<name> label=<profile>
    #aa:exec [P|U|p|u|PU|pu|] profiles...
    #aa:only filters...
    #aa:exclude filters...
    #aa:stack [X] profiles...
```

## Prepare Tasks

### **`synchronise`**

Initialize a new clean `apparmor.d` build directory in `.build/`.

*Enabled by default. Can be disabled in `cmd/prebuild/main.go`*

### **`ignore`**

Ignore profiles and files as defined in the `dist/ignore` directory. See [workflow](workflow.md#ignore-profiles).

*Enabled by default. Can be disabled in `cmd/prebuild/main.go`*

### **`server`**

Configure AppArmor for server. Desktop related groups and profiles that use desktop abstraction are not included. [hotfix](#hotfix) is also disabled, as it is only needed on desktop system. It is mostly intended to be used on server with FSP enabled. E.g: [the play machine](https://github.com/roddhjav/play).

*Enable with the `--server` option in the prebuild command.*

### **`merge`**

Merge profiles from `apparmor.d/group/`, `apparmor.d/profiles-*-*/` to a unified directory in `.build/apparmor.d` that AppArmor can parse.

*Enabled by default. Can be disabled in `cmd/prebuild/main.go`*

### **`configure`**

Set distribution specificities as defined in [`pkg/prebuild/prepare/configure.go`](https://github.com/roddhjav/apparmor.d/blob/main/pkg/prebuild/prepare/configure.go)

*Enabled by default. Can be disabled in `cmd/prebuild/main.go`*

### **`setflags`**

Set flags on profiles as defined in the [flags manifest](workflow.md#profile-flags).

*Enabled by default. Can be disabled in `cmd/prebuild/main.go`*

### **`overwrite`**

Overwrite (dummy) upstream profiles as defined in `dist/overwrite`.

*Enabled by default. Can be disabled in `cmd/prebuild/main.go`*

### **`systemd-default`**

Install systemd unit drop in files from `systemd/default`. They configure the various dbus daemon to use specific profiles.

*Enabled by default. Can be disabled in `cmd/prebuild/main.go`*

### **`systemd-early`**

Install systemd unit drop in files from `systemd/early` to ensure some services start after AppArmor. THis task will be removed in the future, as it will not be needed any more.

*Enabled by default. Can be disabled in `pkg/prebuild/cli/cli.go`*

### **`fsp`**

Configure AppArmor for full system policy.

*Enable with the `--full` option in the prebuild command.*


## Build Tasks

### **`abi3`**

This task will convert all profiles from `abi/4.0` to `abi/3.0`. The rules not supported by `abi/3.0` are commented in the build profiles.

*Enable with the `--abi 3` option in the prebuild command.*

### **`complain | enforce`**

Set or remove the complain flag on all profiles. The `complain` task is enabled by default. When building in enforce mode, it is disabled. Enabling the `enforce` task will enforce **all** profiles including the one set in the [flags manifest](workflow.md#profile-flags). It is intended to be used in specialized system such as a CTF challenge or in (very) high security VM. 

*Enable with the `--complain` or `--enforce` option in the prebuild command.*

### **`userspace`**

Resolve variables in profile attachments. It fixes issues with the userland AppArmor tools (aa-enforce, aa-logprof...) that do not support identical variable in the profiles attachments.

*Enabled by default. Can be disabled in `cmd/prebuild/main.go`*

### **`attach`**

This task reattaches disconnected paths. See the [Re-attached path](internal.md#re-attached-path) page. It will:

- Add the `attach_disconnected.path` flag on all profiles with the `attach_disconnected` flag
- Add the `<abstractions/attached/base>` abstraction in the profile
- For compatibility, non-disconnected profile will have the `@{att}` variable set to `/`

*Enabled when abi >= 4.0*

### **`hotfix`**

Temporary fix for #74, #80 & #235. Only an issue on Gnome, can be disabled on server.

*Enabled by default. Can be disabled in `cmd/prebuild/main.go`*

### **`fsp`**

Prevent unconfined transitions in profile rules.

*Enable with the `--full` option in the prebuild command.*
