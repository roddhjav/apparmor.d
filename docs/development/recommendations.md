---
title: Recommendations
---

## Renaming of profiles

For security reason, once loaded into the kernel, a profile cannot get fully removed. Therefore, by renaming a profile, you create a second profile with the same attachment. AppArmor will not be able to determine witch one to use leading to breakage.

A reboot is required to fully remove the profile from the kernel.


## Programs to not confine

Some programs should not be confined by themselves. For example, tools such as `ls`, `rm`, `diff` or `cat` do not have profiles in this project. Let's see why.

These are general tools that in a general context can legitimately access any file in the system. Therefore, the confinement of such tools by a global profile would at best be minimal at worst be a security theatre.

It gets even worse. Let's say, we write a profile for `cat`. Such a profile would need access to `/etc/`. We will add the following rule:
```sh
  /etc/{,**} rw,
```

However, as `/etc` can contain sensitive files, we now want to explicitly prevent access to these sensitive files. Problems:

1. How do we know the exhaustive list of *sensitive files* in `/etc`?
2. How do we ensure access to these sensitive files is not required?
3. This breaks the principle of mandatory access control.
   See the [first rule of this project](index.md#project-rules) which is to only allow
   what is required. Here we allow everything and blacklist some paths.

It creates even more issues when we want to use this profile in other profiles. Let's take the example of `diff`. Using this rule: `@{bin}/diff rPx,` this will restrict access to the very generic and not very confined `diff` profile. Whereas most of the time, we want to restrict `diff` to some specific file in our profile:

* In `dpkg`, an internal child profile (`rCx -> diff`), allows `diff` to only access etc config files:

!!! note ""

    [apparmor.d/apparmor.d/groups/apt/dpkg](https://github.com/roddhjav/apparmor.d/blob/accf5538bdfc1598f1cc1588a7118252884df50c/apparmor.d/groups/apt/dpkg#L123)
    ``` aa linenums="123"
    profile diff {
      include <abstractions/base>
      include <abstractions/consoles>

      @{bin}/       r,
      @{bin}/pager mr,
      @{bin}/less  mr,
      @{bin}/more  mr,
      @{bin}/diff  mr,

      owner @{HOME}/.lesshs* rw,

      # Diff changed config files
      /etc/** r,

      # For shell pwd
      /root/ r,

    }
    ```

* As it is a dependency of pass, `diff` inherits the `pass' profile and has the same access as the pass profile, so it will be allowed to diff password files because more than a generic `diff`, it is a `diff` "version" for the pass password manager:

!!! note ""

    [apparmor.d/apparmor.d/profiles-m-r/pass](https://github.com/roddhjav/apparmor.d/blob/accf5538bdfc1598f1cc1588a7118252884df50c/apparmor.d/profiles-m-r/pass#L20
    )
    ``` aa linenums="20"
      @{bin}/diff      rix,
    ```

**What if I still want to protect these programs?**

You do not protect these programs. *Protect the usage you have of these programs*. In practice, it means that you should put your terminal in a sandbox managed environment with a sandboxing tool such as Toolbox.

!!! example "To sum up"

    1. Do not create a profile for programs such as: `rm`, `ls`, `diff`, `cd`, `cat`
    2. Do not create a profile for the shell: `bash`, `sh`, `dash`, `zsh`
    3. Use [Toolbox](https://containertoolbx.org/)
