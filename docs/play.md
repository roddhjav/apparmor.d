---
title: Play Machine
---

<div class="grid cards" markdown>

-   :material-play: &nbsp; **[play.pujol.io](https://play.pujol.io)**

</div>

!!! info

    A Play Machine is what a system with root as the guest account with only Apparmor to restrict access is called.

## Free root access on an Apparmor machine!

To access my Ubuntu 24.04 play machine ssh to `play.pujol.io` as root, the password is `apparmor`: `ssh root@play.pujol.io`

The aim of this is to:

- Demonstrate that necessary security can be provided by Apparmor without any Unix permissions (however it is still recommended that you use Unix permissions as well for real servers).
- Show that root is not everything in modern security.
- Give a demo machine with [apparmor.d](https://github.com/roddhjav/apparmor.d) fully integrated.

This server is running the [apparmor.d](https://github.com/roddhjav/apparmor.d) project with the [Full System Policies (FSP)](https://apparmor.pujol.io/full-system-policy/) mode enabled and [enforced](https://apparmor.pujol.io/enforce/). Some profiles that only make sense in such a system have been made on purpose.

## Discuss

- Matrix channel available on https://matrix.to/#/#apparmor.d:matrix.org
- Github discussion on: https://github.com/roddhjav/apparmor.d/discussions/722
- If you find a security issue, please report it privately at security@pujol.io

!!! success "Acknowledgement"

    It is a 2025 and Apparmor version of [Russell Coker's SELinux play machine](https://doc.coker.com.au/computers/se-linux-play-machine/).

