% aa-install(1)
% aa-install was written by Alexandre Pujol (alexandre@pujol.io)
% July 2026

# NAME

aa-install - Install and manage AppArmor profiles from apparmor.d.

# SYNOPSIS

**aa-install** [*options...*] [**-s**|**-l**|**-i**|**-u**] [**-a**] [**-e**|**-c**]

# DESCRIPTION

Install and manage the AppArmor profiles shipped by apparmor.d. Profiles are
built from a source directory, filtered and flagged according to the
configuration, then deployed to the AppArmor policy directory and apparmor is
reloaded.

Installed profiles are recorded in a manifest so that later runs can report the
installation status, list the deployed profiles, or uninstall them.

With no action flag, print the installation status summary.

# OPTIONS

`--status`, `-s`

: Show the installation status summary. This is the default action. It reports
  how many profiles are recorded in the manifest, and how many are up-to-date,
  drifted (modified since install), or missing.

`--list`, `-l`

: List the installed profile paths from the manifest, one per line, sorted.

`--install`, `-i`

: Build and install the profiles, then reload them.

`--all`, `-a`

: Install all the profiles, including the ones of programs not installed on
  the system.

`--complain`, `-c`

: Set the **complain** flag on all installed profiles. Combined with an action,
  it is the default deploy mode. Mutually exclusive with **--enforce**.

`--enforce`, `-e`

: Set the **enforce** flag on all installed profiles. Combined with an action,
  it is the default deploy mode. Mutually exclusive with **--complain**.

`--uninstall`, `-u`

: Remove all profiles recorded in the manifest, then remove the manifest.

`--no-reload`

: Do not reload the profiles after modifying them.

`--config DIR`

: Select an alternate configuration directory (default: `/etc/apparmor/`).

`--magic DIR`

: Select an alternate apparmor.d policy directory (default: `/etc/apparmor.d`).

`--src DIR`

: Select an alternate source directory (default: `/usr/share/apparmor.d`).

`--help`, `-h`

: Print the program usage.

# CONFIGURATION

Configuration is read from two tiers. Vendor defaults in `/usr/share/apparmor/`
are read first, then the local admin configuration in `/etc/apparmor/` (or the
directory given by **--config**), which overrides the vendor defaults.

Each tier may contain:

*modes*

: Modes of enforcement for all profiles. The `default` key sets the default
  deploy mode (**enforce** or **complain**). Overridden by **--enforce** or
  **--complain**. Defaults to **complain**.

*flags.d/\*.conf*

: Set per-profile flags.

*ignore.d/\*.conf*

: Set the (groups of) profiles to ignore.

*include.d/\*.conf*

: Only install the given (groups of) profiles. When non-empty, all other
  profiles are excluded.

# USAGE

To install all profiles in the default (complain) mode:
```sh
aa-install --install
```

To install and enforce all profiles:
```sh
aa-install --install --enforce
```

To show the installation status:
```sh
aa-install
```

To list the installed profiles:
```sh
aa-install --list
```

To remove all installed profiles:
```sh
aa-install --uninstall
```

# SEE ALSO

`apparmor_parser(8)`, `apparmor(7)`, `apparmor.d(5)`, `aa-log(1)`, `aa-mode(1)`, `aa-enforce(1)`, `aa-complain(1)`, and
https://apparmor.pujol.io.
