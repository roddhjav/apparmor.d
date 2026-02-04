---
title: Security Hardening
icon: material/bandage
---

A careful reader would have noticed that these sections do not mention anything about security hardening, whilst it is a common thing in modern secure system.

<figure markdown>

**What is the problem with security hardening?**

</figure>

In one sentence, *It is the opposite of security by design.* Security hardening means retrofitting security controls onto an existing system, while *secure by design* means building security into the architecture from the beginning.

Hardening a system usually means installing a set of security tools and configurations on top of an existing system because they look cool, secure or trendy. The problem is that when doing so, you usually do not have a clear idea of what you are trying to achieve, what are the threats you are trying to mitigate, and what are the trade-offs you are making.

!!! tip "To sum up"

    Hardening is like adding additional reinforcement to a bridge because after finishing it, you realized it could be dangerous. It would have been way better to design the bridge with more structural support in mind from the start.

**Example**

- Instead of disabling some kernel modules, it is better to build the kernel without the module at all in such a way that even if an attacker manage to load the module, it would not be possible.

- Instead of disabling USB storage devices because they could be used to exfiltrate data, it would be way better to design the system in such a way that even if a USB storage device is connected, it cannot be used to exfiltrate data.

In other words, it is way more secure, and stable to design a system with security in mind from the start, rather than trying to patch it afterwards.

This is why these sections focus on the complex question of *what do you want to achieve?* rather than the simpler question of *how to harden your system?*.
