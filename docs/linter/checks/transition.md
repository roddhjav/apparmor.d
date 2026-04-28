---
title: Transition
---

# `transition`

`Pix` transition leads to unmaintainable profile.

## Problematic rule

```sh
# WRONG
@{bin}/foo Pix,
```

## Correct rule

```sh
@{bin}/foo ix,
```

Or, if the transition is needed:

```sh
@{bin}/foo Px,
```

## Rationale

The actual enforced transition will depend on the presence of other profiles and is therefore unpredictable. If a profile exists, it will transition and may allow more access than the original profile. If no profile exists, the program will be run inherited in the same profile, which may lead to breakage and maintenance issues.

It is **also** a security risk when used alongside a wildcard (`@{bin}/* Pix`) as, when a lot of profiles are present (like in apparmor.d) it pretty much allow to transition to any program in the system.

## Exceptions

It can be used in profile for an interactive shell environments. In this case, as long as profiles like `apt` or `apparmor_parser` it may be equivalent to giving full admin access to the user.

## Related Resources

* The [`role_play` profile of the AppArmor Play machine](https://github.com/roddhjav/play/blob/e81baf3b42513983112f3e82250710003c0dd95a/apparmor.d/groups/roles/role_play#L50-L55) use it to provide a fully confined admin role.
