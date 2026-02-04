---
title: Abstractions
tags:
  - abstractions
---

Abstractions enable resources from one profile to be shared with another and with the system. The table below lists currently supported interfaces, with links to further details for each interface.

This project and the official apparmor-profiles project provide a large selection of abstractions to be included in profiles. They should always be used as they target wide compatibility across hardware and distributions while only allowing the bare minimum access.

!!! example

    For instance, to allow download directory access instead of read and write permissions:
    ```sh
    owner @{HOME}/@{XDG_DOWNLOAD_DIR}/{,**} rw,
    ```

    You should write:
    ```sh
    include <abstractions/user-download-strict>
    ```

!!! warning

    https://snapcraft.io/docs/supported-interfaces

find apparmor.d/abstractions/ -maxdepth 1 -type f  | wc -l

## Architecture

Abstraction are structured in layers as follows:

<div class="grid cards" markdown>

-   **[Layer 0](core.md)**

    ---

    For core atomic functionalities.

    ---

    *This resource uses* `mesa`, `openssl`, `bash-strict`, `gtk-strict`...

-   **[Layer 1](generic.md)**

    ---

    For generic access.

    ---

    *This program needs this resource.* `nameservice-strict`, `authentication`, ...

-   **[Layer 2](kind.md)**

    ---

    For common kind of program.

    ---

    *This program kind is* a game, an electron app

-   **[Layer 3](app.md)**

    ---

    For application

    ---

    *This program is* `sudo`, `firefox`

</div>

## System abstractions

In addition to the above layers, there are abstractions that provide access to system specific part of the system resources.

To use the terminology detailed earlier, these abstractions are `layer -1`

<div class="grid cards" markdown>

-   **[Dbus](dbus.md)**

    ---

    Specific to a dbus interface

    ---

    *This interfaces needs* `bus/system/org.freedesktop.ModemManager1` ...

-   **[sys](sys.md)**

    ---

    sys filesystem access

    ---

    *This program needs/has this resource.* `sys/input`, `sys/hmon`, ...

-   **[udev](udev.md)**

    ---

    For udev device access

    ---

    *This program kind is* `udev/input`

-   **[dev](dev.md)**

    ---

    For device access

    ---

    *This program kind is* `/dev/input`

</div>

