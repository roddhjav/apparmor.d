---
title: Directory Mark
---

# `directory-mark`

Missing directory mark (trailing slash) in well-known directory paths.

## Problematic rule

```sh
# WRONG
owner @{HOME} r,
```

## Correct rule

```sh
owner @{HOME}/ r,
```

## Rationale

In AppArmor profiles, a directory path **must** be explicitly marked with a trailing slash (`/`) to indicate that it refers to a directory rather than to a file.

## Exceptions

None
