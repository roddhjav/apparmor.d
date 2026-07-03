---
title: Tunables
tags:
  - tunables
---

Tunables are global variables used to adjust AppArmor policies to match the local system environment without modifying the core profiles themselves (e.g., globally defining the `@{HOME}` directory path).

This project and the official apparmor-profiles project provide a large selection of tunables to be included in profiles. They should always be used as they target wide compatibility across hardware and distributions while only allowing the bare minimum access.

!!! example

    For instance, to allow download directory access instead of read and write permissions:
    ```sh
    owner /home/me/Downloads/{,**} rw,
    ```

    You should write:
    ```sh
    owner @{HOME}/@{XDG_DOWNLOAD_DIR}/{,**} rw,
    ```

<div class="grid cards" markdown>

-   **[Default](default/base.md)**

    ---

    Default tunables from the upstream apparmor project. They are available by any profile even without using apparmor.d as a dependency.

-   **[apparmor.d](base.md)**

    ---

    All new tunables provided by the apparmor.d project.

</div>
