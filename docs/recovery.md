---
title: System Recovery
---

# System Recovery

Issue in some core profiles like the systemd suite, or the desktop environment
can fully break your system. This should not happen a lot, but if it does here
is the process to recover your system on Archlinux:

1. Boot from a Archlinux live USB
1. If you root partition is encryped, decrypt it: `cryptsetup open /dev/<your-disk-id> vg0`
1. Mount your root partition: `mount /dev/<your-plain-disk-id> /mnt`
1. Chroot into your system: `arch-chroot /mnt`
1. Check the AppArmor messages to see what profile is faulty: `aa-log`
1. Temporarily fix the issue with either:
    - When only one profile is faultily, remove it: `rm /etc/apparmor.d/<profile-name>`
    - Otherwise, you can also remove the package: `pacman -R apparmor.d`
    - Alternatively, you may temporarily disable apparmor as it will allow you to
      boot and study the log: `systemctl disable apparmor`
1. Exit, unmount, and reboot:
   ```sh
   exit
   umount -R /mnt
   reboot
   ```
1. Create an issue and report the output of `aa-log`
