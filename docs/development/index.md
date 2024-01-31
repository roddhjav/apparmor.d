---
title: Development
---

You want to contribute to `apparmor.d`, **thanks a lot for this.** Feedbacks, contributors, pull requests are all very welcome. You will find in this page all the useful information needed to contribute.

??? info "How to contribute"

    1. If you don't have git on your machine, [install it](https://help.github.com/articles/set-up-git/).
    2. Fork this repo by clicking on the fork button on the top of the [project Github][project] page.
    3. Clone the repository and go to the directory:
    ```sh
    git clone https://github.com/this-is-you/apparmor.d.git
    cd apparmor.d
    ```
    4. Create a branch:
    ```
    git checkout -b my_contribution
    ```
    5. Make the changes and commit:
    ```
    git add <files changed>
    git commit -m "A message for sum up my contribution"
    ```
    6. Push changes to GitHub:
    ```
    git push origin my_contribution
    ```
    7. Submit your changes for review: If you go to your repository on GitHub,
    you'll see a Compare & pull request button, fill and submit the pull request.


## Project rules

`Rule 1: Mandatory Access Control`

:   As these are mandatory access control policies only what is explicitly required
    should be authorized. Meaning, you should **not** allow everything (or a large area)
    and blacklist some sub areas.

`Rule 2: Do not break a program`

:   A profile **should not break a normal usage of the confined software**. It can
    be complex as simply running the program for your own use case is not always
    exhaustive of the program features and required permissions.

`Rule 3: Do not confine everything`

:   Some programs should not be confined by a MAC policy.

`Rule 4: Distribution and devices agnostic`

:   A profile should be compatible with all distributions, software and devices
    in the Linux world. You cannot deny access to resources you do not use on
    your devices or for your use case.


## Add a profile

!!! danger "Warning"

    Following the [profile guidelines](guidelines.md) is **mandatory** for all new profiles.


1. To add a new profile `foo`, add the file `foo` in [`apparmor.d/profile-a-f`][profiles-a-f]. 
   If your profile is part of a large group of profiles, it can also go in
   [`apparmor.d/groups`][groups].

2. Write the profile content, the rules depend on the confined program,
   Here is the bare minimum for the program `foo`:
``` sh
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2023 You <your@email>
# SPDX-License-Identifier: GPL-2.0-only

abi <abi/3.0>,

include <tunables/global>

@{exec_path} = @{bin}/foo
profile foo @{exec_path} {
  include <abstractions/base>

  @{exec_path} mr,

  include if exists <local/foo>
}
```


3. You can automatically set the `complain` flag on your profile by editing the file [`dists/flags/main.flags`][flags] and add a new line with: `foo complain`

4. Build & install for your distribution.


[project]: https://github.com/roddhjav/apparmor.d

[flags]: https://github.com/roddhjav/apparmor.d/blob/main/dists/flags/main.flags
[profiles-a-f]: https://github.com/roddhjav/apparmor.d/blob/main/apparmor.d/profiles-a-f
[groups]: https://github.com/roddhjav/apparmor.d/blob/main/apparmor.d/groups
