% aa-mode(1)
% aa-mode was written by Alexandre Pujol (alexandre@pujol.io)
% April 2026

# NAME

aa-mode - Switch AppArmor profiles mode.

# SYNOPSIS

**aa-mode** [*options...*] (**-e**|**-c**|**-k**|**-a**|**-u**|**-p**) [*profiles...*]

# DESCRIPTION

Switch AppArmor profiles mode. It modifies the profile flags and reloads the profiles with `apparmor_parser(8)`.

If a profile name is given without a path, it is looked up in `/etc/apparmor.d/`. If a directory is given, all profiles in it are processed recursively.

Exactly one mode option must be given.

# OPTIONS

*profiles*

: One or more profile paths or names to switch.

`--enforce`, `-e`

: Set the profile in **enforce** mode.

`--complain`, `-c`

: Set the profile in **complain** mode.

`--kill`, `-k`

: Set the profile in **kill** mode.

`--default-allow`, `-a`

: Set the profile in **default_allow** mode.

`--unconfined`, `-u`

: Set the profile in **unconfined** mode.

`--prompt`, `-p`

: Set the profile in **prompt** mode.

`--no-reload`

: Do not reload the profile after modifying it.

`--help`, `-h`

: Print the program usage.

# USAGE

To switch a profile to complain mode:
```sh
aa-mode -c dnsmasq
```

To switch a profile to enforce mode:
```sh
aa-mode --enforce /etc/apparmor.d/dnsmasq
```

To switch all profiles in a directory to complain mode without reloading:
```sh
aa-mode --complain --no-reload /etc/apparmor.d/
```

# SEE ALSO

`apparmor_parser(8)`, `apparmor(7)`, `apparmor.d(5)`, `aa-log(1)`, `aa-enforce(1)`, `aa-complain(1)`, and
https://apparmor.pujol.io.
