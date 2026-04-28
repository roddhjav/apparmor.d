---
title: Abi
---

# `abi`

Use of incorrect or missing ABI version.

## Problematic rule

```sh
# WRONG
abi <abi/3.0>,
```

```sh
# WRONG
# missing 'abi <abi/4.0>,'
```

## Correct rule

```sh
abi <abi/4.0>,
```

## Rationale

All profiles in the project must use the same ABI version to ensure compatibility with the AppArmor kernel module and features. The current default ABI targeted by this project is `4.0`.

## Exceptions

None
