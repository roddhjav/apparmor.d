---
title: Report AppArmor logs
---

# Report AppArmor logs

The **[aa-log](usage.md#apparmor-log)** tool reports all AppArmor `DENIED` and `ALLOWED`. It should be used to fix AppArmor related issues.

While testing, if something get wrong, you need to put the profile in complain mode, to that you can investigate and it does not block your program.

When creating [an issue on Github][newissue]. Please ensure you post a link to the [paste] of the AppArmor audit log: `/var/log/audit/audit.log`.

[newissue]: https://github.com/roddhjav/apparmor.d/issues/new
[paste]: https://pastebin.com/
