---
title: Integration Tests
---

# Integration Tests

!!! danger "Work in Progress"

The purpose of integration testing in apparmor.d is to ensure the profiles are not going to break a program when used in the Linux distribution and desktop environment we support.

**Workflow**

1. Build some tests VM
2. Start the VM, do some dev
3. Run the integration test against a given test VM
4. Ensure no new logs have been raised


## Test Virtual Machines

The test VMs are build using [`cloud-init`][cloud-init], [`packer`][packer], and [`vagrant`][vagrant] on Qemu/KVM using Libvirt. No other hypervisor will be targeted for these tests. The files that generate these images can be found in the **[tests/packer](https://github.com/roddhjav/apparmor.d/tree/main/tests/packer)** directory.

[cloud-init]: https://cloud-init.io/
[packer]: https://www.packer.io/
[vagrant]: https://www.vagrantup.com/


### Build

!!! note

    You may need to edit some settings to fit your setup:
     
    - The libvirt configuration in `tests/Vagrantfile` 
    - The default ssh key and ISO directory in `tests/packer/variables.pkr.hcl`

**Build an image**

To build a VM image for development purpose, run the following from the `tests` directory:

| Distribution | Flavor | Build command | VM name |
|:------------:|:------:|:-------------:|:-------:|
| Archlinux | Gnome | `make archlinux flavor=gnome` | `arch-gnome` |
| Archlinux | KDE | `make archlinux flavor=kde` | `arch-kde` |
| Ubuntu | Server | `make ubuntu flavor=server` | `ubuntu-server` |
| Ubuntu | Desktop | `make ubuntu falvor=desktop` | `ubuntu-desktop` |

**VM management**

The development workflow is done through vagrant:

* Star a VM: `vagran up <name>`
* Shutdown a VM: `vagrant halt <name>`
* Reboot a VM: `vagrant reload <name>`

The available VM `name` are defined in the `tests/boxes.yml` file


### Develop

**Credentials**

The admin user is: `user`, its password is: `user`. It has passwordless sudo access. Automatic login is **not** enabled on DE. The root user is not locked.

**Directories**

All the images come pre-configured with the lastest version of `apparmor.d` installed and running in the VM. The apparmor.d is mounted as `/home/user/Projects/apparmor.d`

**Usage**

On all images, `aa-update` can be used to rebuild and install latest version of the profiles. `p`, `pf`, and `pu` are two preconfigured aliases of `ps` that show the security status of processes. `htop` is also configured to show this status.
