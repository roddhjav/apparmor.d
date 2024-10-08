---
title: Known issues
---

!!! info

    Known bugs are tracked on the meta issue **[#75](https://github.com/roddhjav/apparmor.d/issues/74)**.

## Complain mode

A profile in *complain* mode cannot break the program it confines. However, there are some **major exceptions**:

1. `deny` rules are enforced even in *complain* mode,
2. `attach_disconnected` (and `mediate_deleted`) will break the program if they are required and missing in the profile,
3. If AppArmor does not find the profile to transition `rPx`.

## Pacman "could not get current working directory"

```sh
$ sudo pacman -Syu
...
error: could not get current working directory
:: Processing package changes...
...
```

This is **a feature, not a bug!** It can safely be ignored. Pacman tries to get your current directory. You will only get this error when you run pacman in your home directory.

According to the Arch Linux guideline, on Arch Linux, packages cannot install files under `/home/`. Therefore, the [`pacman`][pacman] profile purposely does not allow access of your home directory.

This provides a basic protection against some packages (on the AUR) that may have rogue install script.

[pacman]: https://github.com/roddhjav/apparmor.d/blob/main/apparmor.d/groups/pacman/pacman
