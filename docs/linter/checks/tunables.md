---
title: Tunables
---

# `tunables`

Variables must be used

## Problematic rule

```sh
# WRONG
owner @{HOME}/.config/foo/{,**} rw,
```

## Correct rule

```sh
owner @{user_config_dirs}/foo/{,**} rw,
```

## Rationale

Using variables instead of hardcoding paths allows for better maintainability and compatibility across different systems and user configurations. It also makes the profile more adaptable to changes in directory structures or user environments.

See [Variables](../../variables.md) for more information on available variables.

## Exceptions

None
