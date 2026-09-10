---
title: Enforce Mode
---

The default package configuration installs all profiles in *complain* mode. This is a safety measure to ensure you are not going to break your system on initial installation. Once you have tested it, and it works fine, you can easily switch to *enforce* mode. The profiles that are not considered stable are kept in complain mode, they can be tracked in the [`dists/flags`](https://github.com/roddhjav/apparmor.d/tree/main/dists/flags) directory.

!!! danger

    - You **must** test in complain mode first and ensure your system works as expected.
    - You **must** regularly check AppArmor log with [`aa-log`](usage.md#apparmor-log) and [report](report.md) issues first.
    - When reporting an issue, you **must** ensure the affected profiles are in complain mode.

**Configuration**

Set the default mode of `aa-install` to enforce in `/etc/apparmor/modes`

```sh
cat <<-EOF | sudo tee /etc/apparmor/modes
default enforce
EOF
```

**Installation**

To apply the change, start aa-install manually:

```sh
sudo aa-install
```