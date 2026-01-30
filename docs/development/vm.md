---
title: Development VM
---

To ensure compatibility across distribution, this project ships a wide range of development and tests VM images.

The test VMs can be built locally using [cloud-init](https://cloud-init.io/), [packer](https://www.packer.io/) on Qemu/KVM using Libvirt. No other hypervisor will be targeted for these tests. The files that generate these images can be found in the **[tests/packer](https://github.com/roddhjav/apparmor.d/tree/main/tests/packer)** directory.

The VMs are fully managed using a [Justfile](https://github.com/casey/just) that provides an integration environment helper for `apparmor.d`.

```sh
$ just
```

```
Available recipes:
    help                              # Show this help message
    clean                             # Remove all build artifacts

    ...

    [vm]
    img dist release flavor           # Build the VM image
    create osinfo flavor              # Create the machine
    up osinfo flavor                  # Start a machine
    halt osinfo flavor                # Stops the machine
    reboot osinfo flavor              # Reboot the machine
    destroy osinfo flavor             # Destroy the machine
    snapshots osinfo flavor           # List all snapshots for a machine
    snapshot osinfo flavor snapname   # Snapshot a machine
    restore osinfo flavor snapname    # Restore a machine to a specified snapshot
    delete osinfo flavor snapname     # Delete a specified snapshot from a machine
    ssh osinfo flavor                 # Connect to the machine
    mount osinfo flavor               # Mount the shared directory on the machine
    umount osinfo flavor              # Unmout the shared directory on the machine
    list                              # List the machines
    images                            # List the VM images
    available                         # List the VM images that can be created

    ...

Build variables available:
    build        # Build directory (default: .build)
    destdir      # Installation destination (default: /)
    pkgdest      # Package output directory (default: /home/alex/06_Projects/Security/apparmor.d/.pkg)
    opt          # Prebuild option, only used for the dev install target (default: complain)

Development variables available:
    username     # VM username (default: user)
    password     # VM password (default: user)
    disk_size    # VM disk size (default: 40G)
    vcpus        # VM CPU (default: 6)
    ram          # VM RAM (default: 4096)

See https://apparmor.pujol.io/development/ for more information.
```

## Requirements

* [docker](https://www.docker.com/)
* [just](https://github.com/casey/just)
* [packer](https://www.packer.io/)
* [libvirt](https://libvirt.org/)
* [qemu](https://www.qemu.org/)

!!! note

    You may need to edit some settings to fit your setup:

    - The default ssh key and ISO directory in `tests/packer/variables.pkr.hcl`

## Build

One can see the available images by running:

```sh
$ just available
```

```
Distribution  Release  Flavor
archlinux     -        gnome
archlinux     -        kde
debian        13       gnome
debian        13       server
debian        13       test
opensuse      -        gnome
opensuse      -        kde
ubuntu        24.04    server
ubuntu        25.05    desktop
ubuntu        25.05    kubuntu
ubuntu        25.10    test

...
```

A VM image can be build with:

```sh
$ just img archlinux - gnome
```

The image will then be showed in the list of images:

```sh
$ just images
```

```
OsInfo       Flavor   Size   Date      
archlinux    gnome    3.5GB  Sep   25  23:25
```

The VM can then be created with:

```sh
$ just create archlinux gnome
```

And connected to with:

```sh
$ just ssh archlinux gnome
```

## Develop

**Credentials**

The admin user is: `user`, its password is: `user`. It has passwordless sudo access. Automatic login is **not** enabled on DE. The root user is not locked.

**Directories**

All the images come pre-configured with the latest version of `apparmor.d` installed and running in the VM. The apparmor.d project directory is mounted as `/home/user/Projects/apparmor.d`

**Usage**

On all images, `aa-update` can be used to rebuild and install the latest version of the profiles. `p`, `pf`, and `pu` are two pre-configured aliases of `ps` that show the security status of processes. `htop` is also configured to show this status.
