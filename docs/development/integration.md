---
title: Integration Tests
---

!!! danger "Work in Progress"

The purpose of integration testing in apparmor.d is to ensure the profiles are not going to break programs found in Linux distributions and Desktop Environment that we support.

**Workflow**

1. Create a testing VM
2. Start the VM, do some dev
3. Run the integration tests against the testing VM
4. Ensure no new logs have been raised


## Test Virtual Machines

The test VMs are built using [`cloud-init`][cloud-init] (when available), [`packer`][packer], and [`vagrant`][vagrant] on Qemu/KVM using Libvirt. No other hypervisor will be targeted for these tests. The files that generate these images can be found in the **[tests/packer](https://github.com/roddhjav/apparmor.d/tree/main/tests/packer)** directory.

[cloud-init]: https://cloud-init.io/
[packer]: https://www.packer.io/
[vagrant]: https://www.vagrantup.com/

### Requirements

* docker
* [packer]
* [vagrant]
* vagrant plugin install vagrant-libvirt

!!! note

    You may need to edit some settings to fit your setup:

    - The libvirt configuration in `tests/Vagrantfile` 
    - The default ssh key and ISO directory in `tests/packer/variables.pkr.hcl`

### Build

**Build an image**

To build a VM image for development purpose, run the following from the `tests` directory:

| Distribution | Flavor | Build command | VM name |
|:------------:|:------:|:-------------:|:-------:|
| Arch Linux | Gnome | `make archlinux flavor=gnome` | `arch-gnome` |
| Arch Linux | KDE | `make archlinux flavor=kde` | `arch-kde` |
| Debian | Server | `make debian flavor=server` | `debian-server` |
| openSUSE | KDE | `make opensuse flavor=kde` | `opensuse-kde` |
| Ubuntu | Server | `make ubuntu flavor=server` | `ubuntu-server` |
| Ubuntu | Desktop | `make ubuntu falvor=desktop` | `ubuntu-desktop` |

**VM management**

The development workflow is done through vagrant:

* Star a VM: `vagran up <name>`
* Shutdown a VM: `vagrant halt <name>`
* Reboot a VM: `vagrant reload <name>`

The available VM `name` is defined in the `tests/boxes.yml` file


### Develop

**Credentials**

The admin user is: `user`, its password is: `user`. It has passwordless sudo access. Automatic login is **not** enabled on DE. The root user is not locked.

**Directories**

All the images come pre-configured with the latest version of `apparmor.d` installed and running in the VM. apparmor.d is mounted as `/home/user/Projects/apparmor.d`

**Usage**

On all images, `aa-update` can be used to rebuild and install the latest version of the profiles. `p`, `pf`, and `pu` are two pre-configured aliases of `ps` that show the security status of processes. `htop` is also configured to show this status.


## Tests

!!! warning

    The test suite is expected to be run in a [VM](#test-virtual-machines)

### Getting started

Prepare the test environment:
```sh
cd tests
make <dist> falvor=<flavor>
AA_INTEGRATION=true vagrant up <name>
```

Run the integration tests on the test VM:
```sh
make integration box=<dist> IP=<ip>
```

### Create integration tests

**Test suite usage**

Initialise the tests with:
```sh
./aa-test --bootstrap
```

List the tests scenarios to be run
```sh
./aa-test --list
```

Start the tests and collect the results
```sh
./aa-test --run
```

**Tests manifest**

A basic set of test is generated on initialization. More tests can be manually written in yaml file. They must have the following structure:

```yaml
- name: acpi
  profiled: true
  root: false
  require: []
  arguments: {}
  tests:
    - dsc: Show battery information
      cmd: acpi
      stdin: []
    - dsc: Show thermal information
      cmd: acpi -t
      stdin: []
    - dsc: Show cooling device information
      cmd: acpi -c
      stdin: []
    - dsc: Show thermal information in Fahrenheit
      cmd: acpi -tf
      stdin: []
    - dsc: Show all information
      cmd: acpi -V
      stdin: []
    - dsc: Extract information from `/proc` instead of `/sys`
      cmd: acpi -p
      stdin: []
```
