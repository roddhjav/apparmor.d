---
title: Usage
---

## Enabled profiles

Once installed and with the rules enabled, you can ensure the rules are loaded with:
```sh
sudo aa-status
```

It should give something like:
```
apparmor module is loaded.
1613 profiles are loaded.
1050 profiles are in enforce mode.
   ...
563 profiles are in complain mode.
   ...
0 profiles are in kill mode.
0 profiles are in unconfined mode.
170 processes have profiles defined.
140 processes are in enforce mode.
   ...
30 processes are in complain mode.
   ...
0 processes are in prompt mode.
0 processes are in kill mode.
0 processes are unconfined but have a profile defined.
0 processes are in mixed mode.
```

You can also list the current processes alongside with their security profile with:
```sh
ps auxZ
```

Most of the processes should then be confined:
```
unconfined                      root        /usr/lib/systemd/systemd --switched-root --system --deserialize 33
systemd-udevd (complain)        root        /usr/lib/systemd/systemd-udevd
systemd-journald (complain)     root        /usr/lib/systemd/systemd-journald
rngd (complain)                 root        /usr/bin/rngd -f
systemd-timesyncd (complain)    systemd+    /usr/lib/systemd/systemd-timesyncd
auditd (complain)               root        /sbin/auditd
acpid (complain)                root        /usr/bin/acpid --foreground --netlink
dbus-daemon (complain)          dbus        /usr/bin/dbus-daemon --system --address=systemd: --nofork --nopidfile --systemd-activation --syslog-only
power-profiles-daemon (complain) root       /usr/lib/power-profiles-daemon
systemd-logind (complain)       root        /usr/lib/systemd/systemd-logind
systemd-machined (complain)     root        /usr/lib/systemd/systemd-machined
NetworkManager (complain)       root        /usr/bin/NetworkManager --no-daemon
polkitd (complain)              polkitd     /usr/lib/polkit-1/polkitd --no-debug
gdm (complain)                  root        /usr/bin/gdm
accounts-daemon (complain)      root        /usr/lib/accounts-daemon
rtkit-daemon (complain)         rtkit       /usr/lib/rtkit-daemon
packagekitd (complain)          root        /usr/lib/packagekitd
colord (complain)               colord      /usr/lib/colord
unconfined                      user        /usr/lib/systemd/systemd --user
unconfined                      user        (sd-pam)
gdm-wayland-session (complain)  user        /usr/lib/gdm-wayland-session /usr/bin/gnome-session
gnome-session-binary (complain) user        /usr/lib/gnome-session-binary
gnome-session-ctl (complain)    user        /usr/lib/gnome-session-ctl --monitor
gnome-session-binary (complain) user        /usr/lib/gnome-session-binary --systemd-service --session=gnome
gnome-shell (complain)          user        /usr/bin/gnome-shell
...
ps (complain)                   user        ps auxZ
```

??? info "Display the process hierarchy and hide the kernel thread in `ps`"
 
    In order to list above processes with displaying the process hierarchy you can, alternatively, use `ps auxfZ`. 
    
    To hide the kernel thread in `ps` use `LIBPROC_HIDE_KERNEL=1 ps auxfZ`. You can
    add an alias in your shell:
    ```sh
    alias p="LIBPROC_HIDE_KERNEL=1 ps auxfZ"
    ```


## AppArmor Log

Ensure that `Auditd` is installed and running on your system in order to read AppArmor log from `/var/log/audit/audit.log`. Then you can see the log with the provided command `aa-log` allowing you to review AppArmor generated messages in a colourful way.

Other AppArmor userspace tools such as `aa-enforce`, `aa-complain`, and `aa-logprof` should work as expected. You can also configure [a desktop notification on denied actions](https://wiki.archlinux.org/title/AppArmor#Get_desktop_notification_on_DENIED_actions).


### Basic use

To read the AppArmor log from `/var/log/audit/audit.log`:
```sh
$ aa-log
```

To optionally filter a given profile name: `aa-log <profile-name>` (your shell will autocomplete the profile name):
```
$ aa-log dnsmasq
DENIED  dnsmasq open /proc/sys/kernel/osrelease comm=dnsmasq requested_mask=r denied_mask=r
DENIED  dnsmasq open /proc/1/environ comm=dnsmasq requested_mask=r denied_mask=r
DENIED  dnsmasq open /proc/cmdline comm=dnsmasq requested_mask=r denied_mask=r
```

To generate AppArmor rule:
```sh
$ aa-log -r dnsmasq
profile dnsmasq {
  @{PROC}/@{pid}/environ r,
  @{PROC}/cmdline r,
  @{PROC}/sys/kernel/osrelease r,
}
```

!!! info

    Other logs file in `/var/log/audit/` can easily be checked: `aa-log -f 1`
    parses `/var/log/audit/audit.log.1`.


### Help

```
aa-log [-h] [--systemd] [--file file] [--load] [--rules | --raw] [--since] [--namespace] [profile]

    Review AppArmor generated messages in a colorful way. It supports logs from
    auditd, systemd, syslog as well as dbus session events.

    It can be given an optional profile name to filter the output with.

    Default logs are read from '/var/log/audit/audit.log'. Other files in
    '/var/log/audit/' can easily be checked: 'aa-log -f 1' parses 'audit.log.1'
    Use 'aa-log -f -' to read from standard input.

    Logs written with 'aa-log' can be read again with 'aa-log -l'.

Options:
    -h, --help         Show this help message and exit.
    -f, --file FILE    Set a logfile or a suffix to the default log file.
    -s, --systemd      Parse systemd logs from journalctl.
    -n, --namespace NS Filter the logs to the specified namespace.
    -r, --rules        Convert the log into AppArmor rules.
    -R, --raw          Print the raw log without any formatting.
    -b, --boot NUM     Show entries from the specified boot.
    -S, --since DATE   Show entries not older than the specified date.
    -l, --load         Load logs from the default aa-log output.
```
