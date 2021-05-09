[<img src="https://gitlab.com/uploads/-/system/project/avatar/25600351/logo.png" align="right" height="110"/>][project]

# apparmor.d [![][build]][project]

**Full set of AppArmor profiles**

> Warning: This project is still in early development.


## Description 

A set of over 800 AppArmor profiles which aims is to confine most of Linux base applications and processes.

**Goals & Purpose**
- Support all distribution that support AppArmor (currenlty Archlinux and Debian),
- Target both desktop and server,
- Confine all root processes (bluetooth, dbus, polkit, networkmanager, systemd...),
- Confine all Desktop environments (currently only Gnome),
- Should not break a normal usage of the confined software.
- Fully tested (Work in progress),

**Note:** This work is part of a bigger linux security project.

> This project is based on the excellent work from [Morfikov][upstream] and aims
to extend it to more Linux distributions and desktop environements.


## Concepts

There are over 50000 Linux packages and even more applications. It is simply not possible to write an AppArmor profile for all of them. Therefore a question arises: *What to confine and why?*

We take inspiration from the [Android/ChromeOS Security Model][android_model] and we apply it to the Linux world. Modern [linux security implementation][clipos] usually consider a core base image with a carefully set of selected applications. Everything else should be sandboxed. Therefore, this project tries to confine all the *core* applications you will usually find in a Linux system: all systemd services, xwayland, network, bluetooth, your desktop environment... Non-core user applications are out of scope as they should be sandboxed using a dedicated tool (minijail, bubblewrap...).

This is fundamentally different from how AppArmor is used on Linux server as it is common to only confine the applications that face the internet and/or the users.


## Tests

A full test suite to ensure compatibility across distributions and softwares is
still a work in progress.

## Installation

**Requirements**
* An `apparmor` based linux distribution.
* A `systemd` based linux distribution.
* Base profiles and abstractions shipped with AppArmor are supposed to be
  installed.

**Archlinux**

Build and install the package with:
```sh
makepkg -si
```

**Debian**

Build using standard Debian package build tools:
```sh
dpkg-buildpackage -b -d -us -ui --sign-key=<gpg-id>
```

## Contribution

Feedbacks, contributors, pull requests, are all very welcome.


## License

This program is based on Mikhail Morfikov's [apparmor profiles project][upstream] and thus has the same license (GPL2).

```
Copyright (C)  Alexandre PUJOL & Mikhail Morfikov

This program is free software; you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation; version 2 of the License.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License along
with this program; if not, write to the Free Software Foundation, Inc.,
51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.
```

[upstream]: https://gitlab.com/morfikov/apparmemall
[project]: https://gitlab.com/roddhjav/apparmor.d
[build]: https://gitlab.com/roddhjav/apparmor.d/badges/master/pipeline.svg?style=flat-square

[android_model]: https://arxiv.org/pdf/1904.05572
[clipos]: https://clip-os.org/en/
