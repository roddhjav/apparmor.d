---
title: Security Implementation
icon: material/source-commit-local
---

!!! warning

    This security implementation is still a work in progress. Comments and feedbacks are welcome. [Discuss it on Github](https://github.com/roddhjav/apparmor.d/discussions/1013)

Apparmor is a Mandatory Access Control (MAC) tool. Its only task is to enforce security rules defined by the distribution, or the system administrator. It cannot by itself implement the full security model of an operating system. However, it is an essential tool to achieve such a security goal. This page presents how AppArmor help in the implementation of the [security architecture](architecture.md).
