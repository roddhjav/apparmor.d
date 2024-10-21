---
title: AppArmor.d
---

<style>.md-typeset .md-content__button { display: none; }</style>

**Full set of AppArmor profiles**

!!! danger "Help Wanted"

    This project is still in its early development. Help is very welcome; see [Development](development/index.md)

**AppArmor.d** is a set of over 1500 AppArmor profiles whose aim is to confine most Linux based applications and processes.

### Purpose

- Confine all root processes such as all `systemd` tools, `bluetooth`, `dbus`, `polkit`, `NetworkManager`, `OpenVPN`, `GDM`, `rtkit`, `colord`
- Confine all Desktop environments
- Confine all user services such as `Pipewire`, `Gvfsd`, `dbus`, `xdg`, `xwayland`
- Confine some *"special"* user applications: web browsers, file managers, etc
- Should not break a normal usage of the confined software

See the [Concepts](concepts.md)' page for more detail on the architecture.

### Goals

- Target both desktops and servers
- Support for all distributions that support AppArmor:
    * [:material-arch: Arch Linux](install.md#archlinux)
    * [:material-ubuntu: Ubuntu 24.04/22.04](install.md#ubuntu)
    * [:material-debian: Debian 12](install.md#debian)
    * [:simple-suse: openSUSE Tumbleweed](install.md#opensuse)
- Support for all major desktop environments:
    - [x] :material-gnome: Gnome (GDM)
    - [x] :simple-kde: KDE (SDDM)
    - [ ] :simple-xfce: XFCE (Lightdm) *(work in progress)*
- Fully tested *(work in progress)*

### Presentations

Building the largest set of AppArmor profiles:

- [Linux Security Summit North America (LSS-NA 2023)](https://events.linuxfoundation.org/linux-security-summit-north-america/) *([Slide](https://lssna2023.sched.com/event/1K7bI/building-the-largest-working-set-of-apparmor-profiles-alexandre-pujol-the-collaboratory-tudublin), [Video](https://www.youtube.com/watch?v=OzyalrOzxE8))*
- [Ubuntu Summit 2023](https://events.canonical.com/event/31/) *([Slide](https://events.canonical.com/event/31/contributions/209/), [Video](https://www.youtube.com/watch?v=GK1J0TlxnFI))*

### Chat

A development chat is available on https://matrix.to/#/#apparmor.d:matrix.org
