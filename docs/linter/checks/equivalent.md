---
title: Equivalent
---

# `equivalent`

Missing rules to equivalent paths.

## Problematic rule

```sh
# WRONG
@{bin}/grep ix,
```

## Correct rule

```sh
@{bin}/{,e}grep ix,
```

## Rationale

In AppArmor profiles, certain binaries may have equivalent paths that need to be explicitly allowed to ensure proper functionality.

For example, the `grep` binary can be accessed as both `grep` and `egrep` (`exec grep -E "$@"`). Failing to include rules for all equivalent paths **will** lead to unexpected denials in some distributions.

The following equivalent paths exist:

- `@{bin}/awk` -> `@{bin}/{m,g,}awk`
- `@{bin}/grep` -> `@{bin}/{,e}grep`
- `@{bin}/gs` -> `@{bin}/gs{,.bin}`
- `@{bin}/which` -> `@{bin}/which{,.debianutils}`
- `@{sbin}/xtables-legacy-multi` -> `@{sbin}/xtables-{nft,legacy}-multi`
- `@{bin}/xtables-nft-multi` -> `@{sbin}/xtables-{nft,legacy}-multi`

## Exceptions

None
