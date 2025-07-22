---
title: Roadmap
---

## Toward a stable release

This is the current list of features that must be implemented to get to a stable release

- [x] **Play machine**

- [ ] **[Sub packages](https://github.com/roddhjav/apparmor.d/issues/464)** 
    - [x] Move most profiles into groups such that 
    - [ ] New simplified build system to generate the packages with profile dependencies check

- [ ] **Tests**
    - [x] Tests VM for all supported targets (see [tests/vm](vm.md))
    - [ ] Small integration tests for all core profiles (see [tests/integration](integration.md))

- [ ] **Documentation**
    - [ ] Initial draft of the security model and goal
    - [ ] General documentation improvements

- [ ] **General improvements**
    - [ ] Provide a proper fix for [#74](https://github.com/roddhjav/apparmor.d/issues/74), [#80](https://github.com/roddhjav/apparmor.d/issues/80) & [#235](https://github.com/roddhjav/apparmor.d/issues/235)
    - [x] The apt/dpkg profiles needs to be reworked

- [ ] Build system
    - [ ] Continuous release on the main branch, ~2 releases per week
    - [ ] Provide packages repo for ubuntu/debian
    - [ ] Provide complain/enforced packages version
    - [x] Add a `just` target to install the profiles in the right place
    - [x] Fully drop the Makefile in favor of `just`

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

## Done

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

