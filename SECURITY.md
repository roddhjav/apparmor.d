# Security Policy

## Supported versions

Only the latest release and the `main` branch receive security fixes.

## Reporting a vulnerability

Please do not open a public issue for a security problem. Report it privately
by email to **security@pujol.io**.

Include, when possible:

- the profile(s) and rules involved (file and line),
- the chain from the starting profile to the end position (escape, policy
  tampering, access to data outside the profile),
- the distribution and the AppArmor version.

You should receive an answer within a few days. Fixes are released as soon
as possible, and you are credited unless you ask otherwise.

## Scope

The asset is the policy: the profiles, abstractions and tunables in
`apparmor.d/`. What counts as a vulnerability, and how severity is rated, is
described in the [threat model](.oss-scanner/threat_model.md) and in the
[security documentation](https://apparmor.pujol.io/security/).

A profile that breaks a program is a functional bug, not a security issue:
report it as a normal issue.
