---
title: Local include
---

# `include`

Missing inclusion of local rule additions.

## Problematic rule

```sh
# WRONG
profile pass {
  ...
}

```

## Correct rule

```sh
profile pass {
  include if exists <local/pass>
}
```

## Rationale

To allow for easier customization and extension of AppArmor profiles and subprofiles, all profiles and abstractions **must** include local rule additions.

## Exceptions

None
