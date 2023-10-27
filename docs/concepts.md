---
title: Concepts
---

# Concepts

*One profile a day keeps the hacker away*

There are over 50000 Linux packages and even more applications. It is simply not possible to write an AppArmor profile for all of them. Therefore, a question arises:

**What to confine and why?**

We take inspiration from the [Android/ChromeOS Security Model][android_model], and we apply it to the Linux world. Modern [Linux security distributions][clipos] usually consider an immutable core base image with a carefully selected set of applications. Everything else should be sandboxed. Therefore, this project tries to confine all the *core* applications you will usually find in a Linux system: all systemd services, xwayland, network, bluetooth, your desktop environment...  Non-core user applications are out of scope as they should be sandboxed using a dedicated tool (minijail, bubblewrap, toolbox...).

This is fundamentally different from how AppArmor is usually used on Linux servers as it is common to only confine the applications that face the internet and/or the users.


[android_model]: https://arxiv.org/pdf/1904.05572v2.pdf
[clipos]: https://clip-os.org/en/
[write xor execute]: https://en.wikipedia.org/wiki/W%5EX

