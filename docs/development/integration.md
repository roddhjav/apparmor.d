---
title: Integration Tests
---

The purpose of integration testing in apparmor.d is to ensure the profiles are not going to break programs found in Linux distributions and Desktop Environment that we support.

Although the integration test suite is intended to be run in a [Development VM](vm.md), it is also deployed the GitHub Action pipeline.

**Workflow**

1. Create a testing VM
2. Run the integration tests against the testing VM
3. Ensure no new logs have been raised

## Getting started

**Prepare the test environment:**
```sh
just img <dist> <flavor>
just create <dist> <flavor>
```

Example:
```sh
just img ubuntu25 desktop
just create ubuntu25 desktop
```

**Install dependencies for the integration tests**
```sh
just tests-init <dist> <flavor>
```

Example:
```sh
just tests-init ubuntu25 desktop
```

**Run the integration tests**

It: synchronizes the tests, unmount the shared directory, then run the tests.
```sh
just tests-run <dist> <flavor>
```

Example:
```sh
just tests-run ubuntu25 desktop
```

Partial tests can also be run. For example the following command will only run the tests in the `tests/integration/apt` directory on the `ubuntu25` `desktop` machine:
```sh
just tests-run ubuntu25 desktop apt
```

## Create integration tests

All integration tests are written in [Bats](https://github.com/bats-core/bats-core) and are located in the `tests/integration` directory. The initial tests have been generated using [tldr page](https://tldr.sh/) with the following command:

```sh
go run ./tests/cmd --bootstrap
```
