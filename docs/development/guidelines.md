---
title: Guidelines
---

## Common structure

AppArmor profiles can be written without any specific guidelines. However, when you work with over 1400 profiles, you need a common structure among all the profiles. 

The logic behind it is that if a rule is present in a profile, it should only be in one place, making profile review easier. 

For example, if a program needs to run executables binary. The rules allowing it can only be in a specific rule block (just after the `@{exec_path} mr,` rule). It is therefore easy to ensure some profile features such as:

* A profile has access to a given resource 
* A profile enforces a strict [write xor execute] (W^X) policy. 

It also improves compatibilities and makes personalization easier thanks to the use of more variables.

## Guidelines

!!! note

    This profile guideline is still evolving, feel free to propose improvements
    as long as they do not vary too much from the existing rules.

In order to ensure a common structure across the profiles, all new profile **must** follow the guidelines presented here.

The rules in the profile should be sorted in the rule ***block*** as follows:

| Order | Name | Example |
|:-----:|:----:|:-------:|
| **1** | [`include`](https://man.archlinux.org/man/apparmor.d.5##include_mechanism) | [:octicons-link-external-24:](https://github.com/search?q=repo%3Aroddhjav%2Fapparmor.d+NOT+path%3A*.md+include+%3Cabstractions%2F&type=code) |
| **2** | [`set rlimit`](https://man.archlinux.org/man/apparmor.d.5#rlimit_rules) | [:octicons-link-external-24:](https://github.com/search?q=repo%3Aroddhjav%2Fapparmor.d+NOT+path%3A*.md+set+rlimit&type=code) |
| **3** | [`userns`](https://gitlab.com/apparmor/apparmor/-/wikis/unprivileged_userns_restriction) | [:octicons-link-external-24:](https://github.com/search?q=repo%3Aroddhjav%2Fapparmor.d+NOT+path%3A*.md+userns&type=code) |
| **4** | [`capability`](https://gitlab.com/apparmor/apparmor/-/wikis/AppArmor_Core_Policy_Reference#capability-rules) | [:octicons-link-external-24:](https://github.com/search?q=repo%3Aroddhjav%2Fapparmor.d+NOT+path%3A*.md+capability&type=code) |
| **5** | [`network`](https://gitlab.com/apparmor/apparmor/-/wikis/AppArmor_Core_Policy_Reference#network-rules) | [:octicons-link-external-24:](https://github.com/search?q=repo%3Aroddhjav%2Fapparmor.d+NOT+path%3A*.md+%22+network+%22&type=code) |
| **6** | [`mount`](https://gitlab.com/apparmor/apparmor/-/wikis/AppArmor_Core_Policy_Reference#mount-rules-apparmor-28-and-later) | [:octicons-link-external-24:](https://github.com/search?q=repo%3Aroddhjav%2Fapparmor.d+NOT+path%3A*.md+%22++mount+%22&type=code) |
| **7** | [`remount`](https://gitlab.com/apparmor/apparmor/-/wikis/AppArmor_Core_Policy_Reference#remount) | [:octicons-link-external-24:](https://github.com/search?q=repo%3Aroddhjav%2Fapparmor.d+NOT+path%3A*.md+remount&type=code) |
| **8** | [`umount`](https://gitlab.com/apparmor/apparmor/-/wikis/AppArmor_Core_Policy_Reference#umount)| [:octicons-link-external-24:](https://github.com/search?q=repo%3Aroddhjav%2Fapparmor.d+NOT+path%3A*.md+%22umount+%22&type=code) |
| **9** | [`pivot_root`](https://gitlab.com/apparmor/apparmor/-/wikis/AppArmor_Core_Policy_Reference#pivot_root)| [:octicons-link-external-24:](https://github.com/search?q=repo%3Aroddhjav%2Fapparmor.d+NOT+path%3A*.md+pivot_root&type=code) |
| **10** | [`change_profile`](https://gitlab.com/apparmor/apparmor/-/wikis/AppArmor_Core_Policy_Reference#change_profile)| [:octicons-link-external-24:](https://github.com/search?q=repo%3Aroddhjav%2Fapparmor.d+NOT+path%3A*.md+change_profile+&type=code) |
| **11** | `mqueue` | [:octicons-link-external-24:]() |
| **12** | [`signal`](https://gitlab.com/apparmor/apparmor/-/wikis/AppArmor_Core_Policy_Reference#signals)| [:octicons-link-external-24:](https://github.com/search?q=repo%3Aroddhjav%2Fapparmor.d+NOT+path%3A*.md+%22signal+%22&type=code) |
| **13** | [`ptrace`](https://man.archlinux.org/man/apparmor.d.5#PTrace_rules) | [:octicons-link-external-24:](https://github.com/search?q=repo%3Aroddhjav%2Fapparmor.d+NOT+path%3A*.md+%22ptrace+%22&type=code) |
| **14** | [`unix`](https://man.archlinux.org/man/apparmor.d.5#Unix_socket_rules) | [:octicons-link-external-24:](https://github.com/search?q=repo%3Aroddhjav%2Fapparmor.d+NOT+path%3A*.md+%22unix+%22&type=code) |
| **15** | `io_uring` | [:octicons-link-external-24:]() |
| **16** | [`dbus`](https://gitlab.com/apparmor/apparmor/-/wikis/AppArmor_Core_Policy_Reference#dbus-rules) | [:octicons-link-external-24:](https://github.com/search?q=repo%3Aroddhjav%2Fapparmor.d+NOT+path%3A*.md++NOT+path%3A*.go+%22+dbus+%22&type=code) |
| **17** | [`file`](https://gitlab.com/apparmor/apparmor/-/wikis/AppArmor_Core_Policy_Reference#file-access-rules) | [:octicons-link-external-24:](https://github.com/roddhjav/apparmor.d/blob/2e4788c51ef73798c0ac94993af3cd769723e8e4/apparmor.d/groups/gnome/gnome-shell#L481-L663) |
| **18** | local include | [:octicons-link-external-24:](https://github.com/search?q=repo%3Aroddhjav%2Fapparmor.d+NOT+path%3A*.md+include+if+exists+%3Clocal&type=code) |


This rule order is taken from AppArmor with minor changes as we tend to:

- Divide the file block in multiple subcategories
- Put the block with the longer rules (`files`, `dbus`) after the other blocks

### The file block

The file block should be sorted as follows:

| Order | Description | Example | Link |
|:-----:|:-----------:|:-------:|:------:|
| **1** | The entry point of the profile | `@{exec_path} mr,` | [:octicons-link-external-24:](https://github.com/roddhjav/apparmor.d/blob/2e4788c51ef73798c0ac94993af3cd769723e8e4/apparmor.d/groups/gnome/gdm#L67) |
| **2** | The binaries and library required | `@{bin}/`, `@{lib}/`, `/opt/`. It is the only place where you can have `mr`, `rix`, `rPx`, `rUx`, `rPUX` rules. | [:octicons-link-external-24:](https://github.com/roddhjav/apparmor.d/blob/2e4788c51ef73798c0ac94993af3cd769723e8e4/apparmor.d/groups/gnome/gdm#L69-L76) | 
| **3** | The shared resources | `/usr/share` | [:octicons-link-external-24:](https://github.com/roddhjav/apparmor.d/blob/2e4788c51ef73798c0ac94993af3cd769723e8e4/apparmor.d/groups/network/NetworkManager#L111-L120) | 
| **4** | The system configuration | `/etc` | [:octicons-link-external-24:](https://github.com/roddhjav/apparmor.d/blob/2e4788c51ef73798c0ac94993af3cd769723e8e4/apparmor.d/groups/network/NetworkManager#L111-L120) | 
| **5** | The system data | `/`, `/var`, `/boot` | [:octicons-link-external-24:](https://github.com/roddhjav/apparmor.d/blob/2e4788c51ef73798c0ac94993af3cd769723e8e4/apparmor.d/groups/gnome/tracker-extract#L83-L93) | 
| **6** | The user data | `owner @{HOME}/` | [:octicons-link-external-24:](https://github.com/roddhjav/apparmor.d/blob/2e4788c51ef73798c0ac94993af3cd769723e8e4/apparmor.d/groups/gnome/tracker-extract#L96-L98) | 
| **7** | The user configuration, cache and dotfiles | `@{user_cache_dirs}`, `@{user_config_dirs}`,  `@{user_share_dirs}` | [:octicons-link-external-24:](https://github.com/roddhjav/apparmor.d/blob/2e4788c51ef73798c0ac94993af3cd769723e8e4/apparmor.d/groups/browsers/firefox#L179-L202) | 
| **8** | Temporary and runtime data | `/tmp/`, `@{run}/`, `/dev/shm/` | [:octicons-link-external-24:]() | 
| **9** | Sys files | `@{sys}/` | [:octicons-link-external-24:]() | 
| **10** | Proc files | `@{PROC}/` | [:octicons-link-external-24:]() | 
| **11** | Dev files | `/dev/` | [:octicons-link-external-24:]() | 
| **12** | Deny rules | `deny` | [:octicons-link-external-24:]() | 

### The dbus block


The dbus block should be sorted as follows:

- The system bus should be sorted *before* the session bus
- The bind rules should be sorted *after* the send & receive rules

For DBus, try to determine peer's label when possible. E.g.:
```
dbus send bus=session path=/org/freedesktop/DBus
     interface=org.freedesktop.DBus
     member={RequestName,ReleaseName}
     peer=(name=org.freedesktop.DBus, label=dbus-session),
```
If there is no predictable label it can be omitted.

### Profile rules

`bin, sbin & lib`

:   - Do not use: `/usr/lib` or `/usr/bin` but `@{bin}/` or `@{lib}/`
    - Do not use: `/usr/sbin` or `/sbin` but `@{bin}/`.

`Variables`

:   Always use the apparmor variables.

`Sort`

:   In a rule block, the rules must be alphabetically sorted.

`Sub profile`

:   Sub profile should come at the end of a profile.

`Similar purpose`

:   When some rules share similar purpose, they may be sorted together. Eg:
    ```
    /etc/machine-id r,
    /var/lib/dbus/machine-id r,
    ```


## Additional recommended documentation

* [The AppArmor Core Policy Reference](https://gitlab.com/apparmor/apparmor/-/wikis/AppArmor_Core_Policy_Reference)
* [The OpenSUSE Documentation](https://doc.opensuse.org/documentation/leap/security/html/book-security/part-apparmor.html)
* https://documentation.suse.com/sles/12-SP5/html/SLES-all/cha-apparmor-intro.html
* [The AppArmor.d man page](https://man.archlinux.org/man/apparmor.d.5)
* [F**k AppArmor](https://presentations.nordisch.org/apparmor/#/)
* [A Brief Tour of Linux Security Modules](https://www.starlab.io/blog/a-brief-tour-of-linux-security-modules)

[write xor execute]: https://en.wikipedia.org/wiki/W%5EX
