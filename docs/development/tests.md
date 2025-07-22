---
title: Overview
---

Misconfigured AppArmor profiles is one of the most effective ways to break someone's system. This section present the various tests applied to the profiles as well as their current stage of deployment.

**Current**

- [x] **[Build:](https://gitlab.com/roddhjav/apparmor.d/-/pipelines)** `just complain`
    - Build the profiles for all supported distributions.
    - All CI jobs validate the profiles syntax and ensure they can be safely loaded into a kernel.
    - Ensure the profile entry point (`@{exec_path}`) is defined.

- [x] **[Checks:](https://github.com/roddhjav/apparmor.d/blob/main/tests/check.sh)** `just check` checks basic style of profiles:
    - Ensure apparmor.d header & licence
    - Ensure 2 spaces indentation
    - Ensure local include for profile and subprofiles
    - Ensure abi 4 is used
    - Ensure modern profile naming
    - Ensure `vim:syntax=apparmor`

- [x] **[Integration Tests:](integration.md)** `just test-run <dist> <flavor>`
    - Run simple CLI commands to ensure no logs are raised.
    - Uses the [bats](https://github.com/bats-core/bats-core) test system.
    - Run in the Github Action as well as in all local [test VM](vm.md).

**Plan**

For more complex software suite, more integration tests need to be done. The plan is to run existing integration suite from these very software in an environment with `apparmor.d` profiles.

- [ ] Systemd
    - They use mkosi to generate a VM image to run their own integration tests. 
    - See https://www.codethink.co.uk/articles/2024/systemd-integration-testing-part-1/

- [ ] Gnome
    - They use openQA to run their integration tests. 
    - See https://gitlab.gnome.org/GNOME/openqa-tests/
