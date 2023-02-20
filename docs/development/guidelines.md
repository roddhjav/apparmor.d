---
title: Guidelines
---

## Common structure

AppArmor profiles can be written without any specific guidelines. However,
when you work with over 1400 profiles, you need a common structure among all the
profiles. 

The logic behind it is that if a rule is present in a profile, it should only be
in one place, making profile review easier. 

For example, if a program needs to run executables binary. The rules allowing it
can only be in a specific rule block (just after the `@{exec_path} mr,` rule). It
is therefore easy to ensure some profile features such as:

* A profile has access to a given resource 
* A profile enforces a strict [write xor execute] (W^X) policy. 

It also improves compatibilities and makes personalization easier thanks to the
use of more variables.
 
## Guidelines

!!! note

    This profile guideline is still evolving, feel free to propose improvements
    as long as they do not vary too much from the existing rules.

In order to ensure a common structure across the profiles, all new profile **must**
follow the guidelines presented here.

The rules in the profile should be sorted in the rule ***block*** as follows:

1. [`include`](https://gitlab.com/apparmor/apparmor/-/wikis/AppArmor_Core_Policy_Reference#include-statements)
1. [`set rlimit`](https://gitlab.com/apparmor/apparmor/-/wikis/AppArmor_Core_Policy_Reference#rlimit-rules)
1. [`capability`](https://gitlab.com/apparmor/apparmor/-/wikis/AppArmor_Core_Policy_Reference#capability-rules)
1. [`network`](https://gitlab.com/apparmor/apparmor/-/wikis/AppArmor_Core_Policy_Reference#network-rules)
1. [`mount`](https://gitlab.com/apparmor/apparmor/-/wikis/AppArmor_Core_Policy_Reference#mount-rules-apparmor-28-and-later)
1. [`remount`](https://gitlab.com/apparmor/apparmor/-/wikis/AppArmor_Core_Policy_Reference#remount)
1. [`umount`](https://gitlab.com/apparmor/apparmor/-/wikis/AppArmor_Core_Policy_Reference#umount)
1. [`pivot_root`](https://gitlab.com/apparmor/apparmor/-/wikis/AppArmor_Core_Policy_Reference#pivot_root)
1. [`change_profile`](https://gitlab.com/apparmor/apparmor/-/wikis/AppArmor_Core_Policy_Reference#change_profile)
1. [`signal`](https://gitlab.com/apparmor/apparmor/-/wikis/AppArmor_Core_Policy_Reference#signals)
1. `ptrace`
1. `unix`
1. [`dbus`](https://gitlab.com/apparmor/apparmor/-/wikis/AppArmor_Core_Policy_Reference#dbus-rules)
1. [`file`](https://gitlab.com/apparmor/apparmor/-/wikis/AppArmor_Core_Policy_Reference#file-access-rules)
1. local include

This rule order is taken from AppArmor with minor changes as we tend to:

- Divide the file block in multiple subcategories
- Put the block with the longer rules (`files`, `dbus`) after the other blocks

### The file block

The file block should be sorted as follow:

1. `@{exec_path} mr`, the entry point of the profile
1. The binaries and library required:
    - `/{usr/,}bin/`, `/{usr/,}lib/`, `/opt/`...
    - It is the only place where you can have `mr`, `rix`, `rPx`, `rUx`, `rPUX` rules.
1. The shared resources: `/usr/share`...
1. The system configuration: `/etc`...
1. The system data: `/var`...
1. The user data: `owner @{HOME}/`...
1. The user configuration, cache and in general all dotfiles
1. Temporary and runtime data: `/tmp/`, `@{run}/`, `/dev/shm/`...
1. Sys files: `@{sys}/`...
1. Proc files: `@{PROC}/`... 
1. Dev files: `/dev/`...
1. Deny rules: `deny`...

### The dbus block


The dbus block should be sorted as follow:

- The system bus should be sorted *before* the session bus
- The bind rules should be sorted *after* the send & receive rules

For DBus, try to determine peer's label when possible. E.g.:
```
dbus send bus=session path=/org/freedesktop/DBus
     interface=org.freedesktop.DBus
     member={RequestName,ReleaseName}
     peer=(name=org.freedesktop.DBus, label=dbus-daemon),
```
If there is no predictable label it can be omitted.

### Profile rules

`bin, sbin & lib`

:   - Do not use: `/usr/lib` or `/usr/bin` but `/{usr/,}bin/` or `/{usr/,}lib/`
    - Do not use: `/usr/sbin` or `/sbin` but `/{usr/,}{s,}bin/`.

`Variables`

:   Always use the apparmor variables.

`Sort`

:   In a rule block, the rules must be alphabetically sorted.

`Sub profile`

:   Sub profile should comes at the end of a profile.

`Similar purpose`

:   When some rules share similar purpose, they may be sorted together. Eg:
    ```
    /etc/machine-id r,
    /var/lib/dbus/machine-id r,
    ```


## Additional recommended documentation

* [The AppArmor Core Policy Reference](https://gitlab.com/apparmor/apparmor/-/wikis/AppArmor_Core_Policy_Reference)
* [The AppArmor.d man page](https://man.archlinux.org/man/apparmor.d.5)
* [F**k AppArmor](https://presentations.nordisch.org/apparmor/#/)
* [A Brief Tour of Linux Security Modules](https://www.starlab.io/blog/a-brief-tour-of-linux-security-modules)

[write xor execute]: https://en.wikipedia.org/wiki/W%5EX
