---
title: AppArmor.d
---

# AppArmor.d

**Full set of AppArmor profiles**

!!! danger "Help Wanted"

    This project is still in its early development. Help is very welcome; 
    see [Development](development/)

**AppArmor.d** is a set of over 1400 AppArmor profiles whose aim is to confine
most Linux based applications and processes.

**Purpose**

- Confine all root processes such as all `systemd` tools, `bluetooth`, `dbus`,
  `polkit`, `NetworkManager`, `OpenVPN`, `GDM`, `rtkit`, `colord`
- Confine all Desktop environments
- Confine all user services such as `Pipewire`, `Gvfsd`, `dbus`, `xdg`, `xwayland`
- Confine some *"special"* user applications: web browser, file browser...
- Should not break a normal usage of the confined software

See the [Concepts](concepts) page for more detail on the architecture.

**Goals**

- Target both desktops and servers
- Support all distributions that support AppArmor:
    * [:material-arch: Archlinux](/install/#archlinux)
    * [:material-ubuntu: Ubuntu 22.04](/install/#ubuntu-debian)
    * [:material-debian: Debian 11](/install/#ubuntu-debian)
    * [:simple-suse: OpenSUSE Tumbleweed](/install/#opensuse)
- Support all major desktop environments:
    * Currently only :material-gnome: Gnome
- Fully tested (Work in progress)
