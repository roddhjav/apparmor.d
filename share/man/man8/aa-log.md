% aa-log(8)
% aa-log was written by Alexandre Pujol (alexandre@pujol.io)
% December 2025

# NAME

aa-log — Review AppArmor generated messages in a colorful way.

# SYNOPSIS

**aa-log** [*options…*] [*profile*]

# DESCRIPTION

Review AppArmor generated messages in a colourful way. Support logs from *auditd*, *systemd*, *syslog* as well as *dbus session* events.

It can be given an optional profile name to filter the output with.

It can be used to generate AppArmor rules from the logs and it therefore an alternative to `aa-logprof(8)`. The generated rules should be manually reviewed and inserted into the profile.

Default logs are read from `/var/log/audit/audit.log`. Other files in `/var/log/audit/` can easily be checked: **aa-log -f 1** parses `audit.log.1`

Use `aa-log -f -` to read from standard input.

Logs written with `aa-log` can be read again with `aa-log -l`.
  
# OPTIONS

**aa-log** [*options…*] [*profile*]

[*profile*]

: Optional profile name to filter the output with.

`--file`, `-f`

: Set a logfile or a suffix to the default log file.

`--systemd`, `-s`

: Parse systemd logs from journalctl. Provides all AppArmor logs since the last boot.

`--namespace`, `-n`

: Filter the log to the specified AppArmor namespace.

`--rules`, `-r`

: Convert the log into AppArmor rules.

`--raw`, `-R`

: Print the raw log without any formatting. Useful for reporting logs.

`--since`, `-S`

: Show entries not older than the specified date. It currently only supports log from journalctl (with `--systemd`)

`--boot`, `-b`

: Show entries from the specified boot ID.

`--load`, `-l`

: Load logs from the default `aa-log` output.

`--help`, `-h`

: Print the program usage.


# USAGE

To read the AppArmor log from `/var/log/audit/audit.log`:
```sh
aa-log
```

To optionally filter a given profile name: `aa-log <profile-name>` (your shell will autocomplete the profile name):
```
$ aa-log dnsmasq
DENIED  dnsmasq open /proc/sys/kernel/osrelease comm=dnsmasq requested_mask=r denied_mask=r
DENIED  dnsmasq open /proc/1/environ comm=dnsmasq requested_mask=r denied_mask=r
DENIED  dnsmasq open /proc/cmdline comm=dnsmasq requested_mask=r denied_mask=r
```

To generate AppArmor rule:
```
$ aa-log -r dnsmasq
profile dnsmasq {
  @{PROC}/@{pid}/environ r,
  @{PROC}/cmdline r,
  @{PROC}/sys/kernel/osrelease r,
}
```

# SEE ALSO

`aa-logprof(8)`, `apparmor(7)`, `apparmor.d(5)`, `aa-genprof(1)`, `aa-enforce(1)`, `aa-complain(1)`, `aa-disable(1)`, and
https://apparmor.pujol.io.
