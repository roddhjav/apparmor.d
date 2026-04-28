---
title: Rule too wide
---

# `too-wide`

Rule too wide may lead to confinement escape or data leaks.

## Problematic rule

```sh
# WRONG
/tmp/** rw,
```

```sh
# WRONG
/etc/** rw,
```

## Correct rule

Limit access to only required files as much as you can. For example:

```sh
/tmp/<profile>@{rand6}/{,**} rw,
```

```sh
/etc/<profile>/** rw,
```

## Rationale

Full access to entire config and temporary directories is dangerous as it may allow confinement escape or data leaks. It is better to restrict access to only the required files or subdirectories.

## Exceptions

When a profile needs access to the full system, because it is a package manager for example.

## Related Resources

* [Access to `/tmp` breaks program isolation](https://github.com/roddhjav/apparmor.d/discussions/294)
* [Abusing Ubuntu 24.04 features for root privilege escalation](https://labs.snyk.io/resources/abusing-ubuntu-root-privilege-escalation/)
