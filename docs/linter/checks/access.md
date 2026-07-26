---
title: Access
---

# `access`

File access permissions written in the wrong order.

## Problematic rule

```sh
# WRONG
/usr/** rwkl,
```

## Correct rule

```sh
/usr/** rwlk,
```

## Rationale

File access permissions must be written in a consistent, canonical order. Although the linter only checks the `rwkl`/`rwlk` case, all access permissions should follow this order: `rwlkmix`.

## Exceptions

None
