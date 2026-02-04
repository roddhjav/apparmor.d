---
title: Flatpak abstractions
tags:
  - abstractions
  - flatpak
  - system
---

These abstractions should only be used by the flatpak profiles. They provide the necessary rules to run Flatpak applications confined with AppArmor. They are designed to very closely match the [Flatpak Sandbox Permissions](https://docs.flatpak.org/en/latest/sandbox-permissions.html). Therefore, they are different to they host equivalents, as flatpak apps do not have access to the full host filesystem.
