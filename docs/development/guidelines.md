---
title: Guidelines
---

## Common structure

AppArmor profiles can be written without any specific guidelines. However, when you work with over 1500 profiles, you need a common structure among all the profiles.

The logic behind it is that if a rule is present in a profile, it should only be in one place, making it easier to review profiles.

For example, if a program needs to run executable binaries then the rules allowing it can only be in a specific rule block (just after the `@{exec_path} mr,` rule). It is therefore easy to ensure some profile features such as:

* A profile has access to a given resource 
* A profile enforces a strict [write xor execute](https://en.wikipedia.org/wiki/W%5EX) (W^X) policy. 

It also improves compatibility and makes personalization easier thanks to the use of more variables.

!!! danger "Why does it matter?"

    The structure presented below is more important than it looks. Without this it would be easy to miss some rules when reviewing profiles, making the whole system more opaque and less secure.

## Guidelines

In order to ensure a common structure across the profiles, all new profile **must** follow the guidelines presented here.

The rules in the profile should be written in *blocks* of same rule type. Each *block* should be sorted in the profile as follows:

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

- Divide the file block into multiple subcategories
- Put the block with the longer rules (`files`, `dbus`) after the other blocks

### File sub-blocks

The file block is usually the longest block in a profile. It should be written in sub-blocks of identical file types.
The file sub blocks **must** be sorted as follows.

| Order | Description | Example | Link |
|:-----:|:-----------:|:-------:|:------:|
| **1** | The entry point of the profile | `@{exec_path} mr,` | [:octicons-link-external-24:](https://github.com/roddhjav/apparmor.d/blob/2e4788c51ef73798c0ac94993af3cd769723e8e4/apparmor.d/groups/gnome/gdm#L67) |
| **2** | The binaries and library required | `@{bin}/`, `@{lib}/`, `/opt/`. It is the only place where you can have `mr`, `rix`, `rPx`, `rUx`, `rPUX` rules. | [:octicons-link-external-24:](https://github.com/roddhjav/apparmor.d/blob/2e4788c51ef73798c0ac94993af3cd769723e8e4/apparmor.d/groups/gnome/gdm#L69-L76) |
| **3** | The shared resources | `/usr/share` | [:octicons-link-external-24:](https://github.com/roddhjav/apparmor.d/blob/2e4788c51ef73798c0ac94993af3cd769723e8e4/apparmor.d/groups/network/NetworkManager#L111-L120) |
| **4** | The system configuration | `/etc` | [:octicons-link-external-24:](https://github.com/roddhjav/apparmor.d/blob/2e4788c51ef73798c0ac94993af3cd769723e8e4/apparmor.d/groups/network/NetworkManager#L111-L120) |
| **5** | The system data | `/`, `/var`, `/boot` | [:octicons-link-external-24:](https://github.com/roddhjav/apparmor.d/blob/2e4788c51ef73798c0ac94993af3cd769723e8e4/apparmor.d/groups/gnome/tracker-extract#L83-L93) |
| **6** | The user data | `owner @{HOME}/` | [:octicons-link-external-24:](https://github.com/roddhjav/apparmor.d/blob/2e4788c51ef73798c0ac94993af3cd769723e8e4/apparmor.d/groups/gnome/tracker-extract#L96-L98) |
| **7** | The user configuration, cache and dotfiles | `@{user_cache_dirs}`, `@{user_config_dirs}`, `@{user_share_dirs}` | [:octicons-link-external-24:](https://github.com/roddhjav/apparmor.d/blob/2e4788c51ef73798c0ac94993af3cd769723e8e4/apparmor.d/groups/browsers/firefox#L179-L202) |
| **8** | Temporary and runtime data | `/tmp/`, `@{run}/`, `/dev/shm/` | [:octicons-link-external-24:]() |
| **9** | Sys files | `@{sys}/` | [:octicons-link-external-24:]() |
| **10** | Proc files | `@{PROC}/` | [:octicons-link-external-24:]() |
| **11** | Dev files | `/dev/` | [:octicons-link-external-24:]() |
| **12** | Deny rules | `deny` | [:octicons-link-external-24:]() |


### Dbus block

Dbus rules must be written in a dedicated block just before the file block. The dbus block **must** be sorted as follows:

- The system bus should be sorted *before* the session bus
- The bind rules should be sorted *after* send & receive rules

For DBus, try to determine peer's label when possible. E.g.:
```
dbus send bus=session path=/org/freedesktop/DBus
     interface=org.freedesktop.DBus
     member={RequestName,ReleaseName}
     peer=(name=org.freedesktop.DBus, label="@{p_dbus_session}"),
```

If there is no predictable label it can be omitted.

### Profile rules

#### :material-numeric-1-circle: Variables

:   Always use the apparmor.d [variables](../variables.md).
    Example:

    - `/usr/lib` or `/usr/bin` become `@{bin}` or `@{lib}`
    - `/usr/sbin` or `/sbin` become `@{bin}`.

#### :material-numeric-2-circle: Sort

:   In a rule block, all rules must be alphabetically sorted.

#### :material-numeric-3-circle: Sub-profiles

:   Sub-profiles should come at the end of a profile.

#### :material-numeric-4-circle: Similar purpose

:   When some rules share similar purposes, they may be sorted together. E.g.:
    ```
    /etc/machine-id r,
    /var/lib/dbus/machine-id r,
    ```

#### :material-numeric-5-circle: Limit the use of `deny`

:   The use of `deny` should be limited to the minimum:

    - In MAC policies, we only allow access ([Rule :material-numeric-1-circle:](index.md#rule-mandatory-access-control "Mandatory Access Control"))
    - `deny` rules are enforced even in complain mode,
    - If it works on your machine does not mean it will work on others ([Rule :material-numeric-4-circle:](index.md#rule-distribution-and-devices-agnostic "Distribution and devices agnostic")).

#### :material-numeric-6-circle: Comments

:   Ensure you only have useful comments. E.g.:
    ```
    # Config files for foo
    owner @{user_config_dirs}/foo/{,**} r,
    ```
    Does not help, and if generalized it would add a lot of complexity to any profiles.

#### :material-numeric-7-circle: Clarity over cleverness

:   Always prefer clarity to cleverness. E.g., if a rule is more explicit but longer, prefer it over a shorter but less explicit one.
