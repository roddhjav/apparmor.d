# Threat model: apparmor.d

This file condenses the security documentation of the project
(`docs/security/`: ecosystem, threat model, security model, architecture,
implementation) for the scanner. The documentation is the reference; it is
still a work in progress, discussed at
https://github.com/roddhjav/apparmor.d/discussions/1013.

## What this project does

apparmor.d is a set of 1500+ AppArmor profiles that together define the
Mandatory Access Control (MAC) policy of a whole Linux system: systemd and its
services, D-Bus, polkit, NetworkManager, Bluetooth, desktop environments
(GNOME, KDE, XFCE), user services, and common applications.

It is not a collection of independent profiles. It is one architecture: many
small, specific profiles that transition into each other (`Px`), sharing
reusable building blocks (`abstractions/`) and resource labels (`tunables/`,
the XDG-like variables). A profile defines a *world*: everything a confined
task is allowed to perceive and act upon. Abstractions are shared pieces of
that truth.

**The asset is the policy.** The Go code (`cmd/`, `pkg/`) is build and
installation tooling; it matters much less (see below).

## Attacker model

- The profiles are trusted. Confined programs are not.
- The attacker has already compromised a confined program: arbitrary code
  execution inside its process, with whatever uid it runs as, **including
  root**. AppArmor restrictions apply to root like to anyone else.
- The attacker's goal is to leave the world of that profile: run code
  unconfined or under a more privileged profile, change or unload the policy,
  persist across reboot, or reach data and resources outside the world the
  profile defines.
- A rogue or mistaken rule must not lower security below no confinement at
  all; it only widens a world.

