---
title: Enforce Mode
---

# Enforce Mode

The default package configuration installs all profiles in *complain* mode.
Once you tested have them and it works fine, you can easily switch to *enforce* mode.
To do this, edit `PKGBUILD` on Archlinux or `debian/rules` on Debian and remove 
the `--complain` option to the configure script. Then build the package as usual:
```diff
-  ./configure --complain
+  ./configure
```

Do not worry, the profiles that are not considered stable are kept in complain mode.
They can be tracked in the [`dists/flags`](https://github.com/roddhjav/apparmor.d/tree/master/dists/flags) directory.
