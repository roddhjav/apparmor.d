---
title: Programs
tags:
  - tunables
  - default
---

## dovecot

@{DOVECOT_MAILSTORE} is a space-separated list of all directories
where dovecot is allowed to store and read mails

The default value is quite broad to avoid breaking existing setups.
Please change @{DOVECOT_MAILSTORE} to (only) contain the directory
you use, and remove everything else.

### @{DOVECOT_MAILSTORE}

[:material-file-eye:](https://gitlab.com/apparmor/apparmor/-/blob/9a402e9e6a04b45dd788c7dfb0106dae45443b01/profiles/apparmor.d/tunables/dovecot#L19 "View source"){ .abs-source }

```
@{DOVECOT_MAILSTORE}=@{HOME}/Maildir/ @{HOME}/mail/ @{HOME}/Mail/ /var/vmail/ /var/mail/ /var/spool/mail/
```

## gs

### @{gs_file_ext}

[:material-file-eye:](https://gitlab.com/apparmor/apparmor/-/blob/9a402e9e6a04b45dd788c7dfb0106dae45443b01/profiles/apparmor.d/tunables/gs#L12 "View source"){ .abs-source }

```
@{gs_file_ext}=[pP][dD][fF] [pP][sS] [eE][pP][sS] [eE][pP][sS][iI] [pP][nN][gG] [jJ][pP][gG] [jJ][pP][eE][gG] [pP][nN][mM] [tT][iI][fF] [tT][iI][fF][fF] [bB][mM][pP] [pP][cC][xX] [pP][sS][dD] [tT][xX][tT] [pP][xX][lL] [dD][oO][cC][xX] [xX][pP][sS] [xX][mM][lL] [iI][cC][cC] [rR][dD][fF]
```

## ntpd

### @{NTPD_DEVICE}

[:material-file-eye:](https://gitlab.com/apparmor/apparmor/-/blob/9a402e9e6a04b45dd788c7dfb0106dae45443b01/profiles/apparmor.d/tunables/ntpd#L14 "View source"){ .abs-source }

Add your ntpd devices here eg. if you have a DCF clock

```
@{NTPD_DEVICE}="/dev/tty10"
```

## print-devices

@{print_devices} is a space-separated list of all devices
representing locally connected printers

### @{print_devices}

[:material-file-eye:](https://gitlab.com/apparmor/apparmor/-/blob/9a402e9e6a04b45dd788c7dfb0106dae45443b01/profiles/apparmor.d/tunables/print-devices#L15 "View source"){ .abs-source }

```
@{print_devices}=/dev/lp* /dev/ttyS* /dev/ttyUSB* /dev/usb/lp* /dev/parport*
```
