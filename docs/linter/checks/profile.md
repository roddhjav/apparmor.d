---
title: Profile / Subprofile
---

# `profile` / `subprofile`

Missing or incorrect profile name.

## Problematic rule

```sh
cat /etc/apparmor.d/foo
# WRONG
profile myfoo {
  ...
}
```

```sh
cat /etc/apparmor.d/foo
# WRONG
profile @{bin}/foo {
  ...
}
```

## Correct rule

```sh
cat /etc/apparmor.d/foo
profile foo {
  ...
}
```

```sh
cat /etc/apparmor.d/foo
# WRONG
profile foo @{bin}/foo {
  ...
}
```

## Rationale

AppArmor profiles and subprofiles **must** have a name that matches the filename of the profile. Old syntax that includes profile attachment instead of a profile name **must** be avoided.

## Exceptions

None
