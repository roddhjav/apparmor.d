---
title: Security Architecture
icon: material/wallpaper
---

!!! warning

    This security architecture is still a work in progress. Comments and feedbacks are welcome. [Discuss it on Github](https://github.com/roddhjav/apparmor.d/discussions/1013)

## Preliminaries

<figure markdown>

**The best is the enemy of the good.**

</figure>

The architecture presented here can be seen as a general overview of what any modern Linux security construction may want to achieve.

Example of current Linux distribution implementing something similar with various use case in mind are:

<div class="grid cards" markdown>

- :material-package: &nbsp; **[ClipOS](https://clip-os.org/en/)**
- :material-ubuntu: &nbsp; **[Ubuntu Core](https://ubuntu.com/core)**
- :material-fedora: &nbsp; **[Fedora Atomic Desktops](https://fedoraproject.org/atomic-desktops/)**
- :material-fedora: &nbsp; **[Fedora Core OS](https://fedoraproject.org/coreos/)**
- :material-atom: &nbsp; **[Particle OS](https://github.com/systemd/particleos)**
- :material-train-car-flatbed-car: &nbsp; **[Flatcar OS](https://www.flatcar.org/)**
- :simple-opensuse: &nbsp; **[openSUSE MicroOS](https://get.opensuse.org/microos/)**
- :simple-opensuse: &nbsp; **[](https://)**

</div>

A careful reader would have noticed that the common ground among these distributions is to be constituted of a fully immutable core system. If such a construction is probably the future of Linux, as of today it can raise some usability concerns (cf [rule :material-numeric-6-circle:](#user-freedom "User freedom.")). Therefore, the current project propose a pragmatic long term solution:

1. We acknowledge the end goal need to be fully compatible with system that respect 100% of the security model presented here.
2. We stay compatible with a *"classic"* Linux construction (i.e. without immutable core), and try to implement as mush as we can on classic distribution. It is considered as a transitional state.

## Security Architecture

As attacker usually comes from the top (a high level application) and goes down to the core system, we present the security architecture similarly.

1. **In application sandboxing:** separation of privilege, least privilege principle, within different process of the application itself. It can be implemented using tools such as: Landlock, bwrap, and Apparmor.
The purpose is to separate highly privileged code from the rest of the application.

1. **Application sandboxing:** Isolate the application from the rest of the system with tools such as. In this context, "sandboxing" does not refer to a special technology, but to the general concept. It can be implemented in various ways: trough VM (Qemu/KVM, Firecracker, Cloudhypervisor, Kata Containers), or container (Docker, gVisor, Flatpak, Snap) with different level of isolation and integration with the rest of the system.

1. **Confined user:** The user is confined in its own environment. Limiting what they can do **[:material-police-badge-outline:{ .pg-red }](../full-system-policy.md "Only for Full System Policy (FSP)")**.

1. **System confinement:** separation of privilege, the least privilege principle, within all part of the core system. For example, Gnome is constituted by about 50 small services running in the background. Each of them are confined independently with granular access to other part of the DE.

<figure markdown>
  ![_security-implementation-architecture](arch1.png#only-light)
  ![_security-implementation-architecture](arch2.png#only-dark)
  <figcaption>Overview of the security architecture</figcaption>
</figure>

## Security Levels

As `apparmor.d` can be used in multiple security model, it provides different mode to fit into the following **proposed** security target.

**`default`**

This confinement level can be summarized as *the most we can do without requiring user configuration*. It provides a good level of security for most use cases while requiring minimal configuration from the user. It is suitable for general purpose use. The goal with this level is that the end user should not be aware of AppArmor and its configuration.

- Principle of least privilege applied
- Strict Principle Of Least Astonishment (POLA)

**`strict`**

Stronger confinement for applications. Suitable for more sensitive use cases. This level includes a few sublevel depending on how the user want to configure it.

**`fsp`**

Full System Policy confinement for applications. Suitable for very high security use cases.

**`extreme`**

Maximum confinement for applications. Suitable for the most extreme security use cases. This level may break some applications and require significant user and application configuration as well as patching some applications to work properly.
The goal is to provide a kind of Multi Category Security (MCS) using apparmor on top of the FSP model.
