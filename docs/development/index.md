---
title: Development
---

If you're looking to contribute to `apparmor.d` you can get started by going to the project [GitHub repository](https://github.com/roddhjav/apparmor.d/)! All contributions are welcome no matter how small. In this page you will find all the useful information needed to contribute to the apparmor.d project.

??? info "How to contribute pull requests?"

    1. If you don't have git on your machine, [install it](https://help.github.com/articles/set-up-git/).
    1. Fork this repo by clicking on the fork button on the top of the [project GitHub](https://github.com/roddhjav/apparmor.d) page.
    1. [Generate a new SSH key](    https://docs.github.com/en/authentication/connecting-to-github-with-ssh/generating-a-new-ssh-key-and-adding-it-to-the-ssh-agent) and add it to your GitHub account.
    1. Clone the forked repository and go to the directory:
    ```sh
    git clone git@github.com:your-github-username/apparmor.d.git 
    cd apparmor.d
    ```
    1. Create a branch:
    ```
    git checkout -b my_contribution
    ```
    1. Make the changes and commit:
    ```
    git add <files changed>
    git commit -m "A message to sum up my contribution"
    ```
    1. Push changes to GitHub:
    ```
    git push origin my_contribution
    ```
    1. Submit your changes for review: If you go to your repository on GitHub,
    you'll see a Compare & pull request button, fill and submit the pull request.

<div class="grid cards" markdown>

-   :material-arrow-right: &nbsp; **[See the workflow to write profiles](workflow.md)**

</div>


## Project rules

#### Rule :material-numeric-1-circle: - Mandatory Access Control

:   As these are mandatory access control policies **only** what is explicitly required
    should be authorized. Meaning, you should **not** allow everything (or a large area)
    and deny some sub areas.

#### Rule :material-numeric-2-circle: - Do not break a program

:   A profile **should not break a normal usage of the confined software**. this can
    be complex as simply running the program for your own use case is not always
    exhaustive of the program features and required permissions.

#### Rule :material-numeric-3-circle: - Do not confine everything

:   Some programs should not be confined by a MAC policy.

#### Rule :material-numeric-4-circle: - Distribution and devices agnostic

:   A profile should be compatible with all distributions, software, and devices
    in the Linux world. You cannot deny access to resources you do not use on
    your devices or for your use case.


## Recommended documentation

* [The AppArmor Core Policy Reference](https://gitlab.com/apparmor/apparmor/-/wikis/AppArmor_Core_Policy_Reference)
* [The openSUSE Documentation](https://doc.opensuse.org/documentation/leap/security/html/book-security/part-apparmor.html)
* [SUSE Documentation](https://documentation.suse.com/sles/12-SP5/html/SLES-all/cha-apparmor-intro.html)
* [The AppArmor.d man page](https://man.archlinux.org/man/apparmor.d.5)
* [F**k AppArmor](https://presentations.nordisch.org/apparmor/#/)
* [A Brief Tour of Linux Security Modules](https://www.starlab.io/blog/a-brief-tour-of-linux-security-modules)
