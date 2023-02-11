---
title: Known issues
---

# Known issues

!!! info 

    Known bugs are tracked on the meta issue **[#75](https://github.com/roddhjav/apparmor.d/issues/74)**.


### Pacman "could not get current working directory"

```sh
$ sudo pacman -Syu
...
error: could not get current working directory
:: Processing package changes...
...
```

This is **a feature, not a bug!** It can safely be ignored. Pacman tries to get
your current directory. You will only get this error when you run pacman in your
home directory.

According the Archlinux guideline, on Archlinux, packages cannot install files
under `/home/`. Therefore the [`pacman`][pacman] profile purposely does not
allow access of your home directory.

This provides a basic protection against some packages (on the AUR) that may have
rogue install script.

[pacman]: https://github.com/roddhjav/apparmor.d/blob/main/apparmor.d/groups/pacman/pacman


### Gnome can be very slow to start.

[Gnome](https://github.com/roddhjav/apparmor.d/issues/80) can be slow to start.
This is a known bug, help is very welcome.

The complexity is that:

- It works fine without AppArmor
- It works fine on most system (including test VM)
- It seems to be dbus related
- On archlinux, the dbus mediation is not enabled. So, there is nothing special to allow.
