[<img src="https://gitlab.com/uploads/-/system/project/avatar/25600351/logo.png" align="right" height="110"/>][project]

# apparmor.d

[![][workflow]][action] [![][build]][project] [![][quality]][goreportcard] [![][matrix]][matrix-link] 

**Full set of AppArmor profiles**

> [!WARNING]
> This project is still in its early development. Help is very welcome; see the [documentation website](https://apparmor.pujol.io/) including its [development](https://apparmor.pujol.io/development) section.


## Description 

**AppArmor.d** is a set of over 1500 AppArmor profiles whose aim is to confine most Linux based applications and processes.

**Purpose**

- Confine all root processes such as all `systemd` tools, `bluetooth`, `dbus`,
  `polkit`, `NetworkManager`, `OpenVPN`, `GDM`, `rtkit`, `colord`
- Confine all Desktop environments
- Confine all user services such as `Pipewire`, `Gvfsd`, `dbus`, `xdg`, `xwayland`
- Confine some *"special"* user applications: web browser, file browser...
- Should not break a normal usage of the confined software

**Goals**

- Target both desktops and servers
- Support all distributions that support AppArmor:
    * Archlinux
    * Ubuntu 22.04
    * Debian 12
    * OpenSUSE Tumbleweed
- Support major desktop environments:
    * Gnome
    * KDE *(work in progress)*
- Fully tested (Work in progress)


> This project is originally based on the work from [Morfikov][upstream] and aims to extend it to more Linux distributions and desktop environments.

## Concepts

*One profile a day keeps the hacker away*

There are over 50000 Linux packages and even more applications. It is simply not possible to write an AppArmor profile for all of them. Therefore, a question arises:

**What to confine and why?**

We take inspiration from the [Android/ChromeOS Security Model][android_model], and we apply it to the Linux world. Modern [Linux security distributions][clipos] usually consider an immutable core base image with a carefully selected set of applications. Everything else should be sandboxed. Therefore, this project tries to confine all the *core* applications you will usually find in a Linux system: all systemd services, xwayland, network, Bluetooth, your desktop environment... Non-core user applications are out of scope as they should be sandboxed using a dedicated tool (minijail, bubblewrap, toolbox...).

This is fundamentally different from how AppArmor is usually used on Linux servers as it is common to only confine the applications that face the internet and/or the users.

**Presentations**

Building large set of AppArmor profiles:

- [Linux Security Summit North America (LSS-NA 2023)](https://events.linuxfoundation.org/linux-security-summit-north-america/) *([Slide](https://lssna2023.sched.com/event/1K7bI/building-the-largest-working-set-of-apparmor-profiles-alexandre-pujol-the-collaboratory-tudublin), [Video](https://www.youtube.com/watch?v=OzyalrOzxE8))*
- [Ubuntu Summit 2023](https://events.canonical.com/event/31/) *([Slide](https://events.canonical.com/event/31/contributions/209/), [Video](https://www.youtube.com/watch?v=GK1J0TlxnFI))*

## Installation

Please see [apparmor.pujol.io/install](https://apparmor.pujol.io/install)

## Configuration

Please see [apparmor.pujol.io/configuration](https://apparmor.pujol.io/configuration)

## Usage

Please see [apparmor.pujol.io/usage](https://apparmor.pujol.io/usage)

## Contribution

Feedbacks, contributors, pull requests are all very welcome. Please read
[apparmor.pujol.io/development](https://apparmor.pujol.io/development) 
for more details on the contribution process.

Development chat available on https://matrix.to/#/#apparmor.d:matrix.org

## License

This Project was initially based on Mikhail Morfikov's [apparmor profiles project][upstream]
and thus has the same license (GPL2).

[upstream]: https://gitlab.com/morfikov/apparmemall
[project]: https://gitlab.com/roddhjav/apparmor.d
[build]: https://gitlab.com/roddhjav/apparmor.d/badges/main/pipeline.svg?style=flat-square
[workflow]: https://img.shields.io/endpoint.svg?url=https%3A%2F%2Factions-badge.atrox.dev%2Froddhjav%2Fapparmor.d%2Fbadge%3Fref%3Dmain&style=flat-square
[action]: https://actions-badge.atrox.dev/roddhjav/apparmor.d/goto?ref=main
[quality]: https://img.shields.io/badge/go%20report-A+-brightgreen.svg?style=flat-square
[goreportcard]: https://goreportcard.com/report/github.com/roddhjav/apparmor.d
[matrix]: https://img.shields.io/badge/Matrix-%23apparmor.d-blue?style=flat-square&logo=matrix
[matrix-link]: https://matrix.to/#/#apparmor.d:matrix.org

[android_model]: https://arxiv.org/pdf/1904.05572
[clipos]: https://clip-os.org/en/
[write xor execute]: https://en.wikipedia.org/wiki/W%5EX