The scan is about the profiles in this repository, as shipped. Live systems
running them (such as the project's public play machine) are out of scope.

## Where untrusted input enters

The untrusted input is the *behaviour of every confined process*: every file
path it opens, writes or executes, every binary it runs, every D-Bus message,
unix socket, signal, ptrace, mount and capability use it attempts. The policy
is the only thing that decides which of these succeed.

Network-facing and user-facing programs (sshd, NetworkManager, browsers,
media players, document viewers, Bluetooth, cups, avahi...) are the most
likely to be compromised first, so their profiles, and everything they can
transition to or talk to, matter most.

The only other untrusted input is the logs parsed by `aa-log`, which contain
strings chosen by confined processes. That attack surface is minimal (see
"The Go tooling").

## Security mode in scope

Only the **default** mode is in scope: the profiles shipped by the package.
Until release 1.0, the package installs them in complain mode by default
(`/usr/share/apparmor/modes`). This is fully configurable, globally
(`/etc/apparmor/modes`) and per profile (`flags.d/`), so **judge every
profile as if it is enforced**. In this mode unconfined processes still exist (the user's shell,
the init system), and some rules such as `@{bin}/* rPUx` fall back to
unconfined on purpose so the system does not break.

The Full System Policy (FSP), `strict` and `extreme` modes are out of scope
for now. Ignore the profiles that only apply to them (`groups/_full/`) and
the FSP build code (`pkg/configure/fsp.go`, `pkg/builder/fsp.go`).

The image installs the default policy in `/etc/apparmor.d` (complain mode,
to be judged as enforced).
Profile sources are in `/src/apparmor.d/` (`groups/`, `profiles-*/`), the
built ones in `/src/.build/apparmor.d/`, with the same subfolders as the
sources (`groups/<group>/<profile>`, `profiles-*/<profile>`). The installed
policy in `/etc/apparmor.d/` is flat (`/etc/apparmor.d/<profile>`).

## What counts as a vulnerability

A rule, or a combination of rules across profiles, abstractions and
transitions, that lets a compromised confined process escape its world.
Typical classes:

1. **Transition to unconfined or to a wider profile.** `ux` or `Ux` from a
   profile whose task does not need it; `ix` of a shell or an
   interpreter on attacker-controlled input; `change_profile` to unconfined or
   to a broader profile; a `Px` target that is itself wider than its parent
   for no reason.
2. **Write then execute (W^X).** Most profiles enforce write xor execute. A
   path writable by one profile and executed (or loaded as a library, plugin
   or configuration that runs code) by the same or another, more privileged
   profile: `@{bin}`, `@{lib}`, `/etc/ld.so.preload`, systemd unit and drop-in
   directories, udev rules, cron, `/etc/profile.d`, shell rc files, desktop
   autostart, `.desktop` files, D-Bus service files, polkit rules, sudoers,
   PAM configuration.
3. **Policy tampering.** `mac_admin` combined with write access to
   `@{sys}/kernel/security/apparmor/`, or write access to `/etc/apparmor.d/`,
   `/etc/apparmor/` or the cache, combined with any way to trigger a reload;
   loading kernel modules (`sys_module`, kmod transitions with writable module
   paths); kernel or boot configuration writes.
4. **Capabilities with their enabling rules.** A capability alone is usually
   not enough: mounting needs `sys_admin`, a matching `mount` rule and access
   to the target path. Report the combination that makes it exploitable, for
   example `mount,` or a mount rule wide enough to remount or bind over
   policy, binaries or configuration, instead of a specific
   `mount fstype=... -> <path>,`.
5. **IPC confused deputy.** `dbus`, `unix`, `signal` and `ptrace` rules that
   let a confined process drive a more privileged peer into acting for it:
   systemd (`StartTransientUnit`, unit management), polkit, PackageKit,
   udisks, logind, NetworkManager, accounts-daemon; `ptrace` or `signal`
   towards unconfined or more privileged peers; abstract or pathname unix
   sockets of privileged daemons.
6. **Data outside the world.** User data is more sensitive than the system.
   Access to user files must go through labelled directories (the XDG-like
   variables, `docs/variables.md`) and the strict user abstractions. Report
   access to secrets (`/etc/shadow`, SSH and GPG keys, keyrings, browser
   profiles, tokens) from a profile whose task does not need them, write
   access to hidden configuration in `@{HOME}`, and missing `owner` on user
   data.
7. **Privacy, when the anonymity option is enabled.** Access to identifiers
   (machine-id, MAC address, PCI and DMI data) in that mode.
8. **Attachment and transition gaps.** A profile whose attachment path does
   not cover where the binary really lives (so it runs under the parent or
   unconfined), or an `#aa:` directive, flag or prebuild step that silently
   turns enforcement into complain mode or drops rules.

## How you rate severity

Rate by the attacker's end position, starting from a compromised confined
process:

- **Critical**: from any confined process, reach unconfined code execution as
  root, or change, replace or unload the AppArmor policy, or load kernel code.
  Any full escape from a network-facing service profile.
- **High**: escape into a strictly more privileged profile (W^X chain,
  confused deputy over D-Bus, ptrace), persistence across reboot, or read
  access to secrets (`/etc/shadow`, private keys, keyrings) from a profile
  that does not need them.
- **Medium**: a W^X violation or overly wide rule with a plausible but not
  demonstrated escape chain; access to user data beyond the labelled
  directories the task needs; privacy leaks under the anonymity option.
- **Low**: a rule wider than needed with no demonstrated impact (for example
  a glob that could be a fixed path), or a deviation from the conventions that
  weakens a profile without an exploit.

A finding should state the starting profile, the exact rules involved (file
and line in `/src/apparmor.d/`), and the chain to the end position.
A demonstrated chain beats a theoretical one; without a chain, cap at medium.

## What is not a vulnerability

- **A profile that breaks a program** (missing rule, denial in the logs).
  That is a functional bug; it is important to us, but report it as a normal
  issue.
- **Missing `deny` rules.** This is a MAC policy: everything not allowed is
  denied. Deny rules are avoided on purpose.
- **A generic world doing what its documented border allows.** Some profiles
  are legitimately wide (file managers, terminals, editors, backup tools,
  sandbox managers). When a profile defines a generic world, its border is
  documented in the profile; only crossing that border is a finding.
- **Default mode fallbacks.** The user's login shell and the init system are
  unconfined, and `rPUx` fallbacks to unconfined exist by design. A finding
  is an escape *through* such a fallback that the profile did not intend, not
  the fallback itself.
- **Anything specific to the FSP, `strict` or `extreme` modes**, for now.
- **Profiles explicitly set to complain mode** per profile in
  `share/apparmor/flags.d/` (installed in `/usr/share/apparmor/flags.d/`).
  The global complain default above does not count: judge those profiles as
  enforced.
- **Ignored profiles**: not built for a distribution (`dists/ignore.d/`) or
  not installed (`share/apparmor/ignore.d/`, installed in
  `/usr/share/apparmor/ignore.d/`).
- **Bugs in the kernel or in AppArmor userspace** (`apparmor_parser`). Mention
  them if they change the impact of a rule, but they are upstream issues.
- **Upstream profiles and abstractions** shipped by AppArmor itself, unless
  apparmor.d overwrites them (`overwrite.d`).
- **Security of the confined program itself.** We do not try to prevent a
  program from being attacked; we limit what it can do once it is.

## The Go tooling

Low priority. Every tool the project ships (`aa-log`, `aa-install`,
`aa-mode`...) is itself confined by a profile in this repository, so a bug in
a tool is bounded by that profile; an escape from it is a policy finding like
any other.

- `aa-log` (`pkg/logs`) parses audit and syslog lines that contain strings
  chosen by confined processes (comm, paths, D-Bus names). The surface is
  minimal: the tool only reads logs and prints. The one case worth reporting
  is input that makes `aa-log --rules` suggest a wider rule than the denial
  justifies (medium). Crashes are low.
- `prebuild` and the `#aa:` directive parser only process trusted sources;
  report them only if they silently weaken the built policy (see class 8).

## How to exercise it

- Read profiles in `/src/apparmor.d/` and their built form in
  `/src/.build/apparmor.d/`; the installed policy in `/etc/apparmor.d/`.
- Resolve transitions by following `Px`/`px` targets to the profile with the
  matching name or attachment.
- Compile a profile without a kernel:
  `apparmor_parser -Q -K -I /etc/apparmor.d /etc/apparmor.d/<profile>`.
  `apparmor_parser -p` prints a profile with its includes expanded.
  `/etc/apparmor.d` is a snapshot made when the image was built; it does not
  reflect source edits. Validate a patch against `.build/apparmor.d` (see the
  patch guidelines).
- `just check` runs the project's static checks.
- The container has no AppArmor kernel: escapes cannot be executed live
  here. Demonstrate them from the rules.
