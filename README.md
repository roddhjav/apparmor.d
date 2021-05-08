[<img src="https://gitlab.com/uploads/-/system/project/avatar/25600351/logo.png" align="right" height="110"/>][project]

# apparmor.d

[![][build]][project]

**Full set of apparmor profiles**

> Warning: This project is still in early development.


## Description 

A set of over 800 apparmor profiles which aims is to confine most of Linux base
applications and processes.

**Goals & Purpose**
- All distribution that support Apparmor (currenlty Archlinux and Debian),
- Target both desktop and server,
- Confine all root services (bluetooth, dbus, polkit, networkmanager...),
- Confine all Desktop environments (currently only Gnome),
- Fully tested (Work in progress),
- Should not break a normal usage of the confined software.

These profiles strive to be fully functional with zero audit log warnings under
proper behavior. Functionality is not ignored. If functionality is not
explicitly blocked, then it's probably a bug in the profile and should be fixed.

**Note:** This work is part of a bigger linux security project.

> This project is based on the excellent work from [Morfikov][upstream] and aims
to extend it to more Linux distributions and desktop environements.


## Tests

A full test suite to ensure compatibility across distributions and softwares is
still a work in progress.

## Installation

**Requirements**
* An `apparmor` based linux distribution.
* A `systemd` based linux distribution.
* Base profiles and abstraction shipped with apparmor are supposed to be installed.

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

This program is based on Mikhail Morfikov's [apparmor profiles project][upstream]
and thus has the same license (GPL2).

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
