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

If the log file is empty, check that Auditd is running:
```sh
sudo systemctl status auditd.service
```

If Auditd is disabled aa-log will not have new results, you can enable Auditd by doing the following command:
```sh
sudo systemctl enable auditd.service --now
```

You can get more logs with:

1. `aa-log -R -s` that will provide all apparmor logs since boot time (if journalctl collect them)
2. `aa-log -R -f <nb>` where `<nb>` is `1`, `2`, `3` and `4` (the rotated audit log file)

[newissue]: https://github.com/roddhjav/apparmor.d/issues/new
[paste]: https://pastebin.com/
