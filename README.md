[<img src="https://gitlab.com/uploads/-/system/project/avatar/25600351/logo.png" align="right" height="110"/>][project]

# apparmor.d

[![][build]][project]

**Full set of apparmor profiles**

## Installation

**Requirements**
* An `apparmor` based linux distribution.
* A `systemd` based linux distribution.

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

[project]: https://gitlab.com/archlex/hardening/apparmor.d
[build]: https://gitlab.com/archlex/hardening/apparmor.d/badges/master/pipeline.svg?style=flat-square
