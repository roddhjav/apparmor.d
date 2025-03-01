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

Prepare the test environment:
```sh
just img <dist> <flavor>
just vm <dist> <flavor>
```

Run the integration tests on the test VM:
```sh
just integration <dist> <flavor>
```

## Create integration tests

All integration tests are written in [Bats](https://github.com/bats-core/bats-core) and are located in the `tests/integration` directory. The initial tests have been generated using [tldr page](https://tldr.sh/) with the following command:

```sh
go run ./tests/cmd --bootstrap
```
