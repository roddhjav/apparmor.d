---
title: Roadmap
---

## Toward a stable release

!!! info

    See https://github.com/roddhjav/apparmor.d/issues/967 for more information about current work

This is the current list of features that must be implemented to get to a stable release

- [ ] **[Packages](https://github.com/roddhjav/apparmor.d/issues/986)**
    - [x] Tool, and core packages
    - [ ] Move most profiles into groups
    - [ ] Move from complain/enforced packages version to local installer
    - [ ] Finalise `aa-install` and `aa-config`

- [ ] **[Tests](tests.md)**

    *Most of this is currently impossible with the current infrastructure, as it goes way beyond what free tiers offer on Gitlab/Github CI*

    - [ ] Integration tests for **all** core profiles usings `bats` (see [tests/integration](integration.md))
    - [ ] Orchestrate all bats tests with spread on all supported distributions
    - [ ] Put autopkgtests in a pipeline for development and stable version of Debian/Ubuntu.
    - [ ] Use integration suite from some software (e.g. systemd, gnome, kde...)

- [ ] **General improvements**
    - [ ] Provide a proper fix for [#74](https://github.com/roddhjav/apparmor.d/issues/74), [#80](https://github.com/roddhjav/apparmor.d/issues/80) & [#235](https://github.com/roddhjav/apparmor.d/issues/235): some programs, started by dbus need `px` instead of `Px`.

- [ ] **Conditions**
    - [x] Integrate the new condition feature in the profiles and restrict them a lot according to the application actually in use. Eg: `Gnome | KDE`, `X11 | Wayland`, etc.
    - [ ] Create a new `aa-config` tool, similar to seboolean, to manage various settings, based on conditions.

- [ ] **Abstractions**
    - [x] Document all abstractions
    - [x] Split and reorganize some big abstractions into set of smaller abstractions.
          Strictly follow the new abstractions guidelines (layer 0, layer 1, etc.)
    - [ ] Abstraction based profiles

          Most of the accesses needed by GUI based application are commons. As such 80-90% of the profile content should be handled by abstractions (internally they will have conditions):

        - [x] `common/gnome`
        - [ ] `common/kde`

    - [x] Test new interface like abstractions:
        - [x] notifications
        - [x] bluetooth-observe
        - [x] bluetooth-control
        - [x] secrets-service
        - [x] media-keys
        - [ ] ...

     - [x] Rewrite the desktop abstraction to only contains other abs. No direct rules in it.
     - [x] Rewrite the DE specific abstraction to be a layer 1 abs
     - [ ] Rewrite some abstraction to fully support conditions based on system state
        - [ ] `X-strict`
        - [ ] `wayland-strict`
        - [ ] `desktop` (and deprecate `gnome-strict`, `kde-strict`)

- [ ] **Access Re-structuration** (follow up of **Abstractions**)

      We already include a lot of new abstractions and provides different layers of abstractions for various use cases. Eg: `desktop` for common GUI access, `icons` if only icon access is needed, etc. It is however not enough. The general concept is that to be distribution and hardware agnostic while not giving up on security, most profiles should only be consisted of abstractions. These abstractions could however be as specific
      as needed.

      The solution would be to add hundreds of new small abstractions to define every
      possible access, and then combine them in the profiles.
        
     *Example:*

      - `udev/input` for input devices
      - `sys/input` for input devices sysfs access
      - `devices/input` for /dev/input/*

      With this structure; current low level abstraction (such as `disks-read`) could:

      - Be a set of smaller abstractions
      - Be directly used by fewer profiles. Profiles that only needs some disks access but not all could use the smaller abstractions instead (eg: `devices/nvme-{read,write}`).

     See [abstractions/sys](../abstractions/sys.md) for an overview of existing work on this subject.

- [ ] **`aa-learn`: automatic personalization**

    SELinux is well known to piss of users that want to use special paths or feature for anything: they need to add they paths to the policies or configure boolean. Meanwhile, it provides a *huge* security improvement, as this also makes policies stricter.

    However, it is already complex to system admin on server, we cannot afford this complexity on desktop. On a default install, desktop users should not even know that apparmor is there. Thus, for this default install, we need a way to automatically learn the user configuration, and adapt the profiles accordingly. We need to automatically discover:

    - Special user paths for things such as themes, icon in a git profile (dotfiles)
    - Use of specific software to integrate with.

    This could take the form of a learning mode (based on the complain mode), that will analyze the logs and discover specific user requirements, and automatically add local extension to the abstractions.

    `aa-learn` would differ from `aa-logprof` or `aa-log` as it is not meant to create new profiles or to update existing ones. It is meant to adapt the rules in abstractions to the user needs.

- [ ] **Security improvements**

    - [x] Bwrap in its own profile, stack the rest of the sandbox
    - [x] Bwrap in its own namespace (only when needed)
    - [x] Systematic use of subprofile for sensitive tasks:
        - [x] sudo
        - [x] ldconfig
        - [x] ldd
        - [ ] kmod (partial, we cannot filter the modules)
        - [ ] systemctl (partial, the subprofile is not that limited)
    - [ ] Limit the use of `abstractions/common/systemd`
    - [ ] Ensure systemctl restart/stop/reload is always confined and filtered by unit (dbus only)
    - [ ] Revisit the use of `systemd-tty-ask-password-agent`


## Next features

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
- [x] **[Play machine](;;https://github.com/roddhjav/play)**

**Build system**

- [x] Continuous release on the main branch, ~2 releases per week
- [x] Provide packages repo for ubuntu/debian
- [x] Add a `just` target to install the profiles in the right place
- [x] Fully drop the Makefile in favor of `just`

**Tests**

- [x] Tests VM for all supported targets (see [tests/vm](vm.md))
- [x] Small integration tests for a lot of core profiles (see [tests/integration](integration.md))
- [x] Autopkgtest support

**Documentation**

- [x] Initial draft of the security model and goal
- [x] Automatic docstring generated for profiles, abstraction and tunables
- [x] General documentation improvements

**Abstractions**

- [x] New `audio-client` and `audio-server` abstractions
- [x] New desktop agnostic `desktop` abstraction for all common access for any GUI app. 
- [x] New `graphics` abstraction, hardware-agnostic. Fully replace and restrict the old `opencl` abstractions
- [x] All new abstractions are documented in the [abstractions](../abstractions/index.md) page

**Dbus**

- [x] New `dbus-{system,session,accessibility}` profiles. Works regardless of the dbus implementation in use.
- [x] New talk directive: Allow the application to talk to session services. (send to)
- [x] New own directive: Allow the application to own session services under the given name. (receive, send, bind)
- [x] New `bus-{system,session,accessibility}` abstraction to be used in the profiles

**Directives**

- [x] Add directive. See the [directive](directives.md) page

