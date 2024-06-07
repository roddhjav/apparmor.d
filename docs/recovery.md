---
title: System Recovery
---

An issue in some core profiles like the systemd suite, or the desktop environment can prevent your system from starting correctly. This is rare, but if it does happen this is the process to recover your system on an Arch Linux system **without subvolumes**:

1. Boot from an Arch Linux live USB
1. If you root partition is encrypted, decrypt it: `cryptsetup open /dev/<your-disk-id> vg0`
1. Mount your root partition: `mount /dev/<your-plain-disk-id> /mnt`
1. Chroot into your system: `arch-chroot /mnt`
1. Check the AppArmor logs to see which profile is faulty: `aa-log`
1. Temporarily fix the issue with either:
    - When only one profile is causing problems, remove it: `rm /etc/apparmor.d/<profile-name>`
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
