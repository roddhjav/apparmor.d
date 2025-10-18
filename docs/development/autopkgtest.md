---
title: Autopkgtest
---

**autopkgtest** is Debian's automated package testing framework that validates packages work correctly after installation in clean VM. To ensure real-world functionality, it performs integration testing by installing packages and running tests defined in `debian/tests/`. It is thus a good method to validate the apparmor profiles.

!!! note

    The autopkgtest suite integration in apparmor.d is currently a work in progress

**Workflow**

1. Create a testing VM for autopkgtest
2. Run autopkgtest on a wide range of source package
3. Continuously collect AppArmor logs during the tests


## VM Setup and Management

**Create the test VM**

The test VM is a VM as defined in the [Development VM](vm.md) section, with a specific cloud-init configuration for autopkgtest. We use the same `setup-testbed` script to prepare the VM as `autopkgtest-build-qemu`. In addition, we ensure the VM is built with the lastest `apparmor.d` profiles in test mode.

You can create the image, then the VM, and shut it down with:

```sh
just img <dist> <release> test
just create <dist><release> test
just halt <dist><release> test
```

Example:

```sh
just img ubuntu 25.10 test
just create ubuntu25.10 test
just halt ubuntu25.10 test
```

**Update `apparmor.d` in the VM**

Others VM defined in this project ships with a `aa-update` command that build and update the package. This does not apply to the `test` flavor because:

1. We do not want to mount this project to a VM where the tests can be destructive
2. The `setup-testbed` script gets rid of most build dependencies for `apparmor.d`

To update apparmor.d in the VM without creating a new image, use the `autopkgtest-update` command, it will build the package on the host, and install it in the VM:

```sh
just autopkgtest-update <dist> <release>
```

Example:

```sh
just autopkgtest-update ubuntu 25.10
```

## Test Execution Workflow

The autopkgtest suite runs the tests for all source packages listed in `tests/autopkgtest/src-packages`. It installs each package in the test VM, runs its autopkgtest suite, and monitors AppArmor logs for any policy violations. It is possible to control the range of packages tested using alphabetical start and end points in the `tests/autopkgtest/autopkgtest.sh` script.

To run the full suite for a Debian/Ubuntu system:

```sh
just autopkgtest <dist><release>
```

Example:

```sh
just autopkgtest ubuntu25.10
```

## Log Analysis

The full raw logs are available in the `.logs/autopkgtest/` directory. One can run the following commands to analyze the logs and generate missing rules:

Report all collected logs using `aa-log`
```sh
just autopkgtest-log
```

Generate missing rules using `aa-log --rules`
```sh
just autopkgtest-rules
```
