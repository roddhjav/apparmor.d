---
title: Overview
---

Misconfigured AppArmor profiles is one of the most effective ways to break someone's system. This section present the various tests applied to the profiles as well as their current stage of deployment.

**Current**

<div class="grid cards" markdown>

-   :material-github: &nbsp; **[Build](build.md)** `just complain`

    ---

    Build the profiles for all supported distributions.

    - [x] All CI jobs validate the profiles syntax and,
    - [x] ensure they can be safely loaded into a kernel.

-   :octicons-check-24: &nbsp; **[Checks](check.md)** `just check`

    ---

    Checks for common style and security issues:

    - [x] Security checks
    - [x] Style and maintainability checks

-   :material-package: &nbsp; **[Integration Tests](integration.md)** `just test-run`

    Run commands to ensure no logs are raised.

    ---

    - [x] Uses the [bats](https://github.com/bats-core/bats-core) test system.
    - [x] Run in the Github Action as well as in all local [test VM](vm.md).

-   :material-test-tube: &nbsp; **[Distribution Tests](autopkgtest.md)** `just autopkgtest`

    Run the autopkgtest suite for Ubuntu and Debian.

    ---

    - [x] Setup autopkgtest for Ubuntu.
    - [x] Validate profiles on Ubuntu.

</div>


**Future**

For more complex software suite, more integration tests need to be done. The plan is to run existing integration suite from these very software in an environment with `apparmor.d` profiles.

- [ ] Systemd
    - They use mkosi to generate a VM image to run their own integration tests. 
    - See https://www.codethink.co.uk/articles/2024/systemd-integration-testing-part-1/

- [ ] Gnome
    - They use openQA to run their integration tests. 
    - See https://gitlab.gnome.org/GNOME/openqa-tests/
