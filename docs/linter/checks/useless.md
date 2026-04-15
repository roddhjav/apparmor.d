---
title: Useless rule
---

# `useless`

Rule already included in the base abstraction, remove it.

## Problematic rule

```sh
# WRONG
@{sys}/devices/system/cpu/online r,
```

## Correct rule

```sh
# CORRECT
# Rule already included in the base abstraction, no need to include it again
```

## Exceptions

None

## Related Resources

* The [`base-strict` abstraction](https://github.com/roddhjav/apparmor.d/blob/main/apparmor.d/abstractions/base-strict)
