---
title: Linter
---

The profiles are checked for common style and security issues with `just check`. This page documents all the checks performed.

!!! note "Check system"

    Future implementation will expand this basic check system to a full linter and security analyzer system.


## Overview

| Output | Check ID | Description |
|---|---|---|
| <span class="pg-red">**:fontawesome-solid-ban:**</span> | `abi` | Missing ABI |
| <span class="pg-red">**:fontawesome-solid-ban:**</span> | `abstractions` | Use of dangerous abstraction |
| <span class="pg-red">**:fontawesome-solid-ban:**</span> | `abstractions` | Use of deprecated abstraction |
| <span class="pg-red">**:fontawesome-solid-ban:**</span> | `abstractions` | Use of deprecated, ubuntu only abstraction |
| <span class="pg-red">**:fontawesome-solid-ban:**</span> | `bin` | Use of `@{sbin}` instead of `@{bin}` |
| <span class="pg-red">**:fontawesome-solid-ban:**</span> | `directory-mark` | Missing directory mark |
| <span class="pg-red">**:fontawesome-solid-ban:**</span> | `equivalent` | Missing equivalent program |
| <span class="pg-red">**:fontawesome-solid-ban:**</span> | `header` | Missing header |
| <span class="pg-red">**:fontawesome-solid-ban:**</span> | `include` | Missing include |
| <span class="pg-red">**:fontawesome-solid-ban:**</span> | `indentation` | Invalid indentation |
| <span class="pg-red">**:fontawesome-solid-ban:**</span> | `profile` | Missing profile name |
| <span class="pg-red">**:fontawesome-solid-ban:**</span> | `sbin` | Use of `@{bin}` instead of `@{sbin}` |
| <span class="pg-red">**:fontawesome-solid-ban:**</span> | `subprofiles` | Missing subprofiles |
| <span class="pg-red">**:fontawesome-solid-ban:**</span> | `tabs` | Tabs are not allowed |
| <span class="pg-red">**:fontawesome-solid-ban:**</span> | `trailing` | Line has trailing whitespace |
| <span class="pg-red">**:fontawesome-solid-ban:**</span> | `transition` | `Pix` transition leads to unmaintainable profile |
| <span class="pg-red">**:fontawesome-solid-ban:**</span> | `transition` | Executable should be should be used inherited: `ix` or `Cx` |
| <span class="pg-red">**:fontawesome-solid-ban:**</span> | `transition` | Executable should transition to another (sub)profile with `Px` or `Cx` |
| <span class="pg-red">**:fontawesome-solid-ban:**</span> | `tunables` | Variables must be used |
| <span class="pg-red">**:fontawesome-solid-ban:**</span> | `udev` | Udev data path without a description comment |
| <span class="pg-red">**:fontawesome-solid-ban:**</span> | `useless` | Rule already included in the base abstraction |
| <span class="pg-red">**:fontawesome-solid-ban:**</span> | `vim` | Missing vim syntax |
| <span class="pg-orange">**:fontawesome-solid-warning:**</span> | `transition` | Path `@{bin}/XXX` should transition to a subprofile with 'Cx' |
|<span class="pg-orange">**:fontawesome-solid-warning:**</span> | `too-wide` | Rule too wide may lead to confinement escape or data leaks |

## Directive

We use a special [directive](directives.md) to ignore specific checks:

- Inline directive is supported
- Directive before a paragraph applies to all rules in the paragraph
- Directive within the first 5 lines of a file applies to the whole file

**Format**

```sh
#aa:lint ignore=<check>
```

**`<check>`**

:   Check id to ignore.


**Example**

Ignore the `too-wide` check in the `dpkg` profile:

!!! quote ""

    **[apparmor.d/groups/apt/dpkg](https://github.com/roddhjav/apparmor.d/blob/094795cc6d628923b7454fd3a9289c44891edc62/apparmor.d/groups/apt/dpkg#L52-L61)**
    ``` sh linenums="52"
      #aa:lint ignore=too-wide
      # Install/update packages
      / r,
      /*{,/} rw,
      @{efi}/** rwl -> @{efi}/**,
      /etc/** rwl -> /etc/**,
      /opt/** rwl -> /opt/**,
      /srv/** rwl -> /srv/**,
      /usr/** rwlk -> /usr/**,
      /var/** rwlk -> /var/**,
    ```

## Description Template

    ---
    title: id
    ---

    # `id`

    <description of the check>

    ## Problematic rule

    ```sh
    # WRONG
    <example of problematic rule>
    ```

    ## Correct rule

    ```sh
    <example of correct rule>
    ```

    ## Rationale

    <explanation of why the correct rule is better>

    ## Exceptions

    None

    ## Related Resources
