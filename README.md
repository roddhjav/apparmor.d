[<img src="https://gitlab.com/uploads/-/system/project/avatar/25600351/logo.png" align="right" height="110"/>][project]

# apparmor.d

[![][workflow]][action] [![][build]][project]  [![][quality]][goreportcard]

**Full set of AppArmor profiles**

> Warning: This project is still in early development.


## Description 

A set of over 1000 AppArmor profiles which aims is to confine most of Linux base applications and processes.

**Goals & Purpose**
- Support all distributions that support AppArmor:
  * *Currenlty*: Archlinux, Debian 11 and the last Ubuntu LTS.
  * Not (yet) tested on openSUSE
- Target both desktop and server,
- Confine all root processes. Eg: all systemd tools, bluetooth, dbus, polkit,
  NetworkManager, OpenVPN, GDM, rtkit, colord...
- Confine all Desktop environments:
  * *Currently only Gnome*, see `apparmor.d/groups/gnome`
- Confine all user services: Eg: Pipewire, Gvfsd, dbus, xdg, xwayland...
- Confine some "special" user applications: web browser, file browser...
- Should not break a normal usage of the confined software.
- Fully tested (Work in progress),

> This project is based on the excellent work from [Morfikov][upstream] and aims
to extend it to more Linux distributions and desktop environements.


## Concepts

There are over 50000 Linux packages and even more applications. It is simply not possible to write an AppArmor profile for all of them. Therefore a question arises: *What to confine and why?*

We take inspiration from the [Android/ChromeOS Security Model][android_model] and we apply it to the Linux world. Modern [linux security implementation][clipos] usually consider a core base image with a carefully set of selected applications. Everything else should be sandboxed. Therefore, this project tries to confine all the *core* applications you will usually find in a Linux system: all systemd services, xwayland, network, bluetooth, your desktop environment... Non-core user applications are out of scope as they should be sandboxed using a dedicated tool (minijail, bubblewrap...).

This is fundamentally different from how AppArmor is used on Linux server as it is common to only confine the applications that face the internet and/or the users.



## Installation

**Requirements**
* An `apparmor` based linux distribution.
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
dpkg-buildpackage -b -d --no-sign
sudo dpkg --install ../apparmor.d_*_all.deb
```

**Partial install**

For test purpose, you can install a specific profile with the following commands. The tool will also install required abstractions and tunables:
```
sudo ./pick <profiles-name>
```


## Usage

**Enabled profiles**

Once installed and with the rules enabled, you can ensure the rules are loaded
with `sudo aa-satus`, it should give something like:
```
apparmor module is loaded.
1137 profiles are loaded.
794 profiles are in enforce mode.
   ...
343 profiles are in complain mode.
   ...
0 profiles are in kill mode.
0 profiles are in unconfined mode.
130 processes have profiles defined.
108 processes are in enforce mode.
   ...
22 processes are in complain mode.
   ...
0 processes are unconfined but have a profile defined.
0 processes are in mixed mode.
0 processes are in kill mode.
```

You can also list the current processes alongside with their security profile with
`ps auxZ`. Most of the process should then be confined.

**AppArmor Log**

The provided command `aa-log` allow you review AppArmor generated messages in a
colorfull way:

```
$ aa-log
   ...
```

`aa-log` can optionally be given a profile name as argument to
only shows the log for a given profile:
```
$ aa-log dnsmasq
DENIED  dnsmasq open /proc/sys/kernel/osrelease comm=dnsmasq requested_mask=r denied_mask=r
DENIED  dnsmasq open /proc/1/environ comm=dnsmasq requested_mask=r denied_mask=r
DENIED  dnsmasq open /proc/cmdline comm=dnsmasq requested_mask=r denied_mask=r
```


## Tests

A full test suite to ensure compatibility across distributions and softwares is still a work in progress.

Here an overview of the current CI jobs:

**On Gitlab CI**
- Package build for all supported distribution
- Profiles preprocessing verification for all supported distribution
- Go based command linting and unit tests

**On Github Action**
- Integration test on the ubuntu-latest VM: run a simple list of tasks with
  all the rules enabled and ensure no new issue has been raised. Github Action
  is used as it offers a direct access to a VM with AppArmor included.


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
[workflow]: https://img.shields.io/endpoint.svg?url=https%3A%2F%2Factions-badge.atrox.dev%2Froddhjav%2Fapparmor.d%2Fbadge&style=flat-square
[action]: https://actions-badge.atrox.dev/roddhjav/apparmor.d/goto
[quality]: https://img.shields.io/badge/go%20report-A+-brightgreen.svg?style=flat-square
[goreportcard]: https://goreportcard.com/report/github.com/roddhjav/apparmor.d

[android_model]: https://arxiv.org/pdf/1904.05572
[clipos]: https://clip-os.org/en/
