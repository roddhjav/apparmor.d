---
title: Known issues
---

!!! info

    Known bugs are tracked on the meta issue **[#75](https://github.com/roddhjav/apparmor.d/issues/74)**.

## Ubuntu

### Dbus

Ubuntu fully supports dbus mediation with apparmor. If it is a value added by Ubuntu from other distributions, it can also lead to some breakage if you enforce some profiles. *Do not enforce the rules on Ubuntu Desktop.*

Note: Ubuntu server has been more tested and will work without issues with enforced rules.

### Snap

Apparmor.d needs to be fully integrated with snap, otherwise your snap applications may not work properly. As of today, it is a work in progress.


## Complain mode

A profile in *complain* mode cannot break the program it confines. However, there are some **major exceptions**:

1. `deny` rules are enforced even in *complain* mode,
2. `attach_disconnected` (and `mediate_deleted`) will break the program if they are required and missing in the profile,
3. If AppArmor does not find the profile to transition `rPx`.

