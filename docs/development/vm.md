---
title: Development VM
---

To ensure compatibility across distribution, this project ships a wide range of development and tests VM images.

The test VMs can be built locally using [cloud-init](https://cloud-init.io/), [packer](https://www.packer.io/) on Qemu/KVM using Libvirt. No other hypervisor will be targeted for these tests. The files that generate these images can be found in the **[tests/packer](https://github.com/roddhjav/apparmor.d/tree/main/tests/packer)** directory.
The VMs are fully managed using a [justfile](https://github.com/casey/just) that provides an integration environment helper for `apparmor.d`.

```sh
$ just
```

```
Available recipes:
    help                    # Show this help message
    clean                   # Remove all build artifacts

    [build]
    build                   # Build the go programs
    enforce                 # Prebuild the profiles in enforced mode
    complain                # Prebuild the profiles in complain mode
    fsp                     # Prebuild the profiles in FSP mode
    fsp-complain            # Prebuild the profiles in FSP mode (complain)
    fsp-debug               # Prebuild the profiles in FSP mode (debug)

    [install]
    install                 # Install prebuild profiles
    local +names            # Locally install prebuild profiles
    dev name                # Prebuild, install, and load a dev profile

    [packages]
    pkg                     # Build & install apparmor.d on Arch based systems
    dpkg                    # Build & install apparmor.d on Debian based systems
    rpm                     # Build & install apparmor.d on OpenSUSE based systems
    package dist            # Build the package in a clean OCI container

    [tests]
    tests                   # Run the unit tests
    init dist flavor        # Install dependencies for the bats integration tests
    integration dist flavor # Run the integration tests on the machine

    [linter]
    lint                    # Run the linters
    check                   # Run style checks on the profiles

    [docs]
    man                     # Generate the man pages
    docs                    # Build the documentation
    serve                   # Serve the documentation

    [vm]
    img dist flavor         # Build the VM image
    create dist flavor      # Create the machine
    up dist flavor          # Start a machine
    halt dist flavor        # Stops the machine
    reboot dist flavor      # Reboot the machine
    destroy dist flavor     # Destroy the machine
    ssh dist flavor         # Connect to the machine
    list                    # List the machines
    images                  # List the VM images
    available               # List the VM images that can be created

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
Distribution       Flavor    
archlinux          gnome
archlinux          kde
archlinux          server
archlinux          xfce
debian12           gnome
debian12           kde
debian12           server
ubuntu24           server
...
```

A VM image can be build with:

```sh
$ just img archlinux gnome
```

The image will then be showed in the list of images:

```sh
$ just images
```

```
Distribution       Flavor     Size  Date
archlinux          gnome      3.3G  Mar 1 14:49
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
