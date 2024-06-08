---
title: Tests suite
---

A full test suite to ensure compatibility across supported distributions and that software is still considered a work in progress. Here is an overview of the current CI jobs:

**On Gitlab CI**

- Packages build for all supported distributions
- Profiles preprocessing verification for all supported distributions
- Go based command linting, coverage, and unit tests

**On Github Action**

- Integration test on the ubuntu-latest VM: run a simple list of tasks with all the rules enabled and ensure no new issue has been raised. Github Action is used as it offers direct access to a VM with AppArmor included.
