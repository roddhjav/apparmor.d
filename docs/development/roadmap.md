---
title: Roadmap
---

## Toward a stable release

This is the current list of features that must be implemented to get to a stable release

- [x] **[Play machine](https://github.com/roddhjav/play)**

- [ ] **[Sub packages](https://github.com/roddhjav/apparmor.d/issues/464)**
    - [x] Move most profiles into groups
    - [ ] Provide complain/enforced packages version
    - [ ] normal/FSP/server packages variants

- [ ] **Build system**
    - [ ] Continuous release on the main branch, ~2 releases per week
    - [ ] Provide packages repo for ubuntu/debian
    - [x] Add a `just` target to install the profiles in the right place
    - [x] Fully drop the Makefile in favor of `just`

- [ ] **Tests**
    - [x] Tests VM for all supported targets (see [tests/vm](vm.md))
    - [ ] Small integration tests for all core profiles (see [tests/integration](integration.md))

- [ ] **Documentation**
    - [ ] Initial draft of the security model and goal
    - [ ] General documentation improvements

- [ ] **General improvements**
    - [ ] Provide a proper fix for [#74](https://github.com/roddhjav/apparmor.d/issues/74), [#80](https://github.com/roddhjav/apparmor.d/issues/80) & [#235](https://github.com/roddhjav/apparmor.d/issues/235)

- [ ] **Abstractions**
   - [ ] Document all abstractions
   - [ ] Split and reorganize some big abs into set of smaller abstractions.
         Strictly follow the new abstractions guidelines (layer 0, layer 1, etc.)
   - [ ] Abstraction based profiles:
         Most of the accesses needed by GUI based application are commons. As such 80-90% of the profile content should be handled by abstractions (internally they will have conditions).
   - [ ] Test new interface like abstractions
            - notifications
            - audio-bluetooth
            - secrets-service
            - media-keys
            - ...
    - [ ] Rewrite the desktop abstraction to only contains other abs. No direct rules in it.
    - [ ] Rewrite the DE specific abstraction to be a layer 1 abs

- [ ] **Security improvements**
    - [ ] Limit the use of `abstractions/common/systemd`
    - [ ] Ensure systemctl restart/stop/reload is always confined and filtered by unit (dbus only)
    - [ ] Revisit the usae of `systemd-tty-ask-password-agent`

## Next features

- [ ] **Conditions**
    - [ ] Integrate the new condition feature in the profiles and restrict them a lot according to the application actually in use. Eg: `Gnome | KDE`, `X11 | Wayland`, etc.
    - [ ] Create a new `aa-config` tool, similar to seboolean, to manage various settings, based on conditions.

- [ ] **User Data**
    - [ ] Fully rewrite the way user data is allowed / denied. The current implementation requires too much configuration to be usable by everyone.
    - [ ] Add a prompt listener to handle the user data access.

- [x] **[Full System Policy](https://github.com/roddhjav/apparmor.d/issues/252)**
    - [ ] Debug tool to show the profiles transition tree, and ensure no profile is missing
    - [x] Remove the `default` profile

- [ ] **Define roles**
    - [ ] Unrestricted shell role without FSP enabled
    - [ ] Define the roles when FSP is enabled

________________________________________________________________________________

## Done

**General improvements**

- [x] The apt/dpkg profiles has been rewritten

**Abstractions**

- [x] New `audio-client` and `audio-server` abstractions
- [x] New desktop agnostic `desktop` abstraction for all common access for any GUI app. 
- [x] New `graphics` abstraction, hardware-agnostic. Fully replace and restrict the old `opencl` abstractions
- [x] All new abstractions are documented in the [abstractions](abstractions.md) page

**Dbus**

- [x] New `dbus-{system,session,accessibility}` profiles. Works regardless of the dbus implementation in use.
- [x] New talk directive: Allow the application to talk to session services. (send to)
- [x] New own directive: Allow the application to own session services under the given name. (receive, send, bind)
- [x] New `bus-{system,session,accessibility}` abstraction to be used in the profiles

**Directives**

- [x] Add directive. See the [directive](directives.md) page

