# How we would like patches made

## Size and scope of a fix

- As narrow as possible, and in the right place: fix the abstraction or tunable
  when the problem comes from it, the profile otherwise. Do not refactor
  unrelated rules.
- Narrow the rule rather than remove it when the program needs the access:
  a fixed path instead of a glob, `owner`, a specific `mount fstype=... ->
  <path>,` instead of `mount,`, a specific D-Bus interface and member, a `Px`
  to a dedicated child profile instead of `ix` or `Ux`. A fix must not break
  the program (project rule 2: do not break programs).
- If a dedicated child profile is needed, add it as a new file following the
  same conventions.

## Profile style

- No `deny` rules: only allow what is needed.
- Use the variables, never raw paths: `@{bin}`, `@{lib}`, `@{HOME}`,
  `@{user_config_dirs}`, `@{run}`, `@{sys}`, `@{PROC}`... and the labelled user
  directories instead of `@{HOME}/**`.
- Use the project D-Bus abstractions (`bus/...`, `bus-session`,
  `bus-system`) or `#aa:dbus` directives, not the upstream `dbus-*` ones.
- Keep the rule block order: include, set rlimit, userns, capability, network,
  mount/remount/umount, pivot_root, change_profile, signal, ptrace, unix, dbus,
  file, `include if exists <local/...>`.
- File rules order: entry point, binaries and libraries, `/usr/share`, `/etc`,
  `/var`, user data, user config, temp and runtime, `/sys`, `/proc`, `/dev`.
  Sort alphabetically within a block.
- Do not use the em dash character, in rules, comments or commit messages.

## Tests

- The patched profile must compile:
  `just prebuild && apparmor_parser -Q -K -I .build/apparmor.d .build/apparmor.d/<profile>`
  (or `just parser` for all of them).
- `just check` must pass.
- If Go code is changed, add a case to the existing table-driven test of the
  package (see `pkg/configure/flags_test.go` for the style).

## Branch and commit message

- Against `main`.
- Subject in the form `fix(<profile, abs, tool name>): <name>: <what changed>`, for example
  `fix(profile): kmod: restrict module path to signed modules`.
