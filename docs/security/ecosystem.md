---
title: Ecosystem
icon: simple/linux
---

## Use cases & users

Linux's user can be anyone from novices that does not know they are running Linux and that have no idea that AppArmor exist to the most advanced power user that want full control over their devices.

In opposition of other operating system such as Windows or macOS, Linux can be used in a wide variety of application and use case such as: an OS for embedded device, a mobile OS, a workstation, a gaming console, a server, a global fleet of servers.

The ecosystem is so open that there is no compatibility requirements that define Linux in any way. Particularity in terms of:

1. **Diversity:** They are a wide range of package managers, sandboxing & virtualization tools, init system, display server, Desktop environment.
3. **Users:** *everyone*
2. **Use case:** *everything*
4. **Location:** *everywhere*

This diversity is a strength of Linux but also a challenge when trying to define a security model that can fit all these use cases while being easy to use and understand by everyone.

## Ecosystem

**Architecture**

Despite its diversity, only two major system architecture exists in the Linux ecosystem:

1. **Purposed built systems** such as embedded systems or servers that are designed to run a limited set of tasks and thus that can *easily* support strict and limited security policies.
2. **General purpose systems** such as desktops Linux that are designed to run anything, anyhow.

We note there is trend in the Linux world toward purposed built systems even for desktop usage.

**Security consideration**

Due to the diversity of use case and users in the Linux ecosystem, no particular security requirements/ minimum standard can be assumed for applications or users. However, some general hypothesis can be made:

- Open ecosystem, immense, enforced by package maintainer and distributions
- Any language present, secure or not, past and future
- No rules in software quality
- No rules in security requirements
- The packages are mostly installed from *"trusted source"*
- Distribution repository are well known trusted source
- Software usually do not want to spy on you. Still some do it.
- Continuously find new security vulnerabilities

## Stakeholders

In the Linux ecosystem, there are many stakeholders having some kind of power over the system. These stakeholders can be:

- **The end user:** the person using the device.
- **The device owner:** the company/university/organization owning the device. In case of a personal device, the end user is also the device owner.
- **The Linux vendor:** The Linux distribution
- **The device manufacturer:** The hardware manufacturers providing firmware and drivers
- **Third party app developer/companies:** Other proprietary software company that may run on the device

These stakeholders may have different requirements and objectives regarding the security of the system. For example, the end user may want to protect their privacy, while the device owner may want to ensure that the device is used for work purposes only.

## Requirements

From the above, we can derive the following requirements to apply on our security model.

#### Rqr :material-numeric-1-circle: - User Freedom.

:   Explicitly supports extreme personalization[^1]

[^1]: It does not mean the personalization has to be easy to do, just that it must be possible.

#### Rqr :material-numeric-2-circle: - User Privacy

:   The system should be at the service of the end user in its goal to protect its privacy. The system should not be at the service of the developer, the distributor or a third party application.

#### Rqr :material-numeric-3-circle: - Principle Of Least Astonishment (POLA)

:   A component of a system should behave in a way that most users will expect it to behave, and therefore not astonish or surprise users.

#### Rqr :material-numeric-4-circle: - No strict compatibility

:   Whatever application you have, the system should find a way to run it. Native, containerized, virtualized, emulated, etc.[^2]

[^2]: Some way may be explicitly blocked, but other way should be available.
