---
title: Report AppArmor logs
---

The **[aa-log](usage.md#apparmor-log)** tool reports all AppArmor `DENIED` and `ALLOWED`. It should be used to fix AppArmor related issues.

While testing, if something get wrong, you need to put the profile in complain mode, so that you can investigate, and it does not block your program.

When creating [an issue on Github][newissue], please post a link to the [paste] of the audit log generated with:
```sh
aa-log -R
```

If this command produce nothing, try:
```sh
aa-log -s -R
```

[newissue]: https://github.com/roddhjav/apparmor.d/issues/new
[paste]: https://pastebin.com/
