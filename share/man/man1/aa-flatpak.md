% aa-flatpak(1)
% aa-flatpak was written by Alexandre Pujol (alexandre@pujol.io)
% April 2026

# NAME

aa-flatpak - Confine Flatpak applications with AppArmor profiles.

# SYNOPSIS

**aa-flatpak** [*options...*] [*application-id...*]

# DESCRIPTION

**aa-flatpak** generates AppArmor profiles for installed Flatpak applications from their metadata and permissions, tailoring each profile to the application and integrating it with apparmor.d.

Profiles are written to `/etc/apparmor.d/flatpak-apps/` and reloaded with `apparmor_parser(8)`. Flatpak metadata is read from `/var/lib/flatpak/app/<id>/current/active/metadata`, and per-application overrides from `~/.local/share/flatpak/overrides/` are applied.

If one or more *application-id* arguments are given, profiles are generated only for those applications. If no ID is provided, profiles are generated for all installed Flatpak applications.

On Ubuntu, profiles are generated without an attachment path.

# OPTIONS

**aa-flatpak** [*options...*] [*application-id...*]

[*application-id*]

: Optional Flatpak application ID (e.g. `org.mozilla.firefox`). Multiple IDs may be given. If omitted, all installed applications are processed.

`--daemon`, `-d`

: Run as a daemon. Generate profiles for all installed applications, then watch `/var/lib/flatpak/app/` and the user overrides directory for changes and regenerate, reload, or remove profiles accordingly.

`--help`, `-h`

: Print the program usage.

# USAGE

To generate and load profiles for all installed Flatpak applications:
```sh
aa-flatpak
```

To generate a profile for a single application:
```sh
aa-flatpak org.mozilla.firefox
```

To run as a daemon and keep profiles in sync with installed Flatpaks:
```sh
aa-flatpak --daemon
```

# FILES

`/etc/apparmor.d/flatpak-apps/`

: Directory where generated profiles are written. Profile files are named `flatpak.<application-id>`.

`/var/lib/flatpak/app/<id>/current/active/metadata`

: Flatpak application metadata used as input.

`~/.local/share/flatpak/overrides/`

: Per-application Flatpak overrides applied on top of the metadata.

# SEE ALSO

`apparmor_parser(8)`, `apparmor(7)`, `apparmor.d(5)`, `aa-log(1)`, `aa-enforce(1)`, `aa-complain(1)`, `flatpak(1)`, and
https://apparmor.pujol.io/flatpak/.
