---
title: Threat model
icon: material/bomb
---

The importance of the threat depends on the use case & application. For example the fact that *adversaries can get physical access to device* is a bigger concern on mobile & embedded device than on a server or even on a VM.
However, these remain valid anyway.

## Modularity

Across the linux ecosystem, treats can varies. As such not all threats are relevant to all users and some threats only matter for some security models. As such the underlying implementation should be modular enough to allow selection of the relevant threats and to enforce a given security model even if it lead to a trade off in terms of usability.

## Threats

!!! warning "Work in progress"

    This threats are not yet properlly defined. Comments and feedbacks are welcome. [Discuss it on Github](https://github.com/roddhjav/apparmor.d/discussions/1013)

### `PA` Physical access

Adversaries can get physical access to the devices.

### `C` Communication

Network communication is untrusted

### `P` Platform

The Linux system can be targeted.

### `UI` User interaction

Many stakeholders in the ecosystem can act as supply chain attack vectors.

## Out of scope

Some threats are considered out of scope for various reasons.
