---
title: Vim syntax
---

# `vim`

Missing vim syntax: `# vim:syntax=apparmor`

## Problematic rule

A profile without the vim syntax highlighting comment at the end.

## Correct rule

```sh
profile foo {
  ...
}

# vim:syntax=apparmor

```

## Rationale

The vim syntax highlighting comment enables proper syntax highlighting when editing the profile in vim. This improves readability and helps prevent syntax errors.

## Exceptions

None
