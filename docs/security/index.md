---
title: Security
---

There are over 50000 Linux packages and even more applications. It is simply not possible to write an AppArmor profile for all of them. Therefore, a question arises:

<figure markdown>

**What to confine, how, and why?**

</figure>

The security model presented here help us to scope the security policies within the broader context of system security and privacy.

This section presents the security model considered for the profiles in `apparmor.d`. Despite that this security model looks at Linux security in general, we are only focusing on the threats, model, and implementation within the scope of AppArmor.

!!! warning

    This security model is still a work in progress. Comments and feedbacks are welcome. [Discuss it on Github](https://github.com/roddhjav/apparmor.d/discussions/1013)

<div class="grid cards" markdown>

-   :simple-linux: &nbsp; **[:material-numeric-1-circle-outline: Ecosystem Review](ecosystem.md)**

    ---

    What are Linux based systems used for? By whom and how?

-   :material-bomb: &nbsp; **[:material-numeric-2-circle-outline: Threat model](threat.md)**

    ---

    The list of threats the <span class="pg-red">:material-numeric-1-circle-outline: Ecosystem</span> is continuously facing.

-   :material-security: &nbsp; **[:material-numeric-3-circle-outline: Security Model](model.md)**

    ---

    Given the <span class="pg-red">:material-numeric-1-circle-outline: Ecosystem</span> and the <span class="pg-green">:material-numeric-2-circle-outline: Threat model</span>, what are the rules that should be enforced?

-   :material-source-commit-local: &nbsp; **[:material-numeric-4-circle-outline: Security Implementation](architectureimplementation.md)**

    ---

    How AppArmor is used to enforce part of the <span class="pg-blue">:material-numeric-3-circle-outline: Security Model</span>?

</div>


!!! quote "Security Model"

    A computer security model is a scheme for specifying and enforcing security policies. A security model may be founded upon a formal model of access rights, a model of computation, a model of distributed computing, or no particular theoretical grounding at all. 

    *Source: [Wikipedia](https://en.wikipedia.org/wiki/Computer_security_model)*

!!! example "References"

    1. [The Android Platform Security Model (2023)](https://arxiv.org/pdf/1904.05572v3.pdf)
    1. [ClipOS](https://docs.clip-os.org/) - A security OS made by the ANSSI (the French NIST) and used for sensitive French government related activities.
    1. [Spectrum](https://spectrum-os.org) - A step towards usable secure computing
    1. [QubesOS](https://www.qubes-os.org/) - A reasonably secure operating system
    1. [Whonix](https://www.whonix.org/) – An anonymous operating system
    1. [Kairos](https://kairos.io)

